package pkg

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v50/github"
	"github.com/pasha-codefresh/argo-cd-contrib-insights-generator/pkg/util"
)

const DefaultRepo = "argoproj/argo-cd"

type Github struct {
	repo string
}

// NewGithub returns a client for the given repository in "owner/name" format.
func NewGithub(repo string) *Github {
	return &Github{repo: repo}
}

// ValidateRepo checks that repo is in "owner/name" format.
func ValidateRepo(repo string) error {
	owner, name, ok := strings.Cut(repo, "/")
	if !ok || owner == "" || name == "" || strings.Contains(name, "/") {
		return fmt.Errorf("invalid repo %q: expected format owner/name", repo)
	}
	return nil
}

func (g *Github) Repo() string {
	return g.repo
}

func (g *Github) GetCreatedAndClosedIssues(startDate, endDate string) (int, int, error) {
	client := github.NewClient(nil) // No auth client for public repos
	issuesCreatedQuery := fmt.Sprintf("repo:%s is:issue is:open created:%s..%s", g.repo, startDate, endDate)
	issuesCreated, err := executeSearchQuery(context.Background(), client, issuesCreatedQuery)
	if err != nil {
		return 0, 0, err
	}
	issuesClosedQuery := fmt.Sprintf("repo:%s is:issue is:closed created:%s..%s", g.repo, startDate, endDate)
	issuesClosed, err := executeSearchQuery(context.Background(), client, issuesClosedQuery)
	if err != nil {
		return 0, 0, err
	}
	return issuesCreated, issuesClosed, nil
}

func (g *Github) GetCreatedAndClosedPRs() (int, int, error) {
	client := github.NewClient(nil) // No auth client for public repos
	startDate, endDate := util.GetRangeForLastWeek()
	issuesCreatedQuery := fmt.Sprintf("repo:%s is:pr is:open created:%s..%s", g.repo, startDate, endDate)
	issuesCreated, err := executeSearchQuery(context.Background(), client, issuesCreatedQuery)
	if err != nil {
		return 0, 0, err
	}
	issuesClosedQuery := fmt.Sprintf("repo:%s is:pr is:closed created:%s..%s", g.repo, startDate, endDate)
	issuesClosed, err := executeSearchQuery(context.Background(), client, issuesClosedQuery)
	if err != nil {
		return 0, 0, err
	}
	return issuesCreated, issuesClosed, nil
}

func (g *Github) GetStaleIssues(startDate string) (int, error) {
	client := github.NewClient(nil) // No auth client for public repos
	staleIssuesQuery := fmt.Sprintf("repo:%s is:issue is:open updated:<%s", g.repo, startDate)
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
