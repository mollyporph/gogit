package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// PostToGithub http POST to github and returns string bytearray
func PostToGithub(usr string, pwOrToken string, route string, body string) ([]byte, error) {
	resp, err := postToGithub(usr, pwOrToken, route, body)
	if err != nil {
		if gerr, ok := err.(*githubError); ok {
			if gerr.statuscode == 401 {
				mfa := getMfa()
				resp, err = postToGithubMFA(usr, pwOrToken, route, body, mfa)
			}
		}
	}
	return resp, err
}
func postToGithubMFA(usr string, pwOrToken string, route string, body string, mfaToken string) ([]byte, error) {
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

func postToGithub(usr string, pwOrToken string, route string, body string) ([]byte, error) {
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
	baseURL := "https://api.github.com"
	fullURL := baseURL + route
	return http.NewRequest("POST", fullURL, bytes.NewBufferString(body))
}
func handleResponse(resp *http.Response) ([]byte, error) {
	if resp.StatusCode == 401 {
		gerr := &githubError{}
		gerr.statuscode = 401
		gerr.msg = "authentication failed"
		opt := resp.Header.Get("X-GitHub-OTP")
		if opt != "" && strings.HasPrefix(opt, "required") {
			gerr.githubAuthStatus = MfaRequired
		} else {
			gerr.githubAuthStatus = WrongPassword
		}
		return nil, gerr
	} else if resp.StatusCode == 404 {
		gerr := &githubError{}
		gerr.statuscode = 404
		gerr.msg = "Not found"
		return nil, gerr
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()
	return respBody, nil
}
