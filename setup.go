package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
)

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
	config.Token, config.Username = GetGithubPatAndUsername()
	fmt.Println("Github PAT successfully created!")
	config.Colorize = false // askForConfirmation("would you like colorized output?")
	config.DefaultOrg = AskForDefaultOrg(config.Username, config.Token)
	config.DefaultTeamID = AskForDefaultTeam(config.Username, config.Token, config.DefaultOrg)
	config.DefaultContext = AskForDefaultContext()
	json, jsonErr := json.Marshal(&config)
	if jsonErr != nil {
		log.Fatalf("Could not marshal config %v to json", &config)
	}
	fileErr := ioutil.WriteFile(configFilePath, json, 0666)
	if fileErr != nil {
		log.Fatalf("Could not create config file at %v", configFilePath)
	}
	fmt.Println("Your .gogit config file is now complete! run `gogit help` to begin using GoGit!")
	os.Exit(0)
}
