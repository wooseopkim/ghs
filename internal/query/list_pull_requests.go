package query

import (
	"fmt"

	"github.com/machinebox/graphql"
	"github.com/wooseopkim/ghs/internal/github"
)

type ListPullRequestsOptions struct {
	GitHubToken       string
	Username          string
	PullRequestCursor string
}

func (c *client) ListPullRequests(
	options *ListPullRequestsOptions,
) ([]github.PullRequest, error) {
	return c.listPullRequests(
		options,
		[]github.PullRequest{},
	)
}

func (c *client) listPullRequests(
	options *ListPullRequestsOptions,
	prs []github.PullRequest,
) ([]github.PullRequest, error) {
	req := graphql.NewRequest(`
        query ListPullRequests(
            $username: String!
            $pullRequestLimit: Int = 20
            $pullRequestCursor: String
        ) {
            user(login: $username) {
                pullRequests(
                    first: $pullRequestLimit
                    after: $pullRequestCursor
                ) {
                    totalCount
                    pageInfo {
                        startCursor
                        endCursor
                    }
                    nodes {
                        additions
                        deletions
                        changedFiles
						commits(first: 1) {
							totalCount
						}
                        merged
                        closed
                        title
                        body
						number
                        repository {
							nameWithOwner
                            owner {
                                login
                            }
                        }
                    }
                }
            }
        }
    `)
	req.Header.Set("Authorization", "bearer "+options.GitHubToken)
	req.Var("username", options.Username)
	if options.PullRequestCursor != "" {
		req.Var("pullRequestCursor", options.PullRequestCursor)
	}

	var res struct {
		User struct {
			PullRequests struct {
				TotalCount uint
				PageInfo   PageInfo
				Nodes      []struct {
					Additions    uint
					Deletions    uint
					ChangedFiles uint
					Commits      struct {
						TotalCount uint
					}
					Merged     bool
					Closed     bool
					Title      string
					Body       string
					Number     uint
					Repository struct {
						NameWithOwner string
						Owner         struct {
							Login string
						}
					}
				}
			}
		}
	}
	err := c.run(req, &res)
	if err != nil {
		return nil, err
	}

	p := res.User.PullRequests
	for _, v := range p.Nodes {
		pr := github.PullRequest{}
		pr.Id = fmt.Sprintf("%v#%v", v.Repository.NameWithOwner, v.Number)
		pr.Additions = v.Additions
		pr.Deletions = v.Deletions
		pr.ChangedFiles = v.ChangedFiles
		pr.Commits = v.Commits.TotalCount
		pr.Merged = v.Merged
		pr.Closed = v.Closed
		pr.Title = v.Title
		pr.Body = v.Body
		pr.OwnRepository = v.Repository.Owner.Login == options.Username
		prs = append(prs, pr)
	}

	if p.PageInfo.EndCursor == "" {
		return prs, nil
	}

	options.PullRequestCursor = p.PageInfo.EndCursor
	return c.listPullRequests(
		options,
		prs,
	)
}
