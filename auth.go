package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

//GitHubAuthStatus int value of status
type GitHubAuthStatus int

//Github error statuses
const (
	MfaRequired GitHubAuthStatus = 1 + iota
	WrongPassword
	NotAuthorized
)

//GitHubError custom error with statuscode and authstatus
type GitHubError struct {
	Msg              string
	StatusCode       int
	GitHubAuthStatus GitHubAuthStatus
}

func (e *GitHubError) Error() string { return e.Msg }
func getPatPermnissions() []string {
	return []string{"repo", "read:org", "user"}
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

// GetGithubPatAndUsername returns the PAT token generated and the correct username
func GetGithubPatAndUsername() (string, string) {
	patPermissions := getPatPermnissions()
	askForPatConfirmation(patPermissions)
	fmt.Println("Your password will only be used to create a PAT (over https) and will not be stored anywhere")

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	requestBody := TokenRequestBody{
		Scopes: patPermissions,
		Note:   "GoGit PAT " + hostname, //todo: make changeable}
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatal(err)
	}
	pat, username := getGithubPatAndUsername(string(body))
	return pat, username
}
func getGithubPatAndUsername(body string) (string, string) {
	username := getUsername()
	password := getPassword()
	tokenJSON, err := PostToGitHub(username, password, "/authorizations", body)
	if err != nil {
		if gerr, ok := err.(*GitHubError); ok {
			if gerr.StatusCode == 401 {
				if askForConfirmation("Wrong username, password or mfa, would you like to try again?") {
					return getGithubPatAndUsername(body)
				}
			}
		}
		log.Fatal(err)
	}
	var token TokenResponseBody
	jsonErr := json.Unmarshal(tokenJSON, &token)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return token.Token, username
}
func askForPatConfirmation(patPermissions []string) {
	fmt.Println("GoGit needs to create a personal access token to be able to access github's API.")
	fmt.Println(fmt.Sprintf("We will create a PAT with the following permissions: %v. for more info on oauth scopes, visit https://developer.github.com/v3/oauth/#scopes", patPermissions))
	fmt.Println("The token will be saved in your .gogit config file")
	fmt.Println("To create a PAT gogit will need your github username and password, and a multifactor authentication token if you have it enabled.")
	if !askForConfirmation("Create PAT token?") {
		fmt.Println("Cannot continue without correct config")
		os.Exit(1)
	}
}
