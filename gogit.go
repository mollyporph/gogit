package main

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

//Config .gogit config structure
type Config struct {
	Username, Token, preferredcontext, defaultOrg, defaultTeam string
	Colorize                                                   bool
}

//State application state, such as supplied parameters
type State struct {
	Username, Organization, Team, Context string
	Verbose                               bool
}

func main() {
	var config Config
	var state State
	preflight(&config)
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
			Action: func(c *cli.Context) error {
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

//GetGithubPatAndUsername returns the created PAT and the correct username
