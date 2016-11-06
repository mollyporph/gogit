package main

import (
	"encoding/json"
	"fmt"
	"log"
)

//RepoID json struct
type RepoID struct {
	ID string `json:"name"`
}

func getOrgRepos() []string {
	repos, err := getFromGitHub(state.Username, state.Token, fmt.Sprintf("/orgs/%v/repos", state.Organization))
	if err != nil {
		log.Fatal(err)
	}
	return parseRepos(repos)
}
func getTeamRepos() []string {
	repos, err := getFromGitHub(state.Username, state.Token, fmt.Sprintf("/teams/%v/repos", state.TeamID))
	if err != nil {
		log.Fatal(err)
	}
	return parseRepos(repos)
}

func getPersonalRepos() []string {
	repos, err := getFromGitHub(state.Username, state.Token, "/user/repos?affiliation=owner")
	if err != nil {
		log.Fatal(err)
	}
	return parseRepos(repos)
}
func parseRepos(repoJSON []byte) []string {
	var repoIDs []RepoID
	jsonErr := json.Unmarshal(repoJSON, &repoIDs)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return repoIDSsTostrings(repoIDs)
}
func repoIDSsTostrings(repoIDs []RepoID) []string {
	var result []string
	for _, i := range repoIDs {
		result = append(result, i.ID)
	}
	return result
}