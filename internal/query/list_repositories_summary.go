package query

import (
	"time"

	"github.com/machinebox/graphql"
	"github.com/wooseopkim/ghs/internal/github"
)

type ListRepositoriesSummaryOptions struct {
	GitHubToken      string
	Id               string
	Username         string
	repositoryCursor string
}

type summaryRepositoriesReponse struct {
	User struct {
		RepositoriesContributedTo summaryRepositories
		Repositories              summaryRepositories
	}
}

type summaryRepositories struct {
	TotalCount uint
	PageInfo   PageInfo
	Nodes      []struct {
		Name  string
		Owner struct {
			Login string
		}
		NameWithOwner string
		Languages     struct {
			Edges []struct {
				Size uint
			}
			Nodes []struct {
				Name string
			}
		}
		Stargazers struct {
			TotalCount uint
		}
		IsPrivate bool
		CreatedAt time.Time
		UpdatedAt time.Time
	}
}

type summaryRepositoriesFragment = fragment[summaryRepositoriesReponse, summaryRepositories]

func (c *client) ListRepositoriesSummary(
	options *ListRepositoriesSummaryOptions,
) ([]github.Repository, error) {
	return listRepositories(
		func(
			fragment summaryRepositoriesFragment,
			data []github.Repository,
		) ([]github.Repository, error) {
			return c.listRepositoriesSummary(
				options,
				fragment,
				data,
			)
		},
		summaryRepositoriesFragment{
			RepositoriesQuery: `repositories(
				affiliations: [OWNER, COLLABORATOR]
				first: $repositoryLimit
				after: $repositoryCursor
				isFork: false
			)`,
			RetrieveRepositories: func(res summaryRepositoriesReponse) summaryRepositories {
				return res.User.Repositories
			},
		},
		summaryRepositoriesFragment{
			RepositoriesQuery: `repositoriesContributedTo(
				first: $repositoryLimit
				after: $repositoryCursor
				contributionTypes: [COMMIT]
			)`,
			RetrieveRepositories: func(res summaryRepositoriesReponse) summaryRepositories {
				return res.User.RepositoriesContributedTo
			},
		},
	)
}

func (c *client) listRepositoriesSummary(
	options *ListRepositoriesSummaryOptions,
	fragment summaryRepositoriesFragment,
	repos []github.Repository,
) ([]github.Repository, error) {
	req := graphql.NewRequest(`
        query ListRepositories(
            $username: String!
            $repositoryLimit: Int = 25
            $languageLimit: Int = 10
            $repositoryCursor: String
        ) {
            user(login: $username) {
				` + fragment.RepositoriesQuery + ` {
					totalCount
					pageInfo {
						endCursor
						startCursor
					}
					nodes {
						owner {
							login
						}
						nameWithOwner
						languages(first: $languageLimit, orderBy: { field: SIZE, direction: DESC }) {
							edges {
								size
							}
							nodes {
								name
							}
						}
						stargazers {
							totalCount
						}
						isPrivate
						createdAt
						updatedAt
					}
				}
			}
		}`)
	req.Header.Set("Authorization", "bearer "+options.GitHubToken)
	req.Var("id", options.Id)
	req.Var("username", options.Username)
	if options.repositoryCursor != "" {
		req.Var("repositoryCursor", options.repositoryCursor)
	}

	var res summaryRepositoriesReponse
	err := c.run(req, &res)
	if err != nil {
		return nil, err
	}

	r := fragment.RetrieveRepositories(res)
	for _, v := range r.Nodes {
		repo := github.Repository{}
		repo.Name = string(v.Name)
		repo.Owner = string(v.Owner.Login)
		repo.Id = string(v.NameWithOwner)
		for i := range v.Languages.Edges {
			repo.Languages = append(repo.Languages, github.Language{
				Name: string(v.Languages.Nodes[i].Name),
				Size: uint(v.Languages.Edges[i].Size),
			})
		}
		repo.Stars = uint(v.Stargazers.TotalCount)
		repo.Private = bool(v.IsPrivate)
		repo.CreatedAt = v.CreatedAt
		repo.UpdatedAt = v.UpdatedAt
		repos = append(repos, repo)
	}

	if r.PageInfo.EndCursor == "" {
		return repos, nil
	}

	options.repositoryCursor = r.PageInfo.EndCursor
	return c.listRepositoriesSummary(
		options,
		fragment,
		repos,
	)
}
