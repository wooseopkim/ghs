package analyzer

import (
	"github.com/wooseopkim/ghs/internal/counter"
	"github.com/wooseopkim/ghs/internal/github"
	"github.com/wooseopkim/ghs/internal/stats"
)

func AnalyzePullRequests(prs []github.PullRequest) stats.PullRequest {
	s := stats.PullRequest{
		Additions:    counter.New[uint](),
		Deletions:    counter.New[uint](),
		ChangedFiles: counter.New[uint](),
		Commits:      counter.New[uint](),
		TitleLength:  counter.New[uint](),
		BodyLength:   counter.New[uint](),
	}

	for _, v := range prs {
		s.Additions.Increment(v.Id, v.Additions)
		s.Deletions.Increment(v.Id, v.Deletions)
		s.ChangedFiles.Increment(v.Id, v.ChangedFiles)
		s.Commits.Increment(v.Id, v.Commits)
		s.TitleLength.Increment(v.Id, uint(len(v.Title)))
		s.BodyLength.Increment(v.Id, uint(len(v.Body)))

		if v.Merged {
			s.Merged += 1
		}
		if v.Closed {
			s.Closed += 1
		}
		if v.OwnRepository {
			s.OwnRepository += 1
		}
	}

	s.Found = uint(len(prs))

	return s
}
