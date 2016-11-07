package main

type Issue struct{}

func PrintIssues() {
}

func GetIssues() []Issue {
	switch state.Context {
	case "org":
		return getOrgIssues()
	case "team":
		return getTeamIssues()
	case "personal":
		return getMyIssues()
	}
	return nil
}

func getMyIssues() []Issue {
	return nil
}
func getOrgIssues() []Issue {
	return nil
}
func getTeamIssues() []Issue {
	return nil
}
