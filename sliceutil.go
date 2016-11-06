package main

import "strings"

//PivotStringSlice returns an array pivoted from horizontal to vertical
func PivotStringSlice(data []string) [][]string {
	var result [][]string
	for _, s := range data {
		result = append(result, []string{strings.ToLower(s)})
	}
	return result
}

func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
