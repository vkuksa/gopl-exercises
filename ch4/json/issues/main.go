// Exercise 4.10 from gopl.io

package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"gopl-exercises/ch4/json/issues/github" // Exported module from the book
)

const (
	month = 731
	year  = 8766
)

func sortResultsByTimestamp(issues *github.IssuesSearchResult, hours int) *github.IssuesSearchResult {
	now := time.Now()

	timestamped := &github.IssuesSearchResult{TotalCount: 0, Items: nil}
	for _, issue := range issues.Items {
		difference := now.Sub(issue.CreatedAt)
		if int(difference.Hours()) < hours {
			timestamped.Items = append(timestamped.Items, issue)
		}
	}
	timestamped.TotalCount = len(timestamped.Items)

	return timestamped
}

func printResults(issues *github.IssuesSearchResult) {
	fmt.Printf("%d issues:\n", issues.TotalCount)
	for _, item := range issues.Items {
		fmt.Printf("#%-5d %9.9s %.55s\n",
			item.Number, item.User.Login, item.Title)
	}
}

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	monthly := sortResultsByTimestamp(result, month)
	yearly := sortResultsByTimestamp(result, year)
	allTime := sortResultsByTimestamp(result, math.MaxInt)

	fmt.Println("Monthly")
	printResults(monthly)

	fmt.Println("Yearly")
	printResults(yearly)

	fmt.Println("All-time")
	printResults(allTime)

	fmt.Print()
}
