package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func makeRequest(method, url, token string, body []byte) (*http.Response, error) {
	reader := bytes.NewReader(body)
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, fmt.Errorf("Request creation failed")
	}
	if len(token) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	req.Header.Set("Accept", "application/vnd.github.v3.text-match+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request performing failed")
	}

	return resp, nil
}

func performGetRequest(url, token string) (*http.Response, error) {
	resp, _ := makeRequest("GET", url, token, nil)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("Request execution returned non OK value: %s\n", resp.Status)
	}

	return resp, nil
}

func performPatchRequest(url, token string, body []byte) (*http.Response, error) {
	resp, _ := makeRequest("PATCH", url, token, body)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("Request execution returned non Ok value: %s\n", resp.Status)
	}

	return resp, nil
}

func performPostRequest(url, token string, body []byte) (*http.Response, error) {
	resp, _ := makeRequest("POST", url, token, body)
	if resp.StatusCode != http.StatusCreated {
		resp.Body.Close()
		return nil, fmt.Errorf("Request execution returned non Created value: %s\n", resp.Status)
	}

	return resp, nil
}

func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := performGetRequest(SearchIssuesURL+"?q="+q, "")
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed issues search: " + err.Error())
	}

	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}

func CreateIssue(token, owner, repository string, fields map[string]string) (*Issue, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(fields)
	if err != nil {
		return nil, fmt.Errorf("Encoding of issue creation fields failed: " + err.Error())
	}

	url := fmt.Sprintf(RepoIssuesURL, owner, repository)
	resp, err := performPostRequest(url, token, buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("Issue creation " + url + " failed: " + err.Error())
	}

	var result *Issue
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}

	resp.Body.Close()
	return result, nil
}

func UpdateIssue(token, owner, repository string, number uint64, fields map[string]string) (*Issue, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(fields)
	if err != nil {
		return nil, fmt.Errorf("Encoding of issue update fields failed: " + err.Error())
	}

	url := fmt.Sprintf(RepoIssueURL, owner, repository, number)
	resp, err := performPatchRequest(url, token, buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("Issue update " + url + " failed: " + err.Error())
	}

	var result *Issue
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}

	resp.Body.Close()
	return result, nil
}

func GetIssue(token, owner, repository string, number uint64) (*Issue, error) {
	url := fmt.Sprintf(RepoIssueURL, owner, repository, number)
	resp, err := performGetRequest(url, token)
	if err != nil {
		return nil, fmt.Errorf("Retrieving issue " + url + " failed: " + err.Error())
	}

	var result *Issue
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}

	resp.Body.Close()
	return result, nil
}

func SearchIssuesOf(token string) ([]Issue, error) {
	resp, err := performGetRequest(AccountIssuesURL, token)
	if err != nil {
		return nil, fmt.Errorf("Failed User issues search: " + err.Error())
	}

	var result []Issue = make([]Issue, 0) //! We initialize with make, so if returned JSON has empty array - the value of result would be [], but not nil
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}

	resp.Body.Close()
	return result, nil
}
