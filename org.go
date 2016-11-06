package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

//Org name json struct
type Org struct {
	Login string `json:"login"`
}

// Team json struct
type Team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

//AskForDefaultContext returns the default context that the user would like to use
func AskForDefaultContext() string {
	fmt.Println("GoGit can operate out three different contexts: Organizational, Team and Personal. Which filters the resultsets of the various commands accordingly")
	fmt.Println("This is overridable with the global flag --context team/org/Personal (or --c team/org/personal)")
	context := getString("What default context would you like to use? (team,org,personal): ")
	if !stringInSlice(context, []string{"org", "team", "personal"}) {
		log.Fatalln("input not in available options")
	}
	return context
}

//AskForDefaultOrg returns the org that the user chooses out of her available orgs
func AskForDefaultOrg(username string, token string) string {
	orgs := ListMyOrgs(username, token)
	fmt.Println("Which organisation would you want as default? (this option can be overridden in individual calls)")
	PrintTable([]string{"OrgName"}, PivotStringSlice(orgs))
	org := strings.ToLower(getString("Organisation name: "))
	if !stringInSlice(org, orgs) {
		log.Fatalln("input not in available options")
	}
	return org
}

//AskForDefaultTeam returns the team that the user chooses out of her available teams in the given org
func AskForDefaultTeam(username string, token string, org string) int {
	teams := ListOrgTeams(username, token, org)
	fmt.Printf("%v", teams)
	fmt.Println("Which team would you want as default? (this option can be overridden in individual calls)")
	PrintTable([]string{"TeamId", "TeamName"}, pivotTeamList(teams))
	team := getInt("Team ID (number): ")
	if !intInSlice(team, teamSliceToIntSlice(teams)) {
		log.Fatalln("input not in available options")
	}
	return team
}

//ListMyOrgs returns a list of all organisations that you are part of
func ListMyOrgs(username string, token string) []string {
	resp, err := GetFromGitHub(username, token, "/user/orgs")
	if err != nil {
		log.Fatal(err)
	}
	var orgList []Org
	jsonErr := json.Unmarshal(resp, &orgList)
	if jsonErr != nil {
		log.Fatal(err)
	}
	var returnList []string
	for _, s := range orgList {
		returnList = append(returnList, s.Login)
	}
	return returnList
}

//ListOrgTeams returns a list of all teams in an org
func ListOrgTeams(username string, token string, org string) []Team {
	resp, err := GetFromGitHub(username, token, fmt.Sprintf("/orgs/%v/teams", org))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(resp))
	var result []Team
	jsonErr := json.Unmarshal(resp, &result)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return result
}

func pivotTeamList(data []Team) [][]string {
	var result [][]string
	for _, s := range data {
		row := []string{strconv.Itoa(s.ID), s.Name}
		result = append(result, row)
	}
	return result
}

func teamSliceToIntSlice(list []Team) []int {
	var result []int
	for _, i := range list {
		result = append(result, i.ID)
	}
	return result
}
