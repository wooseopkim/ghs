package ghs

import (
	"context"

	"github.com/wooseopkim/ghs/internal/analyzer"
	"github.com/wooseopkim/ghs/internal/query"
	"github.com/wooseopkim/ghs/internal/stats"
)

func GetStats(gitHubToken string) (stats.Records, error) {
	ctx := context.Background()
	client := query.NewClient(ctx, "https://api.github.com/graphql")

	user, err := client.GetUser(&query.GetUserOptions{
		GitHubToken: gitHubToken,
	})
	if err != nil {
		return stats.Records{}, err
	}

	repos, err := client.ListRepositoriesSummary(&query.ListRepositoriesSummaryOptions{
		Id:          user.Id,
		Username:    user.Username,
		GitHubToken: gitHubToken,
	})
	if err != nil {
		return stats.Records{}, err
	}

	commits, err := client.ListRepositoriesDetailed(&query.ListRepositoriesDetailedOptions{
		Id:          user.Id,
		Username:    user.Username,
		GitHubToken: gitHubToken,
	})
	if err != nil {
		return stats.Records{}, err
	}

	prs, err := client.ListPullRequests(&query.ListPullRequestsOptions{
		Username:    user.Username,
		GitHubToken: gitHubToken,
	})
	if err != nil {
		return stats.Records{}, err
	}

	prStats := analyzer.AnalyzePullRequests(prs)
	commitStats := analyzer.AnalyzeCommits(commits)
	repoStats := analyzer.AnalyzeRepositories(repos)

	return stats.Records{
		PullRequest: prStats,
		Commit:      commitStats,
		Repository:  repoStats,
	}, nil
}
