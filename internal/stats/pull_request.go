package stats

import "github.com/wooseopkim/ghs/internal/counter"

type PullRequest struct {
	Additions     *counter.Counter[uint]
	Deletions     *counter.Counter[uint]
	ChangedFiles  *counter.Counter[uint]
	Commits       *counter.Counter[uint]
	Merged        uint
	Closed        uint
	OwnRepository uint
	TitleLength   *counter.Counter[uint]
	BodyLength    *counter.Counter[uint]
	Found         uint
}
