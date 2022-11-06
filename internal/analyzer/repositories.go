package analyzer

import (
	"github.com/wooseopkim/ghs/internal/counter"
	"github.com/wooseopkim/ghs/internal/github"
	"github.com/wooseopkim/ghs/internal/stats"
)

func AnalyzeRepositories(repos []github.Repository) stats.Repository {
	s := stats.Repository{
		Owners:    counter.New[uint](),
		Languages: counter.New[float32](),
		Stars:     counter.New[uint](),
	}

	for _, v := range repos {
		s.Owners.Increment(v.Owner, uint(1))

		languagesTotal := uint(0)
		for _, vv := range v.Languages {
			languagesTotal += vv.Size
		}
		for _, vv := range v.Languages {
			score := float32(vv.Size) / float32(languagesTotal)
			s.Languages.Increment(vv.Name, score)
		}

		s.Stars.Increment(v.Id, v.Stars)

		if v.Private {
			s.Private += 1
		} else {
			s.Public += 1
		}
	}

	s.Found = uint(len(repos))

	return s
}
