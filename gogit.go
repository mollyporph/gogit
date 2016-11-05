package main

import (
	"os"

	"fmt"

	"github.com/urfave/cli"
)

var state State
var config Config

//Config .gogit config structure
type Config struct {
	Username, Token, DefaultContext, DefaultOrg string
	Colorize                                    bool
	DefaultTeamID                               int
}

//State application state, such as supplied parameters
type State struct {
	Username, Token, Organization, Context string
	Verbose                                bool
	TeamID                                 int
}

func buildState() {
	state.Username = config.Username
	state.Token = config.Token
	if state.Context == "" {
		state.Context = config.DefaultContext
	}
	if state.TeamID == 0 {
		state.TeamID = config.DefaultTeamID
	}
	if state.Organization == "" {
		state.Organization = config.DefaultOrg
	}
}
func main() {
	app := cli.NewApp()
	app.Name = "gogit"
	app.Usage = "tbc"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "context, c",
			Value:       "org",
			Usage:       "Context of the call, org, team or personal",
			Destination: &state.Context,
		},
		cli.StringFlag{
			Name:        "org, o",
			Value:       "",
			Usage:       "The organization you want to query against",
			Destination: &state.Organization,
		},
		cli.IntFlag{
			Name:        "teamid, t",
			Value:       0,
			Usage:       "The team you want to query against (you need to supply `--context team` for this to work)",
			Destination: &state.TeamID,
		},
		cli.BoolFlag{
			Name:        "verbose, vb",
			Usage:       "If you want verbose output.",
			Destination: &state.Verbose,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "pullrequest",
			Aliases: []string{"pr"},
			Usage:   "add a task to the list",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "status, s",
					Value: "open",
					Usage: "Filters the pullrequests based on `status`. Available options are open,closed,all",
				},
			},
			Action: func(c *cli.Context) error {
				buildState()
				PrintPullRequests(c.String("status"))
				return nil
			},
		},
		{
			Name:    "organisations",
			Aliases: []string{"orgs"},
			Usage:   "Lists your available orgs",
			Action: func(c *cli.Context) error {
				buildState()
				orgs := ListMyOrgs(state.Username, state.Token)
				for _, s := range orgs {
					fmt.Println(s)
				}
				return nil
			},
		},
	}
	preflight(&config)
	app.Run(os.Args)
}

//GetGithubPatAndUsername returns the created PAT and the correct username
