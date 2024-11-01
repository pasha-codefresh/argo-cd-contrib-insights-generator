package types

import (
	"fmt"
	"strings"
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

type QueryPayloadFactory struct {
}

func NewQueryPayloadFactory() *QueryPayloadFactory {
	return &QueryPayloadFactory{}
}

func (factory *QueryPayloadFactory) Create(sql string, format string, from string, to string) QueryPayload {
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

func ContributorsToString(contributors []Contributor) string {
	var sb strings.Builder
	for i, contributor := range contributors {
		if i == len(contributors)-1 {
			sb.WriteString(fmt.Sprintf("%s (%d)", contributor.Username, contributor.Total))
		} else {
			sb.WriteString(fmt.Sprintf("%s (%d), ", contributor.Username, contributor.Total))
		}
	}
	return sb.String()
}
