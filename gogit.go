package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"

	"path"

	"net/http"

	"bytes"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"
)

//Config .gogit config structure
type Config struct {
	Token, preferredcontext string
	Colorize                bool
}

func getPatPermnissions() []string {
	return []string{"repo", "read:org", "user"}
}

func main() {
	var config Config
	var verbose bool
	var context string
	preflight(&config)
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

func preflight(config *Config) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	configFilePath := path.Join(usr.HomeDir, ".gogit")
	_, fileErr := ioutil.ReadFile(configFilePath)
	if fileErr != nil {
		handleConfigError(config, configFilePath)
	}
	file, _ := ioutil.ReadFile(configFilePath)
	jsonErr := json.Unmarshal(file, config)
	if jsonErr != nil {
		handleConfigError(config, configFilePath)
	}
	if config.Token == "" {
		handleConfigError(config, configFilePath)
	}
}
func handleConfigError(config *Config, configFilePath string) {
	if askForConfirmation(fmt.Sprintf("Woops! It seems like you don't have a valid .gogit config at %v, want to set one up?", configFilePath)) {
		setup(config, configFilePath)
	} else {
		fmt.Printf("Can't continue without a proper config file, exiting...")
		os.Exit(1)
	}
}
func setup(config *Config, configFilePath string) {

	pat := createGithubPAT()
	config.Token = pat
	json, jsonErr := json.Marshal(&config)
	if jsonErr != nil {
		log.Fatalf("Could not marshal config %v to json", &config)
	}
	fileErr := ioutil.WriteFile(configFilePath, json, 0666)
	if fileErr != nil {
		log.Fatalf("Could not create config file at %v", configFilePath)
	}
	config.Colorize = askForConfirmation("would you like colorized output?")
	fmt.Println("Your .gogit config file is now complete! run `gogit help` to begin using GoGit!")
	os.Exit(0)
}

// TokenRequestBody request for token call
type TokenRequestBody struct {
	Scopes []string `json:"scopes"`
	Note   string   `json:"note"`
}

// TokenResponseBody response from token
type TokenResponseBody struct {
	Token string
}

func createGithubPAT() string {

	reader := bufio.NewReader(os.Stdin)
	patPermissions := getPatPermnissions()
	askForPatConfirmation(patPermissions)
	fmt.Println("To create a PAT gogit will need your github username and password, and a multifactor authentication token if you have it enabled.")
	fmt.Println("Your password will only be used to create a PAT (over https) and will not be stored anywhere")
	fmt.Printf("github username: ")
	username, err := reader.ReadString('\n')
	password := readPassword()
	requestBody := TokenRequestBody{
		Scopes: patPermissions,
		Note:   "GoGit PAT", //todo: make changeable}
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatal(err)
	}
	tokenBody := postToGithub(strings.TrimSpace(username), string(password), "/authorizations", string(body), "")
	var token TokenResponseBody
	jsonErr := json.Unmarshal(tokenBody, &token)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return token.Token
}
func readPassword() string {
	fd := os.Stdin.Fd()
	oldState, err := terminal.MakeRaw(int(fd))
	if err != nil {
		panic(err)
	}
	fmt.Printf("github password: ")
	password, err := terminal.ReadPassword(int(fd))
	terminal.Restore(int(fd), oldState)
	fmt.Println()
	return string(password)
}
func askForPatConfirmation(patPermissions []string) {
	fmt.Println("GoGit needs to create a personal access token to be able to access github's API.")
	s := fmt.Sprintf("We will create a PAT with the following permissions: %v. for more info on oauth scopes, visit https://developer.github.com/v3/oauth/#scopes", patPermissions)
	fmt.Println(s)
	fmt.Println("The token will be saved in your .gogit config file")
	if !askForConfirmation("Create PAT token?") {
		fmt.Println("Cannot continue without correct config")
		os.Exit(1)
	}
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
func postToGithub(usr string, pwOrToken string, route string, body string, mfaToken string) []byte {
	var returnMsg []byte
	baseURL := "https://api.github.com"
	fullURL := baseURL + route
	client := &http.Client{}
	req, err := http.NewRequest("POST", fullURL, bytes.NewBufferString(body))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(usr, pwOrToken)
	if mfaToken != "" {
		req.Header.Add("X-GitHub-OTP", mfaToken)
	}
	resp, err := client.Do(req)
	if resp.StatusCode == 401 {
		opt := resp.Header.Get("X-GitHub-OTP")
		if opt != "" && strings.HasPrefix(opt, "required") { //Handle two-factor auth

			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("Please enter two-factor authentication token: ")
			mfa, mfaErr := reader.ReadString('\n')

			if mfaErr != nil {
				log.Fatal(mfaErr)
			}
			return postToGithub(usr, pwOrToken, route, body, strings.TrimSpace(mfa))
		}
		log.Fatal("Unauthorized! Please check your username and password. Or run gogit setup to set up a new public access token")
	} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		returnMsg = respBody
	} else {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		msg := fmt.Sprintf("status: %v body: %v", resp.StatusCode, string(respBody))
		fmt.Println(msg)
		log.Fatal("something something..darkside")
	}
	return returnMsg
}
