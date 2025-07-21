package utils

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

type Pipelines []Pipeline

func GetPipelineStatus(projectID string, accessToken string) (Pipelines, error) {
	url := fmt.Sprintf("https://gitlab.com/api/v4/projects/%s/pipelines", projectID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("PRIVATE-TOKEN", accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pipelines Pipelines
	if err := json.NewDecoder(resp.Body).Decode(&pipelines); err != nil {
		return nil, err
	}

	return pipelines, nil

}
func (pipe Pipeline) IsFailed() bool{
	return pipe.Status == "failed" || pipe.Status == "canceled"
}

func (pipes Pipelines) HasFailedPipelines() (bool, Pipelines) {
	var failedPipeList Pipelines

	for _, pipe := range pipes {
		if pipe.IsFailed(){
			failedPipeList = append(failedPipeList, pipe)
		}
	}

	return len(failedPipeList) > 0, failedPipeList
}

func (pipes Pipelines)DisplayPipelines() {
	for _, pipe := range pipes {
		fmt.Printf("ID: %d | Status: %s | Created: %s | URL: %s\n",
			pipe.ID, pipe.Status, pipe.Created, pipe.WebURL)
	}
}
