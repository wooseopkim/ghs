package query

import (
	"fmt"

	"github.com/machinebox/graphql"
	"github.com/wooseopkim/ghs/internal/github"
)

type ListRepositoriesDetailedOptions struct {
	GitHubToken      string
	Username         string
	Id               string
	repositoryCursor string
	repositoryLimit  uint
	commitCursor     string
	commitLimit      uint
}

type detailedRepositoriesResponse struct {
	User struct {
		RepositoriesContributedTo detailedRepositories
		Repositories              detailedRepositories
	}
}

type detailedRepositories struct {
	PageInfo PageInfo
	Nodes    []struct {
		NameWithOwner    string
		DefaultBranchRef struct {
			Target struct {
				History struct {
					TotalCount uint
					PageInfo   PageInfo
					Nodes      []struct {
						Oid       string
						Additions uint
						Deletions uint
						Parents   struct {
							TotalCount uint
						}
					}
				}
			}
		}
	}
}

type detailedRepositoriesFragment = fragment[detailedRepositoriesResponse, detailedRepositories]

func (c *client) ListRepositoriesDetailed(
	options *ListRepositoriesDetailedOptions,
) ([]github.Commit, error) {
	return listRepositories(
		func(
			fragment detailedRepositoriesFragment,
			commits []github.Commit,
		) ([]github.Commit, error) {
			return c.listRepositoriesDetailed(
				options,
				fragment,
				commits,
				&Paginator[detailedRepositoriesResponse, ListRepositoriesDetailedOptions]{
					LastPage: func(res detailedRepositoriesResponse) bool {
						return res.User.Repositories.PageInfo.EndCursor == ""
					},
					NextPageOptions: func(
						res detailedRepositoriesResponse,
						options *ListRepositoriesDetailedOptions,
					) *ListRepositoriesDetailedOptions {
						newOptions := *options
						newOptions.repositoryCursor = res.User.Repositories.PageInfo.EndCursor
						return &newOptions
					},
				},
			)
		},
		detailedRepositoriesFragment{
			RepositoriesQuery: `repositories(
				affiliations: [OWNER, COLLABORATOR]
				first: $repositoryLimit
				after: $repositoryCursor
				isFork: false
			)`,
			RetrieveRepositories: func(res detailedRepositoriesResponse) detailedRepositories {
				return res.User.Repositories
			},
		},
		detailedRepositoriesFragment{
			RepositoriesQuery: `repositoriesContributedTo(
				first: $repositoryLimit
				after: $repositoryCursor
				contributionTypes: [COMMIT]
			)`,
			RetrieveRepositories: func(res detailedRepositoriesResponse) detailedRepositories {
				return res.User.RepositoriesContributedTo
			},
		},
	)
}

type Paginator[RES interface{}, OPTIONS interface{}] struct {
	LastPage        func(res RES) bool
	NextPageOptions func(res RES, options *OPTIONS) *OPTIONS
}

func (c *client) listRepositoriesDetailed(
	options *ListRepositoriesDetailedOptions,
	fragment detailedRepositoriesFragment,
	commits []github.Commit,
	paginator *Paginator[detailedRepositoriesResponse, ListRepositoriesDetailedOptions],
) ([]github.Commit, error) {
	req := graphql.NewRequest(`
        query ListCommits(
            $username: String!
            $id: ID
            $repositoryLimit: Int = 5
            $repositoryCursor: String
            $commitLimit: Int = 50
            $commitCursor: String
        ) {
            user(login: $username) {
			` + fragment.RepositoriesQuery + ` {
					...repoData
				}
			}
		}

		fragment repoData on RepositoryConnection {
			pageInfo {
				startCursor
				endCursor
			}
			nodes {
				nameWithOwner
				defaultBranchRef {
					target {
						...on Commit {
							history(
								author: { id: $id }
								first: $commitLimit
								after: $commitCursor
							) {
								totalCount
								pageInfo {
									startCursor
									endCursor
								}
								nodes {
									oid
									additions
									deletions
									parents(first: 1) {
										totalCount
									}
								}
							}
						}
					}
				}
			}
		}`)
	req.Header.Set("Authorization", "bearer "+options.GitHubToken)
	req.Var("username", options.Username)
	req.Var("id", options.Id)
	if options.repositoryCursor != "" {
		req.Var("repositoryCursor", options.repositoryCursor)
	}
	if options.repositoryLimit != 0 {
		req.Var("repositoryLimit", options.repositoryLimit)
	}
	if options.commitCursor != "" {
		req.Var("commitCursor", options.commitCursor)
	}
	if options.commitLimit != 0 {
		req.Var("commitLimit", options.commitLimit)
	}

	var res detailedRepositoriesResponse
	err := c.run(req, &res)
	if err != nil {
		return nil, err
	}

	r := fragment.RetrieveRepositories(res)
	for _, v := range r.Nodes {
		if options.repositoryLimit != 1 {
			moreCommitsOptions := *options
			moreCommitsOptions.repositoryLimit = 1
			moreCommitsOptions.commitLimit = 100
			moreCommitsPaginator := &Paginator[detailedRepositoriesResponse, ListRepositoriesDetailedOptions]{
				LastPage: func(res detailedRepositoriesResponse) bool {
					data := fragment.RetrieveRepositories(res)
					return data.Nodes[0].DefaultBranchRef.Target.History.PageInfo.EndCursor == ""
				},
				NextPageOptions: func(res detailedRepositoriesResponse, options *ListRepositoriesDetailedOptions) *ListRepositoriesDetailedOptions {
					data := fragment.RetrieveRepositories(res)
					newOptions := *options
					newOptions.commitCursor = data.Nodes[0].DefaultBranchRef.Target.History.PageInfo.EndCursor
					return &newOptions
				},
			}
			moreCommits, err := c.listRepositoriesDetailed(
				&moreCommitsOptions,
				fragment,
				[]github.Commit{},
				moreCommitsPaginator,
			)
			if err != nil {
				return nil, err
			}
			commits = append(commits, moreCommits...)
		}

		for _, vv := range v.DefaultBranchRef.Target.History.Nodes {
			if vv.Parents.TotalCount == 0 {
				continue
			}
			commit := github.Commit{}
			commit.Id = fmt.Sprintf("%v#%v", v.NameWithOwner, vv.Oid)
			commit.Additions = vv.Additions
			commit.Deletions = vv.Deletions
			commits = append(commits, commit)
		}
	}

	if err != nil {
		return nil, err
	}

	if paginator.LastPage(res) {
		return commits, nil
	}

	return c.listRepositoriesDetailed(
		paginator.NextPageOptions(res, options),
		fragment,
		commits,
		paginator,
	)
}
