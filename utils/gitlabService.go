package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Pipeline struct{
	ID int `json:"id"`
	Status string `json:"status"`
	Created string `json:"created_at"`
	WebURL string `json:"web_url"`
}

func GetPiplineStatus(projectID string, accessToken string) ([]Pipeline, error){
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

	var pipelines []Pipeline
	if err := json.NewDecoder(resp.Body).Decode(&pipelines); err != nil {
		return nil, err
	}

	return pipelines, nil

}

func DisplayPipelines(pipelines []Pipeline) {
	for _, pipe := range pipelines{
		fmt.Printf("ID: %d | Status: %s | Created: %s | URL: %s\n",
			pipe.ID, pipe.Status, pipe.Created, pipe.WebURL)
	}
}
