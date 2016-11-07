package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

//Issue json struct
type Issue struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	ClosedAt  string    `json:"closed_at"`
	User      struct {
		Name string `json:"login"`
	} `json:"user"`
	Assignee struct {
		Name string `json:"login"`
	} `json:"assignee"`
}

//PrintIssues prints issues
func PrintIssues(status string) {
	issues := GetIssues(status)
	PrintTable([]string{"ID", "Title", "Created by", "assigned to", "State", "Created"}, issuesToTable(issues))
}

//GetIssues gets issues
func GetIssues(status string) []Issue {
	switch state.Context {
	case "org":
		return getOrgIssues(status)
	case "team":
		return getTeamIssues(status)
	case "personal":
		return getMyIssues(status)
	}
	return nil
}

func getMyIssues(status string) []Issue {
	return getIssues(fmt.Sprintf("/issues?state=%v&filter=all&page_size=100", status))
}
func getOrgIssues(status string) []Issue {
	return getIssues(fmt.Sprintf("/orgs/%v/issues?state=%v&filter=all&page_size=100", state.Organization, status))
}
func getTeamIssues(status string) []Issue {
	var result []Issue
	repos := getTeamRepos()
	res := make(chan []Issue)
	for i, xi := range repos {
		go func(i int, xi string) {
			res <- getRepoIssues(xi, state.Organization, status)
		}(i, xi)
	}
	for i := 0; i < len(repos); i++ {
		result = append(result, <-res...)
	}
	return result
}
func getRepoIssues(repo string, owner string, status string) []Issue {
	return getIssues(fmt.Sprintf("/repos/%v/%v/issues?state=%v&filter=all&page_size=100", repo, owner, status))
}
func getIssues(route string) []Issue {
	res, err := getFromGitHub(state.Username, state.Token, route)
	if err != nil {
		log.Fatal(err)
	}
	var issues []Issue
	jsonErr := json.Unmarshal(res, &issues)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return issues
}
func issuesToTable(issues []Issue) [][]string {
	var result [][]string
	for _, i := range issues {
		result = append(result, []string{
			strconv.Itoa(i.ID),
			stringElipse(i.Title, 20), //todo
			i.User.Name,
			i.Assignee.Name,
			i.State,
			i.CreatedAt.Format("2006-01-02 15:04"),
		})
	}
	return result
}
