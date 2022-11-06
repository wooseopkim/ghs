package query

import (
	"context"

	"github.com/machinebox/graphql"
)

type client struct {
	client *graphql.Client
	ctx    context.Context
}

func NewClient(ctx context.Context, endpoint string) *client {
	return &client{
		client: graphql.NewClient(endpoint),
		ctx:    ctx,
	}
}

func (c *client) run(req *graphql.Request, res interface{}) error {
	return c.client.Run(c.ctx, req, &res)
}
