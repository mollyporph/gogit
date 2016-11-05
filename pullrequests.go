package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"
)

type empty struct{}

//Pullrequest json struct
type Pullrequest struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"` //unmarshaller produces error with empty or nil time.Times
	ClosedAt  string    `json:"closed_at"`
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
func PrintPullRequests(status string, orderby string) {
	pullrequests := getPullrequests(status, orderby)
	PrintTable(
		[]string{"Repo", "Name", "Requester", "Assignee", "From", "To", "State", "Created"},
		pullrequestsToTable(pullrequests))
}

func getPullrequests(status string, orderby string) []Pullrequest {
	prCreatedAt := func(p1, p2 *Pullrequest) bool {
		return p1.CreatedAt.Before(p2.CreatedAt)
	}
	prState := func(p1, p2 *Pullrequest) bool {
		return p1.State < p2.State
	}
	prRepo := func(p1, p2 *Pullrequest) bool {
		return p1.Base.Repo.Name < p2.Base.Repo.Name
	}
	var result []Pullrequest
	switch state.Context {
	case "org":
		result = getOrgPullrequests(status)
	case "team":
		result = getTeamPullrequests(status)
	case "personal":
		result = getPersonalPullrequests(status)
	}
	switch orderby {
	case "created":
		By(prCreatedAt).Sort(result)
	case "state":
		By(prState).Sort(result)
	case "repo":
		By(prRepo).Sort(result)
	default:
	}
	return result
}

func getOrgPullrequests(status string) []Pullrequest {
	repos := getOrgRepos()
	var result []Pullrequest
	sem := make(chan empty, len(repos))
	for i, xi := range repos {
		go func(i int, xi string) {
			repoList := <-getRepoPullRequestsFuture(state.Organization, xi, status)
			result = append(result, repoList...)
			sem <- empty{}
		}(i, xi)
	}
	for i := 0; i < len(repos); i++ {
		<-sem
	}
	return result
}

func getTeamPullrequests(status string) []Pullrequest {
	repos := getTeamRepos()
	var result []Pullrequest
	sem := make(chan empty, len(repos))
	for i, xi := range repos {
		go func(i int, xi string) {
			repoList := <-getRepoPullRequestsFuture(state.Organization, xi, status)
			result = append(result, repoList...)
			sem <- empty{}
		}(i, xi)
	}
	for i := 0; i < len(repos); i++ {
		<-sem
	}
	return result
}

func getPersonalPullrequests(status string) []Pullrequest {
	repos := getPersonalRepos()
	var result []Pullrequest
	sem := make(chan empty, len(repos))
	for i, xi := range repos {
		go func(i int, xi string) {
			repoList := <-getRepoPullRequestsFuture(state.Username, xi, status)
			result = append(result, repoList...)
			sem <- empty{}
		}(i, xi)
	}
	for i := 0; i < len(repos); i++ {
		<-sem
	}
	return result
}

func getRepoPullRequests(owner string, repo string, status string) []Pullrequest {
	resp, err := getFromGitHub(state.Username, state.Token, fmt.Sprintf("/repos/%v/%v/pulls?state=%v", owner, repo, status))
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
			stringElipse(i.Base.Repo.Name, 20),
			stringElipse(i.Title, 20),
			i.User.Name,
			i.Assignee.Name,
			stringElipse(i.Head.Branch, 15),
			stringElipse(i.Base.Branch, 15),
			i.State,
			i.CreatedAt.Format("2006-01-02 15:04")})
	}
	return result
}

func getRepoPullRequestsFuture(owner string, repo string, status string) chan []Pullrequest {
	future := make(chan []Pullrequest)
	go func() { future <- getRepoPullRequests(owner, repo, status) }()
	return future
}

func stringElipse(word string, maxLength int) string {
	result := word[:Min(maxLength, len(word))]
	if len(result) == maxLength {
		result = result + "..."
	}
	return result
}

//Min for ints
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

//Max for ints
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

//By sorting by
type By func(p1, p2 *Pullrequest) bool
type pullrequestSorter struct {
	pullrequests []Pullrequest
	by           func(p1, p2 *Pullrequest) bool
}

func (s *pullrequestSorter) Len() int {
	return len(s.pullrequests)
}
func (s *pullrequestSorter) Swap(i, j int) {
	s.pullrequests[i], s.pullrequests[j] = s.pullrequests[j], s.pullrequests[i]
}
func (s *pullrequestSorter) Less(i, j int) bool {
	return s.by(&s.pullrequests[i], &s.pullrequests[j])
}

//Sort sorts pullrequests
func (by By) Sort(pullrequests []Pullrequest) {
	ps := &pullrequestSorter{
		pullrequests: pullrequests,
		by:           by,
	}
	sort.Sort(ps)
}

//SortDesc sorts descending
func (by By) SortDesc(pullrequests []Pullrequest) {
	ps := &pullrequestSorter{
		pullrequests: pullrequests,
		by:           by,
	}
	sort.Sort(sort.Reverse(ps))
}
