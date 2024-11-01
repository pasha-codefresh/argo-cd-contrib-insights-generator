package pkg

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/pasha-codefresh/argo-cd-contrib-insights-generator/pkg/types"
	"github.com/pasha-codefresh/argo-cd-contrib-insights-generator/pkg/util"
)

type Grafana struct {
	queryPayloadFactory              *types.QueryPayloadFactory
	timeSeriesContributorsAggregator ContributorsAggregator
	tableContributorsAggregator      ContributorsAggregator
}

func NewGrafana() *Grafana {
	return &Grafana{
		queryPayloadFactory:              types.NewQueryPayloadFactory(),
		timeSeriesContributorsAggregator: NewTimeSeriesContributorsAggregator(),
		tableContributorsAggregator:      NewTableContributorsAggregator(),
	}
}

func (g *Grafana) TopArgoCDReviewers() ([]types.Contributor, error) {

	sql := "select\n  * \nfrom\n  suser_reviews\nwhere\n  $__timeFilter(time)\n  and period = 'd'\n  and series = 'rev_per_usrargoprojargocd'\norder by\n  time"

	from, to := util.GetRangeForLastWeekAsMilli()

	payloadJSON, err := json.Marshal(g.queryPayloadFactory.Create(sql, "time_series", from, to))
	if err != nil {
		return nil, err
	}

	// Make the POST request
	resp, err := http.Post("https://argo.devstats.cncf.io/api/ds/query", "application/json", bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response
	var response types.Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	contributors := g.timeSeriesContributorsAggregator.aggregate(response)

	if len(contributors) < 10 {
		return contributors, nil
	}
	return contributors[0:10], nil
}

func (g *Grafana) TopArgoRolloutsReviewers() ([]types.Contributor, error) {

	sql := "select\n  * \nfrom\n  suser_reviews\nwhere\n  $__timeFilter(time)\n  and period = 'd'\n  and series = 'rev_per_usrargoprojargorollouts'\norder by\n  time"

	from, to := util.GetRangeForLastWeekAsMilli()

	payloadJSON, err := json.Marshal(g.queryPayloadFactory.Create(sql, "time_series", from, to))
	if err != nil {
		return nil, err
	}

	// Make the POST request
	resp, err := http.Post("https://argo.devstats.cncf.io/api/ds/query", "application/json", bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response
	var response types.Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	contributors := g.timeSeriesContributorsAggregator.aggregate(response)

	return contributors, nil
}

func (g *Grafana) TopArgoRolloutsMergers() ([]types.Contributor, error) {

	sql := "select\n  row_number() over (order by value desc, name asc) as \"Rank\",\n  name,\n  value\nfrom\n  shpr_mergers\nwhere\n  series = 'hpr_mergersargoprojargorollouts'\n  and period = 'w'"

	from, to := util.GetRangeForLastWeekAsMilli()

	payloadJSON, err := json.Marshal(g.queryPayloadFactory.Create(sql, "table", from, to))
	if err != nil {
		return nil, err
	}

	// Make the POST request
	resp, err := http.Post("https://argo.devstats.cncf.io/api/ds/query", "application/json", bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response
	var response types.Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	contributors := g.tableContributorsAggregator.aggregate(response)

	return contributors, nil
}

func (g *Grafana) TopArgoCDMergers() ([]types.Contributor, error) {

	sql := "select\n  row_number() over (order by value desc, name asc) as \"Rank\",\n  name,\n  value\nfrom\n  shpr_mergers\nwhere\n  series = 'hpr_mergersargocd'\n  and period = 'w'"

	from := strconv.FormatInt(time.Now().AddDate(0, 0, -7).UnixMilli(), 10)
	to := strconv.FormatInt(time.Now().UnixMilli(), 10)

	payloadJSON, err := json.Marshal(g.queryPayloadFactory.Create(sql, "table", from, to))
	if err != nil {
		return nil, err
	}

	// Make the POST request
	resp, err := http.Post("https://argo.devstats.cncf.io/api/ds/query", "application/json", bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response
	var response types.Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	contributors := g.tableContributorsAggregator.aggregate(response)

	return contributors, nil
}

type ContributorsAggregator interface {
	aggregate(data types.Response) []types.Contributor
}

type timeSeriesContributorsAggregator struct {
}

func NewTimeSeriesContributorsAggregator() ContributorsAggregator {
	return &timeSeriesContributorsAggregator{}
}

// Aggregates and sorts contributors by total contributions
func (aggregator *timeSeriesContributorsAggregator) aggregate(data types.Response) []types.Contributor {
	contribMap := make(map[string]float64)
	for _, result := range data.Results {
		for _, frame := range result.Frames {
			for i, field := range frame.Schema.Fields {
				if i == 0 { // Skip "Time" field
					continue
				}
				username := field.Name
				for _, contrib := range frame.Data.Values[i] {
					total := contrib.(float64)
					if contrib != nil && total > 0 {
						contribMap[username] += total
					}
				}
			}
		}
	}

	// Convert map to sorted slice
	contributors := make([]types.Contributor, 0, len(contribMap))
	for username, total := range contribMap {
		if total > 0 {
			contributors = append(contributors, types.Contributor{Username: username, Total: int(total)})
		}
	}
	sort.Slice(contributors, func(i, j int) bool {
		return contributors[i].Total > contributors[j].Total
	})

	return contributors
}

type tableContributorsAggregator struct {
}

func NewTableContributorsAggregator() ContributorsAggregator {
	return &tableContributorsAggregator{}
}

func (aggregator *tableContributorsAggregator) aggregate(data types.Response) []types.Contributor {
	var contributors []types.Contributor
	for _, result := range data.Results {
		for _, frame := range result.Frames {
			nameField := 1
			valueField := 2

			// Loop over values and capture each contributor's name and total contributions
			for i := 0; i < len(frame.Data.Values[nameField]); i++ {
				contributors = append(contributors, types.Contributor{
					Username: frame.Data.Values[nameField][i].(string),
					Total:    int(frame.Data.Values[valueField][i].(float64)),
				})
			}
		}
	}

	// Return the aggregated list of contributors
	return contributors
}
