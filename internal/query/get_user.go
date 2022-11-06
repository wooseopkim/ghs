package query

import (
	"github.com/machinebox/graphql"
	"github.com/wooseopkim/ghs/internal/github"
)

type GetUserOptions struct {
	GitHubToken string
}

func (c *client) GetUser(
	options *GetUserOptions,
) (github.User, error) {
	req := graphql.NewRequest(`
        query GetUser {
            viewer {
                id
                login
            }
        }
    `)
	req.Header.Set("Authorization", "bearer "+options.GitHubToken)

	var res struct {
		Viewer struct {
			Id    string
			Login string
		}
	}
	err := c.run(req, &res)
	if err != nil {
		return github.User{}, err
	}

	return github.User{
		Id:       res.Viewer.Id,
		Username: res.Viewer.Login,
	}, nil
}
