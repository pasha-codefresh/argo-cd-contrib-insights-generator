package pkg

import (
	"net/url"
	"testing"
)

func TestFormatLFXGroupRoundsPercentages(t *testing.T) {
	sectionURL := lfxInsightsURL([]string{"argoproj_argo-cd", "argoproj_argo-rollouts"}, activityReviewed, "2026-10-08", "2026-10-15")

	got := formatLFXGroup("Top Reviewers", ` (Choose "Reviewed a pull request" in the activity drop-down):`, sectionURL, []lfxProjectStats{
		{
			Name: "Argo CD",
			People: []lfxContributor{
				{Name: "Blake Pettersson", Contributions: 36, Percentage: 22.22},
				{Name: "Michael Crenshaw", Contributions: 30, Percentage: 18.52},
			},
		},
		{
			Name: "Argo Rollouts",
			People: []lfxContributor{
				{Name: "Kostis Kapelonis", Contributions: 12, Percentage: 92.31},
				{Name: "Zach Aller", Contributions: 1, Percentage: 7.69},
			},
		},
	})

	want := "" +
		"  - [Top Reviewers](" + sectionURL + ") (Choose \"Reviewed a pull request\" in the activity drop-down):\n" +
		"    - **Argo CD:**\n" +
		"      - Blake Pettersson 36 - 22%\n" +
		"      - Michael Crenshaw 30 - 19%\n" +
		"    - **Argo Rollouts:**\n" +
		"      - Kostis Kapelonis 12 - 92%\n" +
		"      - Zach Aller 1 - 8%"

	if got != want {
		t.Fatalf("format mismatch\n got:\n%s\nwant:\n%s", got, want)
	}
}

func TestLFXInsightsURLAppliesFilters(t *testing.T) {
	raw := lfxInsightsURL([]string{"argoproj_argo-cd"}, activityMerged, "2026-10-08", "2026-10-15")
	parsed, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Scheme != "https" || parsed.Host != "insights.linuxfoundation.org" || parsed.Path != "/embed/project/argo" {
		t.Fatalf("unexpected url %s", raw)
	}
	query := parsed.Query()
	checks := map[string]string{
		"widget":       "contributors-leaderboard",
		"timeRange":    "custom",
		"start":        "2026-10-08",
		"end":          "2026-10-15",
		"timeRangeKey": "custom",
		"startDate":    "2026-10-08",
		"endDate":      "2026-10-15",
		"repos":        "argoproj_argo-cd",
		"metric":       "github:pull_request-merged",
	}
	for key, want := range checks {
		if got := query.Get(key); got != want {
			t.Errorf("%s = %q, want %q", key, got, want)
		}
	}
}
