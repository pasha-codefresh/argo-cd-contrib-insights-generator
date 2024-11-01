package pkg

import (
	"fmt"
	"strings"
	"time"

	"github.com/pasha-codefresh/argo-cd-contrib-insights-generator/pkg/types"
	"github.com/pasha-codefresh/argo-cd-contrib-insights-generator/pkg/util"
)

type StatsGenerator interface {
	Generate() (string, string, error)
}

type createdIssuesStatsGenerator struct {
	github *Github
}

func NewCreatedIssuesStatsGenerator() StatsGenerator {
	return &createdIssuesStatsGenerator{
		github: NewGithub(),
	}
}

func (c *createdIssuesStatsGenerator) Generate() (string, string, error) {
	startDate, endDate := util.GetRangeForLastWeek()

	issuesCreated, issuesClosed, err := c.github.GetCreatedAndClosedIssues(startDate, endDate)
	if err != nil {
		return "", "", err
	}
	link := fmt.Sprintf("https://github.com/argoproj/argo-cd/issues?q=is%%3Aissue+is%%3Aopen+created%%%s..%s+", startDate, endDate)
	return fmt.Sprintf("Created Issues: %d open / %d closed", issuesCreated, issuesClosed), link, nil
}

type createdPRsStatsGenerator struct {
	github *Github
}

func NewCreatedPRsStatsGenerator() StatsGenerator {
	return &createdPRsStatsGenerator{
		github: NewGithub(),
	}
}

func (c *createdPRsStatsGenerator) Generate() (string, string, error) {
	startDate, endDate := util.GetRangeForLastWeek()

	prsCreated, prsClosed, err := c.github.GetCreatedAndClosedPRs()
	if err != nil {
		return "", "", err
	}
	link := fmt.Sprintf("https://github.com/argoproj/argo-cd/issues?q=is%%3Aissue+is%%3Aopen+created%%%s..%s+", startDate, endDate)
	return fmt.Sprintf("Created PRs: %d open / %d closed", prsCreated, prsClosed), link, nil
}

type staleIssuesStatsGenerator struct {
	github *Github
}

func NewStaleIssuesStatsGenerator() StatsGenerator {
	return &staleIssuesStatsGenerator{
		github: NewGithub(),
	}
}

func (c *staleIssuesStatsGenerator) Generate() (string, string, error) {
	startDate := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
	staleIssues, err := c.github.GetStaleIssues(startDate)
	if err != nil {
		return "", "", err
	}
	link := fmt.Sprintf("https://github.com/argoproj/argo-cd/issues?q=is%%3Aissue+is%%3Aopen+updated%%3A%%3C%s+", startDate)
	return fmt.Sprintf("Stale Issues: %d", staleIssues), link, nil
}

type topReviewersStatsGenerator struct {
	grafana *Grafana
}

func NewTopReviewersStatsGenerator() StatsGenerator {
	return &topReviewersStatsGenerator{
		grafana: NewGrafana(),
	}
}

func (c *topReviewersStatsGenerator) Generate() (string, string, error) {
	argocdreviewers, err := c.grafana.TopArgoCDReviewers()
	if err != nil {
		return "", "", err
	}

	argorolloutsreviewers, err := c.grafana.TopArgoRolloutsReviewers()
	if err != nil {
		return "", "", err
	}

	var sb strings.Builder
	sb.WriteString("Argo CD: ")
	sb.WriteString(types.ContributorsToString(argocdreviewers))

	sb.WriteString("\nArgo Rollouts: ")
	sb.WriteString(types.ContributorsToString(argorolloutsreviewers))

	link := "https://argo.devstats.cncf.io/d/29/pr-reviews-by-contributor?orgId=1&from=now-7d&to=now&var-period=d&var-repo_name=argoproj%2Fargo-cd"

	return sb.String(), link, nil
}

type topMergersStatsGenerator struct {
	grafana *Grafana
}

func NewTopMergersStatsGenerator() StatsGenerator {
	return &topMergersStatsGenerator{
		grafana: NewGrafana(),
	}
}

func (c *topMergersStatsGenerator) Generate() (string, string, error) {
	argocdmergers, err := c.grafana.TopArgoCDMergers()
	if err != nil {
		return "", "", err
	}
	argorolloutsmergers, err := c.grafana.TopArgoRolloutsMergers()
	if err != nil {
		return "", "", err
	}

	var sb strings.Builder
	sb.WriteString("Argo CD: ")
	sb.WriteString(types.ContributorsToString(argocdmergers))

	sb.WriteString("\nArgo Rollouts: ")
	sb.WriteString(types.ContributorsToString(argorolloutsmergers))

	return sb.String(), "https://argo.devstats.cncf.io/d/75/prs-mergers-table?orgId=1&var-period_name=Last%20week&var-repogroup_name=All", nil
}
