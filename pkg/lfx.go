package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	lfxLeaderboardURL   = "https://insights.linuxfoundation.org/api/widget/contributors/contributor-leaderboard"
	lfxEmbedURL         = "https://insights.linuxfoundation.org/embed/project/argo"
	lfxProject          = "argo"
	lfxLeaderboardLimit = 10
	lfxWidget           = "contributors-leaderboard"

	activityReviewed = "pull_request-reviewed"
	activityMerged   = "pull_request-merged"
)

type LFX struct {
	client *http.Client
}

func NewLFX() *LFX {
	return &LFX{client: &http.Client{Timeout: 30 * time.Second}}
}

type lfxContributor struct {
	Name          string
	Contributions int
	Percentage    float64
}

type lfxProjectStats struct {
	Name   string
	People []lfxContributor
}

func lfxRepoSlug(repo string) string {
	return owner + "_" + repo
}

// lfxInsightsURL opens the leaderboard with the repo, date range, and activity already selected.
// The embed page is what applies metric; the main contributors page only keeps timeRange, start, end, and repos.
func lfxInsightsURL(repoSlugs []string, activityType, startDate, endDate string) string {
	query := url.Values{}
	query.Set("widget", lfxWidget)
	query.Set("timeRange", "custom")
	query.Set("start", startDate)
	query.Set("end", endDate)
	query.Set("timeRangeKey", "custom")
	query.Set("startDate", startDate)
	query.Set("endDate", endDate)
	query.Set("repos", strings.Join(repoSlugs, ","))
	query.Set("metric", "github:"+activityType)
	return lfxEmbedURL + "?" + query.Encode()
}

func (l *LFX) TopContributors(repo, activityType, startDate, endDate string) ([]lfxContributor, error) {
	query := url.Values{}
	query.Set("project", lfxProject)
	query.Set("platform", "github")
	query.Set("activityType", activityType)
	query.Set("repos", fmt.Sprintf("https://github.com/%s/%s", owner, repo))
	query.Set("startDate", startDate)
	query.Set("endDate", endDate)
	query.Set("offset", "0")
	query.Set("limit", fmt.Sprintf("%d", lfxLeaderboardLimit))
	query.Set("includeCollaborations", "false")

	req, err := http.NewRequest(http.MethodGet, lfxLeaderboardURL+"?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := l.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("lfx insights leaderboard: %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var payload struct {
		Data []struct {
			Name          string  `json:"name"`
			Contributions int     `json:"contributions"`
			Percentage    float64 `json:"percentage"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	contributors := make([]lfxContributor, 0, len(payload.Data))
	for _, row := range payload.Data {
		contributors = append(contributors, lfxContributor{
			Name:          row.Name,
			Contributions: row.Contributions,
			Percentage:    row.Percentage,
		})
	}
	return contributors, nil
}

func formatLFXGroup(title, suffix, sectionURL string, projects []lfxProjectStats) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "  - [%s](%s)%s\n", title, sectionURL, suffix)
	for _, project := range projects {
		fmt.Fprintf(&sb, "    - **%s:**\n", project.Name)
		for _, person := range project.People {
			fmt.Fprintf(&sb, "      - %s %d - %d%%\n", person.Name, person.Contributions, int(math.Round(person.Percentage)))
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}
