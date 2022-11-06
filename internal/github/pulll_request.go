package github

type PullRequest struct {
	Id            string
	Additions     uint
	Deletions     uint
	ChangedFiles  uint
	Commits       uint
	Merged        bool
	Closed        bool
	Title         string
	Body          string
	OwnRepository bool
}
