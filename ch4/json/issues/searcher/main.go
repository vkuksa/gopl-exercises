// Exercise 4.10 from gopl.io

package searcher

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"gopl-exercises/ch4/json/issues/github"
)

const (
	month = 731
	year  = 8766
)

func filterResultsByTimeframe(issues *github.IssuesSearchResult, hours int) *github.IssuesSearchResult {
	now := time.Now()

	filtered := &github.IssuesSearchResult{TotalCount: 0, Items: nil}
	for _, issue := range issues.Items {
		difference := now.Sub(issue.CreatedAt)
		if int(difference.Hours()) < hours {
			filtered.Items = append(filtered.Items, issue)
		}
	}
	filtered.TotalCount = len(filtered.Items)

	return filtered
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

	monthly := filterResultsByTimeframe(result, month)
	yearly := filterResultsByTimeframe(result, year)
	allTime := filterResultsByTimeframe(result, math.MaxInt)

	fmt.Println("Monthly")
	printResults(monthly)

	fmt.Println("Yearly")
	printResults(yearly)

	fmt.Println("All-time")
	printResults(allTime)

	fmt.Print()
}
