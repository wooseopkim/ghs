package stats

import "github.com/wooseopkim/ghs/internal/counter"

type Repository struct {
	Owners    *counter.Counter[uint]
	Languages *counter.Counter[float32]
	Stars     *counter.Counter[uint]
	Private   uint
	Public    uint
	Found     uint
}
