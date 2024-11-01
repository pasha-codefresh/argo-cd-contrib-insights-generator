package pkg

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"time"
)

// Data structure to match the JSON payload for the POST request
type QueryPayload struct {
	Queries []struct {
		RefID      string `json:"refId"`
		Datasource struct {
			UID  string `json:"uid"`
			Type string `json:"type"`
		} `json:"datasource"`
		RawSQL        string `json:"rawSql"`
		Format        string `json:"format"`
		DatasourceID  int    `json:"datasourceId"`
		IntervalMs    int    `json:"intervalMs"`
		MaxDataPoints int    `json:"maxDataPoints"`
	} `json:"queries"`
	Range struct {
		From string `json:"from"`
		To   string `json:"to"`
		Raw  struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"raw"`
	} `json:"range"`
	From string `json:"from"`
	To   string `json:"to"`
}

// Data structure to parse response JSON
type Response struct {
	Results map[string]struct {
		Frames []struct {
			Schema struct {
				Fields []struct {
					Name string `json:"name"`
				} `json:"fields"`
			} `json:"schema"`
			Data struct {
				Values [][]interface{} `json:"values"`
			} `json:"data"`
		} `json:"frames"`
	} `json:"results"`
}

// Contributor struct holds the username and total contributions
type Contributor struct {
	Username string
	Total    int
}

type Grafana struct {
}

func NewGrafana() *Grafana {
	return &Grafana{}
}

func (g *Grafana) TopArgoCDReviewers() ([]Contributor, error) {

	sql := "select\n  * \nfrom\n  suser_reviews\nwhere\n  $__timeFilter(time)\n  and period = 'd'\n  and series = 'rev_per_usrargoprojargocd'\norder by\n  time"

	from := strconv.FormatInt(time.Now().AddDate(0, 0, -7).UnixMilli(), 10)
	to := strconv.FormatInt(time.Now().UnixMilli(), 10)

	payloadJSON, err := json.Marshal(buildPayload(sql, "time_series", from, to))
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
	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	contributors := aggregateContributions(response)

	if len(contributors) < 10 {
		return contributors, nil
	}
	return contributors[0:10], nil
}

func (g *Grafana) TopArgoRolloutsReviewers() ([]Contributor, error) {

	sql := "select\n  * \nfrom\n  suser_reviews\nwhere\n  $__timeFilter(time)\n  and period = 'd'\n  and series = 'rev_per_usrargoprojargorollouts'\norder by\n  time"

	from := strconv.FormatInt(time.Now().AddDate(0, 0, -7).UnixMilli(), 10)
	to := strconv.FormatInt(time.Now().UnixMilli(), 10)

	payloadJSON, err := json.Marshal(buildPayload(sql, "time_series", from, to))
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
	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	contributors := aggregateContributions(response)

	return contributors, nil
}

func (g *Grafana) TopArgoRolloutsMergers() ([]Contributor, error) {

	sql := "select\n  row_number() over (order by value desc, name asc) as \"Rank\",\n  name,\n  value\nfrom\n  shpr_mergers\nwhere\n  series = 'hpr_mergersargoprojargorollouts'\n  and period = 'w'"

	from := strconv.FormatInt(time.Now().AddDate(0, 0, -7).UnixMilli(), 10)
	to := strconv.FormatInt(time.Now().UnixMilli(), 10)

	payloadJSON, err := json.Marshal(buildPayload(sql, "table", from, to))
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
	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	contributors := processContributors(response)

	return contributors, nil
}

func (g *Grafana) TopArgoCDMergers() ([]Contributor, error) {

	sql := "select\n  row_number() over (order by value desc, name asc) as \"Rank\",\n  name,\n  value\nfrom\n  shpr_mergers\nwhere\n  series = 'hpr_mergersargocd'\n  and period = 'w'"

	from := strconv.FormatInt(time.Now().AddDate(0, 0, -7).UnixMilli(), 10)
	to := strconv.FormatInt(time.Now().UnixMilli(), 10)

	payloadJSON, err := json.Marshal(buildPayload(sql, "table", from, to))
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
	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	contributors := processContributors(response)

	return contributors, nil
}

func buildPayload(sql string, format string, from string, to string) QueryPayload {
	// Payload for the POST request
	payload := QueryPayload{
		Queries: []struct {
			RefID      string `json:"refId"`
			Datasource struct {
				UID  string `json:"uid"`
				Type string `json:"type"`
			} `json:"datasource"`
			RawSQL        string `json:"rawSql"`
			Format        string `json:"format"`
			DatasourceID  int    `json:"datasourceId"`
			IntervalMs    int    `json:"intervalMs"`
			MaxDataPoints int    `json:"maxDataPoints"`
		}{
			{
				RefID: "A",
				Datasource: struct {
					UID  string `json:"uid"`
					Type string `json:"type"`
				}{
					UID:  "P172949F98CB31475",
					Type: "postgres",
				},
				RawSQL:        sql,
				Format:        format,
				DatasourceID:  1,
				IntervalMs:    3600000,
				MaxDataPoints: 1622,
			},
		},
		From: from,
		To:   to,
	}

	return payload
}

// Aggregates and sorts contributors by total contributions
func aggregateContributions(data Response) []Contributor {
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
	contributors := make([]Contributor, 0, len(contribMap))
	for username, total := range contribMap {
		if total > 0 {
			contributors = append(contributors, Contributor{Username: username, Total: int(total)})
		}
	}
	sort.Slice(contributors, func(i, j int) bool {
		return contributors[i].Total > contributors[j].Total
	})

	return contributors
}

func processContributors(data Response) []Contributor {
	var contributors []Contributor
	for _, result := range data.Results {
		for _, frame := range result.Frames {
			nameField := 1
			valueField := 2

			// Loop over values and capture each contributor's name and total contributions
			for i := 0; i < len(frame.Data.Values[nameField]); i++ {
				contributors = append(contributors, Contributor{
					Username: frame.Data.Values[nameField][i].(string),
					Total:    int(frame.Data.Values[valueField][i].(float64)),
				})
			}
		}
	}

	// Return the aggregated list of contributors
	return contributors
}
