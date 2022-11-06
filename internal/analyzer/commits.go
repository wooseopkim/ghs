package analyzer

import (
	"github.com/wooseopkim/ghs/internal/counter"
	"github.com/wooseopkim/ghs/internal/github"
	"github.com/wooseopkim/ghs/internal/stats"
)

func AnalyzeCommits(commits []github.Commit) stats.Commit {
	s := stats.Commit{
		Additions: counter.New[uint](),
		Deletions: counter.New[uint](),
	}

	for _, v := range commits {
		s.Additions.Increment(v.Id, v.Additions)
		s.Deletions.Increment(v.Id, v.Deletions)
	}

	s.Found = uint(len(commits))

	return s
}
