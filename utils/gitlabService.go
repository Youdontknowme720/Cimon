package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Pipeline struct {
	ID      int    `json:"id"`
	Status  string `json:"status"`
	Created string `json:"created_at"`
	WebURL  string `json:"web_url"`
}

type Job struct{
	ID int `json:"id"`
	Name string `json:"name"`
	Stage string `json:"stage"`
	Status string `json:"status"`
	Duration float64 `json:"duration"`
	WebURL string `json:"web_url"`
}

func GetPiplineStatus(projectID string, accessToken string) ([]Pipeline, error) {
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
func GetJobDetails(projectID string, pipelineID int, accessToken string) ([]Job, error) {
	url := fmt.Sprintf("https://gitlab.com/api/v4/projects/%s/pipelines/%d/jobs", projectID, pipelineID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Erstellen der Anfrage: %w", err)
	}

	req.Header.Add("PRIVATE-TOKEN", accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Ausführen der Anfrage: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitLab API Fehler (Status: %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Lesen der Response: %w", err)
	}

	var jobs []Job
	if err := json.Unmarshal(body, &jobs); err != nil {
		return nil, fmt.Errorf("fehler beim Parsen der JSON-Response: %w", err)
	}

	return jobs, nil
}

func GetJobsLog(projectID string, jobID int, accessToken string) (string, error) {
	url := fmt.Sprintf("https://gitlab.com/api/v4/projects/%s/jobs/%d/trace", projectID, jobID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("PRIVATE-TOKEN", accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Für Logs ist die Response ein plain text, nicht JSON
	buf := new(strings.Builder)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
