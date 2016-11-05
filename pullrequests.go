package main

import "fmt"
import "log"
import "encoding/json"
import "strconv"

//Pullrequest json struct
type Pullrequest struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	State     string `json:"state"`
	CreatedAt string `json:"created_at"` //unmarshaller produces error with empty or nil time.Times
	ClosedAt  string `json:"closed_at"`
	Base      struct {
		Branch string `json:"ref"`
		Repo   struct {
			Name string `json:"name"`
		} `json:"repo"`
	} `json:"base"`
	User struct {
		Name string `json:"login"`
	} `json:"user"`
	Head struct {
		Branch string `json:"ref"`
	} `json:"head"`
	Assignee struct {
		Name string `json:"login"`
	} `json:"assignee"`
}

//PrintPullRequests prints the pull requests that are in `status` status as a table
func PrintPullRequests(status string) {
	pullrequests := getPullrequests(status)
	PrintTable(
		[]string{"Repo", "Id", "Title", "Requester", "Assignee", "From", "To", "State", "Created-date"},
		pullrequestsToTable(pullrequests))
}
func getPullrequests(status string) []Pullrequest {
	switch state.Context {
	case "org":
		return getOrgPullrequests(status)
	case "team":
		return getTeamPullrequests(status)
	case "personal":
		return getPersonalPullrequests(status)
	}
	return nil
}
func getOrgPullrequests(status string) []Pullrequest {
	repos := getOrgRepos()
	var result []Pullrequest
	for _, i := range repos {
		result = append(result, getRepoPullRequests(state.Organization, i, status)...)
	}
	return result
}
func getTeamPullrequests(status string) []Pullrequest {
	repos := getTeamRepos()
	var result []Pullrequest
	for _, i := range repos {
		result = append(result, getRepoPullRequests(state.Organization, i, status)...)
	}
	return result
}
func getPersonalPullrequests(status string) []Pullrequest {
	repos := getPersonalRepos()
	var result []Pullrequest
	for _, i := range repos {
		result = append(result, getRepoPullRequests(state.Username, i, status)...)
	}
	return result
}
func getRepoPullRequests(owner string, repo string, status string) []Pullrequest {
	resp, err := getFromGitHub(state.Username, state.Token, fmt.Sprintf("/repos/%v/%v/pulls?status=%v", owner, repo, status))
	if err != nil {
		log.Fatal(err)
	}
	var result []Pullrequest
	jsonErr := json.Unmarshal(resp, &result)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return result
}
func pullrequestsToTable(pullrequests []Pullrequest) [][]string {
	var result [][]string
	for _, i := range pullrequests {
		result = append(result, []string{
			i.Base.Repo.Name,
			strconv.Itoa(i.ID),
			i.Title,
			i.User.Name,
			i.Assignee.Name,
			i.Head.Branch,
			i.Base.Branch,
			i.State,
			i.CreatedAt})
	}
	return result
}
