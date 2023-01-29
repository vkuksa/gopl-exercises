package github

import "time"

const (
	SearchIssuesURL  = "https://api.github.com/search/issues"
	AccountIssuesURL = "https://api.github.com/issues"
	RepoIssueURL     = "https://api.github.com/repos/%s/%s/issues/%d"
	RepoIssuesURL    = "https://api.github.com/repos/%s/%s/issues"
)

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

type Issue struct {
	Id            int
	Number        int
	HTMLURL       string `json:"html_url"`
	Title         string
	State         string
	User          *User
	CreatedAt     time.Time `json:"created_at"`
	Body          string    // in Markdown format
	RepositoryUrl string    `json:"repository_url"`
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}
