package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//GetFromGitHub http GET to GitHub to GitHub and returns json string bytearray
func GetFromGitHub(usr string, pwOrToken string, route string) ([]byte, error) {
	resp, err := getFromGitHub(usr, pwOrToken, route)
	if err != nil {
		if gerr, ok := err.(*GitHubError); ok {
			if gerr.StatusCode == 401 {
				mfa := getMfa()
				resp, err = getFromGitHubMFA(usr, pwOrToken, route, mfa)
			}
		}
	}
	return resp, err
}
func getFromGitHub(usr string, pwOrToken string, route string) ([]byte, error) {
	req, err := setupGetRequest(route)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(usr, pwOrToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return handleResponse(resp)
}
func getFromGitHubMFA(usr string, pwOrToken string, route string, mfaToken string) ([]byte, error) {
	req, err := setupGetRequest(route)
	if err != nil {
		log.Fatal(err)
	}
	if mfaToken == "" {
		log.Fatalln("MFA empty")
	}
	req.Header.Add("X-GitHub-OTP", mfaToken)
	req.SetBasicAuth(usr, pwOrToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return handleResponse(resp)
}

// PostToGitHub http POST to GitHub and returns json string bytearray
func PostToGitHub(usr string, pwOrToken string, route string, body string) ([]byte, error) {
	resp, err := postToGitHub(usr, pwOrToken, route, body)
	if err != nil {
		if gerr, ok := err.(*GitHubError); ok {
			if gerr.StatusCode == 401 {
				mfa := getMfa()
				resp, err = postToGitHubMFA(usr, pwOrToken, route, body, mfa)
			}
		}
	}
	return resp, err
}
func postToGitHubMFA(usr string, pwOrToken string, route string, body string, mfaToken string) ([]byte, error) {
	req, err := setupPostRequest(route, body)
	if err != nil {
		log.Fatal(err)
	}
	if mfaToken == "" {
		log.Fatalln("MFA empty")
	}
	req.Header.Add("X-GitHub-OTP", mfaToken)
	req.SetBasicAuth(usr, pwOrToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return handleResponse(resp)
}

func postToGitHub(usr string, pwOrToken string, route string, body string) ([]byte, error) {
	req, err := setupPostRequest(route, body)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(usr, pwOrToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return handleResponse(resp)

}
func setupPostRequest(route string, body string) (*http.Request, error) {
	baseURL := "https://api.GitHub.com"
	fullURL := baseURL + route
	return http.NewRequest("POST", fullURL, bytes.NewBufferString(body))
}
func setupGetRequest(route string) (*http.Request, error) {
	baseURL := "https://api.GitHub.com"
	fullURL := baseURL + route
	return http.NewRequest("GET", fullURL, nil)
}
func handleResponse(resp *http.Response) ([]byte, error) {
	if resp.StatusCode == 401 {
		gerr := &GitHubError{}
		gerr.StatusCode = 401
		gerr.Msg = "authentication failed"
		opt := resp.Header.Get("X-GitHub-OTP")
		if opt != "" && strings.HasPrefix(opt, "required") {
			gerr.GitHubAuthStatus = MfaRequired
		} else {
			gerr.GitHubAuthStatus = WrongPassword
		}
		return nil, gerr
	} else if resp.StatusCode == 404 {
		gerr := &GitHubError{}
		gerr.StatusCode = 404
		gerr.Msg = "Not found"
		return nil, gerr
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()
	return respBody, nil
}
