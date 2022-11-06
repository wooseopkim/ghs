package stats

import "github.com/wooseopkim/ghs/internal/counter"

type Commit struct {
	Additions *counter.Counter[uint]
	Deletions *counter.Counter[uint]
	Found     uint
}
