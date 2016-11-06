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
	Name      string    `json:"title"`
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

	s := GetSortFuncs()
	var result []Pullrequest
	switch state.Context {
	case "org":
		result = getOrgPullrequests(status)
	case "team":
		result = getTeamPullrequests(status)
	case "personal":
		result = getPersonalPullrequests(status)
	}
	if orderby != "" && stringInSlice(orderby, []string{"repo", "name", "requester", "assignee", "from", "to", "state", "created"}) {
		By(s[orderby]).Sort(result)
	}

	return result
}

func getTeamPullrequests(status string) []Pullrequest {
	repos := getTeamRepos()
	return getPullrequestsRoutine(repos, state.Organization, status)
}
func getPersonalPullrequests(status string) []Pullrequest {
	repos := getPersonalRepos()
	return getPullrequestsRoutine(repos, state.Username, status)
}
func getOrgPullrequests(status string) []Pullrequest {
	repos := getOrgRepos()
	return getPullrequestsRoutine(repos, state.Organization, status)
}
func getPullrequestsRoutine(repos []string, owner string, status string) []Pullrequest {
	var result []Pullrequest
	res := make(chan []Pullrequest)
	for i, xi := range repos {
		go func(i int, xi string) {
			res <- getRepoPullRequests(owner, xi, status)
		}(i, xi)
	}
	for i := 0; i < len(repos); i++ {
		result = append(result, <-res...)
	}
	return result
}
func getRepoPullRequestsFuture(owner string, repo string, status string) chan []Pullrequest {
	future := make(chan []Pullrequest)
	go func() { future <- getRepoPullRequests(owner, repo, status) }()
	return future
}
func getRepoPullRequests(owner string, repo string, status string) []Pullrequest {
	resp, err := getFromGitHub(state.Username, state.Token, fmt.Sprintf("/repos/%v/%v/pulls?state=%v&page_size=100", owner, repo, status))
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
			stringElipse(i.Name, 20),
			i.User.Name,
			i.Assignee.Name,
			stringElipse(i.Head.Branch, 15),
			stringElipse(i.Base.Branch, 15),
			i.State,
			i.CreatedAt.Format("2006-01-02 15:04")})
	}
	return result
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

//GetSortFuncs returns the less functions required by Sort By
func GetSortFuncs() map[string]func(*Pullrequest, *Pullrequest) bool {
	s := make(map[string]func(*Pullrequest, *Pullrequest) bool)
	s["repo"] = func(p1, p2 *Pullrequest) bool {
		return p1.Base.Repo.Name < p2.Base.Repo.Name
	}
	s["name"] = func(p1, p2 *Pullrequest) bool {
		return p1.Name < p2.Name
	}
	s["requester"] = func(p1, p2 *Pullrequest) bool {
		return p1.User.Name < p2.User.Name
	}
	s["assignee"] = func(p1, p2 *Pullrequest) bool {
		return p1.Assignee.Name < p2.Assignee.Name
	}
	s["from"] = func(p1, p2 *Pullrequest) bool {
		return p1.Head.Branch < p2.Head.Branch
	}
	s["to"] = func(p1, p2 *Pullrequest) bool {
		return p1.Base.Branch < p2.Base.Branch
	}
	s["state"] = func(p1, p2 *Pullrequest) bool {
		return p1.State < p2.State
	}
	s["created"] = func(p1, p2 *Pullrequest) bool {
		return p1.CreatedAt.Before(p2.CreatedAt)
	}
	return s
}
