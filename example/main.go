package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/wooseopkim/ghs"
	"github.com/wooseopkim/ghs/internal/stats"
	"golang.org/x/term"
)

func main() {
	fmt.Print("GitHub token: ")
	token, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println()

	if len(strings.TrimSpace(string(token))) == 0 {
		fmt.Println("You should provide GitHub token")
		os.Exit(1)
	}

	stats, err := ghs.GetStats(string(token))
	if err != nil {
		log.Fatalln(err)
	}

	print(stats)
}

func print(stats stats.Records) {
	var b strings.Builder

	b.WriteString("[REPOSITORIES]")
	b.WriteByte('\n')

	b.WriteString("TOP OWNERS:")
	b.WriteByte('\n')
	for _, v := range stats.Repository.Owners.Top(5) {
		b.WriteString("  ")
		b.WriteString(v.Key)
		b.WriteByte(' ')
		b.WriteByte('(')
		b.WriteString(strconv.Itoa(int(v.Count)))
		b.WriteByte(')')
		b.WriteByte('\n')
	}

	b.WriteString("TOP LANGUAGES:")
	b.WriteByte('\n')
	for _, v := range stats.Repository.Languages.Top(5) {
		b.WriteString("  ")
		b.WriteString(v.Key)
		b.WriteByte(' ')
		b.WriteByte('(')
		b.WriteString(fmt.Sprintf("%.2f", (v.Count)))
		b.WriteByte(')')
		b.WriteByte('\n')
	}

	b.WriteString("TOP STARS:")
	b.WriteByte('\n')
	for _, v := range stats.Repository.Stars.Top(5) {
		b.WriteString("  ")
		b.WriteString(v.Key)
		b.WriteByte(' ')
		b.WriteByte('(')
		b.WriteString(strconv.Itoa(int(v.Count)))
		b.WriteByte(')')
		b.WriteByte('\n')
	}

	b.WriteString("FOUND:")
	b.WriteByte(' ')
	b.WriteString(strconv.Itoa(int(stats.Repository.Found)))
	b.WriteByte('\n')

	b.WriteString("PRIVATE:")
	b.WriteByte(' ')
	b.WriteString(strconv.Itoa(int(stats.Repository.Private)))
	b.WriteByte('\n')

	b.WriteString("PUBLIC:")
	b.WriteByte(' ')
	b.WriteString(strconv.Itoa(int(stats.Repository.Public)))
	b.WriteByte('\n')

	b.WriteString("[COMMITS]")
	b.WriteByte('\n')

	b.WriteString("FOUND:")
	b.WriteByte(' ')
	b.WriteString(strconv.Itoa(int(stats.Commit.Found)))
	b.WriteByte('\n')

	b.WriteString("TOP ADDITIONS:")
	b.WriteByte('\n')
	for _, v := range stats.Commit.Additions.Top(5) {
		b.WriteString("  ")
		b.WriteString(v.Key)
		b.WriteByte(' ')
		b.WriteByte('(')
		b.WriteString(strconv.Itoa(int(v.Count)))
		b.WriteByte(')')
		b.WriteByte('\n')
	}

	b.WriteString("TOP DELETIONS:")
	b.WriteByte('\n')
	for _, v := range stats.Commit.Deletions.Top(5) {
		b.WriteString("  ")
		b.WriteString(v.Key)
		b.WriteByte(' ')
		b.WriteByte('(')
		b.WriteString(strconv.Itoa(int(v.Count)))
		b.WriteByte(')')
		b.WriteByte('\n')
	}

	b.WriteString("AVERAGE ADDITIONS:")
	b.WriteByte(' ')
	b.WriteString(fmt.Sprintf("%.2f", stats.Commit.Additions.Average()))
	b.WriteByte('\n')

	b.WriteString("AVERAGE DELETIONS:")
	b.WriteByte(' ')
	b.WriteString(fmt.Sprintf("%.2f", stats.Commit.Deletions.Average()))
	b.WriteByte('\n')

	b.WriteString("[PULL REQUESTS]")
	b.WriteByte('\n')

	b.WriteString("FOUND:")
	b.WriteByte(' ')
	b.WriteString(strconv.Itoa(int(stats.PullRequest.Found)))
	b.WriteByte('\n')

	b.WriteString("TOP ADDITIONS:")
	b.WriteByte('\n')
	for _, v := range stats.PullRequest.Additions.Top(5) {
		b.WriteString("  ")
		b.WriteString(v.Key)
		b.WriteByte(' ')
		b.WriteByte('(')
		b.WriteString(strconv.Itoa(int(v.Count)))
		b.WriteByte(')')
		b.WriteByte('\n')
	}

	b.WriteString("TOP DELETIONS:")
	b.WriteByte('\n')
	for _, v := range stats.PullRequest.Deletions.Top(5) {
		b.WriteString("  ")
		b.WriteString(v.Key)
		b.WriteByte(' ')
		b.WriteByte('(')
		b.WriteString(strconv.Itoa(int(v.Count)))
		b.WriteByte(')')
		b.WriteByte('\n')
	}

	b.WriteString("TOP CHANGED FILES:")
	b.WriteByte('\n')
	for _, v := range stats.PullRequest.ChangedFiles.Top(5) {
		b.WriteString("  ")
		b.WriteString(v.Key)
		b.WriteByte(' ')
		b.WriteByte('(')
		b.WriteString(strconv.Itoa(int(v.Count)))
		b.WriteByte(')')
		b.WriteByte('\n')
	}

	b.WriteString("TOP COMMITS:")
	b.WriteByte('\n')
	for _, v := range stats.PullRequest.Commits.Top(5) {
		b.WriteString("  ")
		b.WriteString(v.Key)
		b.WriteByte(' ')
		b.WriteByte('(')
		b.WriteString(strconv.Itoa(int(v.Count)))
		b.WriteByte(')')
		b.WriteByte('\n')
	}

	b.WriteString("AVERAGE ADDITIONS:")
	b.WriteByte(' ')
	b.WriteString(fmt.Sprintf("%.2f", stats.PullRequest.Additions.Average()))
	b.WriteByte('\n')

	b.WriteString("AVERAGE DELETIONS:")
	b.WriteByte(' ')
	b.WriteString(fmt.Sprintf("%.2f", stats.PullRequest.Deletions.Average()))
	b.WriteByte('\n')

	b.WriteString("AVERAGE CHANGED FILES:")
	b.WriteByte(' ')
	b.WriteString(fmt.Sprintf("%.2f", stats.PullRequest.ChangedFiles.Average()))
	b.WriteByte('\n')

	b.WriteString("AVERAGE COMMITS:")
	b.WriteByte(' ')
	b.WriteString(fmt.Sprintf("%.2f", stats.PullRequest.Commits.Average()))
	b.WriteByte('\n')

	b.WriteString("AVERAGE TITLE LENGTH:")
	b.WriteByte(' ')
	b.WriteString(fmt.Sprintf("%.2f", stats.PullRequest.TitleLength.Average()))
	b.WriteByte('\n')

	b.WriteString("AVERAGE BODY LENGTH:")
	b.WriteByte(' ')
	b.WriteString(fmt.Sprintf("%.2f", stats.PullRequest.BodyLength.Average()))
	b.WriteByte('\n')

	b.WriteString("MERGED:")
	b.WriteByte(' ')
	b.WriteString(strconv.Itoa(int(stats.PullRequest.Merged)))
	b.WriteByte('\n')

	b.WriteString("CLOSED:")
	b.WriteByte(' ')
	b.WriteString(strconv.Itoa(int(stats.PullRequest.Closed)))
	b.WriteByte('\n')

	b.WriteString("OWN REPOSITORY:")
	b.WriteByte(' ')
	b.WriteString(strconv.Itoa(int(stats.PullRequest.OwnRepository)))
	b.WriteByte('\n')

	fmt.Println(b.String())
}
