package pkg

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/go-github/v50/github"
)

const (
	owner            = "argoproj"
	ArgoCDRepo       = "argo-cd"
	ArgoRolloutsRepo = "argo-rollouts"
)

type Github struct {
	repo string
}

func NewGithub(repo string) *Github {
	return &Github{repo: repo}
}

func (g *Github) SearchLink(query string) string {
	return fmt.Sprintf("https://github.com/%s/%s/issues?q=%s", owner, g.repo, url.QueryEscape(query))
}

func (g *Github) GetCreatedAndClosedIssues(startDate, endDate string) (int, int, error) {
	client := github.NewClient(nil) // No auth client for public repos
	issuesCreatedQuery := fmt.Sprintf("repo:%s/%s is:issue is:open created:%s..%s", owner, g.repo, startDate, endDate)
	issuesCreated, err := executeSearchQuery(context.Background(), client, issuesCreatedQuery)
	if err != nil {
		return 0, 0, err
	}
	issuesClosedQuery := fmt.Sprintf("repo:%s/%s is:issue is:closed created:%s..%s", owner, g.repo, startDate, endDate)
	issuesClosed, err := executeSearchQuery(context.Background(), client, issuesClosedQuery)
	if err != nil {
		return 0, 0, err
	}
	return issuesCreated, issuesClosed, nil
}

func (g *Github) GetCreatedAndClosedPRs(startDate, endDate string) (int, int, error) {
	client := github.NewClient(nil) // No auth client for public repos
	issuesCreatedQuery := fmt.Sprintf("repo:%s/%s is:pr is:open created:%s..%s", owner, g.repo, startDate, endDate)
	issuesCreated, err := executeSearchQuery(context.Background(), client, issuesCreatedQuery)
	if err != nil {
		return 0, 0, err
	}
	issuesClosedQuery := fmt.Sprintf("repo:%s/%s is:pr is:closed created:%s..%s", owner, g.repo, startDate, endDate)
	issuesClosed, err := executeSearchQuery(context.Background(), client, issuesClosedQuery)
	if err != nil {
		return 0, 0, err
	}
	return issuesCreated, issuesClosed, nil
}

func (g *Github) GetStaleIssues(startDate string) (int, error) {
	client := github.NewClient(nil) // No auth client for public repos
	staleIssuesQuery := fmt.Sprintf("repo:%s/%s is:issue is:open updated:<%s", owner, g.repo, startDate)
	staleIssues, err := executeSearchQuery(context.Background(), client, staleIssuesQuery)
	if err != nil {
		return 0, err
	}
	return staleIssues, nil
}

func executeSearchQuery(ctx context.Context, client *github.Client, query string) (int, error) {
	options := &github.SearchOptions{ListOptions: github.ListOptions{}}
	result, _, err := client.Search.Issues(ctx, query, options)
	if err != nil {
		return 0, err
	}
	return result.GetTotal(), nil
}
