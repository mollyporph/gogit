package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"

	"encoding/json"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	"gopkg.in/resty.v0"
)

type config struct {
	token, preferredcontext string
	colorize                bool
}

func getPatPermnissions() []string {
	return []string{"repo", "read:org", "user"}
}

func main() {
	//	var config config
	var verbose bool
	var context string
	//preflight(&config)
	app := cli.NewApp()
	app.Name = "gogit"
	app.Usage = "tbc"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "context, c",
			Value:       "org",
			Usage:       "Context of the call, org, team or personal",
			Destination: &context,
		},
		cli.BoolFlag{
			Name:        "verbose, vb",
			Usage:       "If you want verbose output.",
			Destination: &verbose,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "pullrequest",
			Aliases: []string{"pr"},
			Usage:   "add a task to the list",
			Action: func(c *cli.Context) error {
				fmt.Println(context)
				data := [][]string{
					[]string{"1", "2", "3", "4"},
				}
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"1", "2", "3", "4"})
				table.AppendBulk(data) // Add Bulk Data
				table.Render()
				return nil
			},
		},
	}

	app.Run(os.Args)
}

func preflight(config *config) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	file, fileErr := ioutil.ReadFile(usr.HomeDir + "/.gogit")
	if fileErr != nil {
		handleConfigError(config)
	}
	jsonErr := json.Unmarshal(file, config)
	if jsonErr != nil {
		handleConfigError(config)
	}
}
func handleConfigError(config *config) {
	if askForConfirmation("Woops! It seems like you don't have a valid .gogit config, want to set one up?") {
		setup(config)
	} else {
		fmt.Printf("Can't continue without a proper config file, exiting...")
		os.Exit(1)
	}
}
func setup(config *config) {
	pat := createGithubPAT()
	config.token = pat
}
func createGithubPAT() string {
	fmt.Println("GoGit needs to create a personal access token to be able to access github's API.")
	patPermissions := getPatPermnissions()
	s := fmt.Sprintf("We will create a PAT with the following permissions: %v. for more info on oauth scopes, visit https://developer.github.com/v3/oauth/#scopes", patPermissions)
	fmt.Println(s)
	fmt.Println("The token will be saved in your .gogit config file")
	if !askForConfirmation("Create PAT token?") {
		fmt.Println("Cannot continue without correct config")
		os.Exit(1)
	}
	return ""
}
func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [y/n]: ", s)
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
func postToGithub(usr string, pwOrToken string, route string, body string, mfaToken string) string {
	var respBody string
	baseURL := "https://api.github.com"
	request := resty.R().
		SetBody(body).
		SetResult(&respBody).
		SetBasicAuth(usr, pwOrToken)
	if mfaToken != "" {
		request.SetHeader("X-GitHub-OTP", mfaToken)
	}
	resp, err := request.Post(baseURL + route)
	if err != nil {
		if resp.StatusCode() == 401 {

			opt := resp.Header().Get("X-GitHub-OTP")

			if opt != "" && strings.HasPrefix(opt, "required") { //Handle two-factor auth

				reader := bufio.NewReader(os.Stdin)
				fmt.Printf("Please enter two-factor authentication token: ")
				mfa, mfaErr := reader.ReadString('\n')

				if mfaErr != nil {
					log.Fatal(mfaErr)
				}

				return postToGithub(usr, pwOrToken, route, body, mfa)
			}
			log.Fatal("Unauthorized! Please check your username and password. Or run gogit setup to set up a new public access token")
		}
	}
	return respBody
}
