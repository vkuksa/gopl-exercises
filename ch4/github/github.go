package github

import "time"

const (
	GithubApiURL     = "https://api.github.com"
	SearchIssuesURL  = GithubApiURL + "/search/issues"
	AccountIssuesURL = GithubApiURL + "/issues"
	RepoIssueURL     = GithubApiURL + "/repos/%s/%s/issues/%d"
	RepoIssuesURL    = GithubApiURL + "/repos/%s/%s/issues"
	//label:bug
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
