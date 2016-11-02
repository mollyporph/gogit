package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gogit"
	app.Usage = "tbc"
	app.Action = func(c *cli.Context) error {
		fmt.Println("It's working..somewhat")
		return nil
	}

	app.Run(os.Args)
}
