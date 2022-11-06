package github

import "time"

type Repository struct {
	Owner     string
	Name      string
	Id        string
	Languages []Language
	Stars     uint
	Private   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
