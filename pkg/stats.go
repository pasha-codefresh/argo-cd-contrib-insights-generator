package pkg

import (
	"fmt"
	"time"

	"github.com/pasha-codefresh/argo-cd-contrib-insights-generator/pkg/util"
)

type StatsGenerator interface {
	Generate() (string, error)
}

type createdIssuesStatsGenerator struct {
	github *Github
}

func NewCreatedIssuesStatsGenerator(repo string) StatsGenerator {
	return &createdIssuesStatsGenerator{
		github: NewGithub(repo),
	}
}

func (c *createdIssuesStatsGenerator) Generate() (string, error) {
	startDate, endDate := util.GetRangeForLastWeek()

	issuesCreated, issuesClosed, err := c.github.GetCreatedAndClosedIssues(startDate, endDate)
	if err != nil {
		return "", err
	}
	openLink := c.github.SearchLink(fmt.Sprintf("is:issue is:open created:%s..%s", startDate, endDate))
	closedLink := c.github.SearchLink(fmt.Sprintf("is:issue is:closed created:%s..%s", startDate, endDate))
	return fmt.Sprintf("**Created Issues:** [%d](%s) open / [%d](%s) closed", issuesCreated, openLink, issuesClosed, closedLink), nil
}

type createdPRsStatsGenerator struct {
	github *Github
}

func NewCreatedPRsStatsGenerator(repo string) StatsGenerator {
	return &createdPRsStatsGenerator{
		github: NewGithub(repo),
	}
}

func (c *createdPRsStatsGenerator) Generate() (string, error) {
	startDate, endDate := util.GetRangeForLastWeek()

	prsCreated, prsClosed, err := c.github.GetCreatedAndClosedPRs(startDate, endDate)
	if err != nil {
		return "", err
	}
	openLink := c.github.SearchLink(fmt.Sprintf("is:pr is:open created:%s..%s", startDate, endDate))
	closedLink := c.github.SearchLink(fmt.Sprintf("is:pr is:closed created:%s..%s", startDate, endDate))
	return fmt.Sprintf("**Created PRs:** [%d](%s) open / [%d](%s) closed", prsCreated, openLink, prsClosed, closedLink), nil
}

type staleIssuesStatsGenerator struct {
	github *Github
}

func NewStaleIssuesStatsGenerator(repo string) StatsGenerator {
	return &staleIssuesStatsGenerator{
		github: NewGithub(repo),
	}
}

func (c *staleIssuesStatsGenerator) Generate() (string, error) {
	startDate := time.Now().AddDate(-1, 0, 0).Format("2006-01-02")
	staleIssues, err := c.github.GetStaleIssues(startDate)
	if err != nil {
		return "", err
	}
	link := c.github.SearchLink(fmt.Sprintf("is:issue is:open updated:<%s", startDate))
	return fmt.Sprintf("**Stale Issues:** [%d](%s)", staleIssues, link), nil
}

type topReviewersStatsGenerator struct {
	lfx *LFX
}

func NewTopReviewersStatsGenerator() StatsGenerator {
	return &topReviewersStatsGenerator{
		lfx: NewLFX(),
	}
}

func (c *topReviewersStatsGenerator) Generate() (string, error) {
	return generateLFXGroup(c.lfx, "Top Reviewers", ` (Choose "Reviewed a pull request" in the activity drop-down):`, activityReviewed)
}

type topMergersStatsGenerator struct {
	lfx *LFX
}

func NewTopMergersStatsGenerator() StatsGenerator {
	return &topMergersStatsGenerator{
		lfx: NewLFX(),
	}
}

func (c *topMergersStatsGenerator) Generate() (string, error) {
	return generateLFXGroup(c.lfx, "Top Mergers", `: (Choose "Merged a pull request" in the activity drop-down):`, activityMerged)
}

func generateLFXGroup(lfx *LFX, title, suffix, activityType string) (string, error) {
	startDate, endDate := util.GetRangeForLastWeek()

	argoCD, err := lfx.TopContributors(ArgoCDRepo, activityType, startDate, endDate)
	if err != nil {
		return "", err
	}
	argoRollouts, err := lfx.TopContributors(ArgoRolloutsRepo, activityType, startDate, endDate)
	if err != nil {
		return "", err
	}

	return formatLFXGroup(title, suffix, lfxInsightsURL([]string{lfxRepoSlug(ArgoCDRepo), lfxRepoSlug(ArgoRolloutsRepo)}, activityType, startDate, endDate), []lfxProjectStats{
		{Name: "Argo CD", People: argoCD},
		{Name: "Argo Rollouts", People: argoRollouts},
	}), nil
}
