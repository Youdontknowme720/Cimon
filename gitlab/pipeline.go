// Package gitlab is used for api requests
package gitlab

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Pipeline struct {
	ID      int    `json:"id"`
	Status  string `json:"status"`
	Created string `json:"created_at"`
	WebURL  string `json:"web_url"`
}

const baseURL = "https://gitlab.com/api/v4"

type Pipelines []Pipeline

func GetAllPipelines(projectID, token string, perPage int) ([]Pipeline, error) {
	url := fmt.Sprintf("%s/projects/%s/pipelines?per_page=%d", baseURL, projectID, perPage)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("PRIVATE-TOKEN", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pipelines []Pipeline
	if err := json.NewDecoder(resp.Body).Decode(&pipelines); err != nil {
		return nil, err
	}

	return pipelines, nil
}
