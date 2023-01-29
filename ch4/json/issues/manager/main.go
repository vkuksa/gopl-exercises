// Exercise 4.11 from gopl.io

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"

	"gopl-exercises/ch4/json/issues/github"
)

const (
	usage = `usage: 
		create [token_id] <owner> <repository_name>
		get [token_id] <owner> <repository_name> <issue_number>
		update [token_id] <owner> <repository_name> <issue_number>
		search [token_id]`
)

func editIssueWithEditor(issue *github.Issue) (map[string]string, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "gedit"
	}
	editorPath, err := exec.LookPath(editor)
	if err != nil {
		return nil, err
	}
	tempfile, err := os.CreateTemp("", "issue_")
	if err != nil {
		return nil, err
	}
	defer tempfile.Close()
	defer os.Remove(tempfile.Name())

	input := map[string]string{
		"title": "",
		"state": "",
		"body":  "",
	}
	if issue != nil {
		input["title"] = issue.Title
		input["state"] = issue.State
		input["body"] = issue.Body
	}

	err = json.NewEncoder(tempfile).Encode(input)
	if err != nil {
		return nil, err
	}

	cmd := &exec.Cmd{
		Path:   editorPath,
		Args:   []string{editor, tempfile.Name()},
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	fields := make(map[string]string)
	content, err := os.ReadFile(tempfile.Name())
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, &fields)

	return fields, err
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println(usage)
		os.Exit(1)
	}

	var token string = os.Args[2]

	switch os.Args[1] {
	case "create":
		{
			if len(os.Args) < 5 {
				fmt.Println(usage)
				os.Exit(0)
			}

			fields, err := editIssueWithEditor(nil)
			if err != nil {
				log.Fatal("Retrieving issue details failed: ", err)
			}

			issue, err := github.CreateIssue(token, os.Args[3], os.Args[4], fields)
			if err != nil || issue == nil {
				log.Fatal("Retrieving issue details failed: ", err)
			}

			fmt.Printf("Created issue: #%-5d\t%d\t%9.9s\t%.55s\n",
				issue.Id, issue.Number, issue.User.Login, issue.Title)
		}
	case "get":
		{
			if len(os.Args) < 6 {
				fmt.Println(usage)
				os.Exit(0)
			}

			number, err := strconv.ParseUint(os.Args[5], 10, 64)
			if err != nil {
				log.Fatal("issue id should be a numeric value")
			}

			issue, err := github.GetIssue(token, os.Args[3], os.Args[4], number)
			if err != nil {
				log.Fatal("Retrieving issue details failed: ", err)
			}

			fmt.Printf("Retrieved issue: #%-5d\t%d\t%9.9s\t%.55s\n",
				issue.Id, issue.Number, issue.User.Login, issue.Title)
		}
	case "update":
		{
			if len(os.Args) < 6 {
				fmt.Println(usage)
				os.Exit(0)
			}

			number, err := strconv.ParseUint(os.Args[5], 10, 64)
			if err != nil {
				log.Fatal("issue id should be a numeric value")
			}

			issue, err := github.GetIssue(token, os.Args[3], os.Args[4], number)
			if err != nil {
				log.Fatal("Retrieving issue details failed: ", err) //Create new here?
			}

			fields, err := editIssueWithEditor(issue)
			if err != nil {
				log.Fatal("Retrieving issue details failed: ", err)
			}

			issue, err = github.UpdateIssue(token, os.Args[3], os.Args[4], number, fields)
			if err != nil || issue == nil {
				log.Fatal("Retrieving issue details failed: ", err)
			}

			fmt.Printf("Updated issue: #%-5d\t%d\t%9.9s\t%.55s\n",
				issue.Id, issue.Number, issue.User.Login, issue.Title)
		}
	case "search":
		{
			issues, err := github.SearchIssuesOf(token)
			if err != nil {
				log.Fatal("Retrieving issues failed: ", err)
			}

			for _, issue := range issues {
				fmt.Printf("%s\t%s\t#%d\n",
					issue.User.Login, path.Base(issue.RepositoryUrl), issue.Number)
			}
		}
	default:
		fmt.Println(usage)
	}
}
