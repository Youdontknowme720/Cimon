package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type JobRunResponse struct {
	TotalCount int   `json:"total_count"`
	Jobs       []Job `json:"jobs"`
}

type Job struct {
	ID int `json:"id"`
	Status string `json:"status"`
	Conclusion string `json:"conclusion"`
	Name string `json:"name"`
	Steps []Step `json:"steps"`
}

type Step struct {
	Name string `json:"name"`
	Status string `json:"status"`
	Conclusion string `json:"conclusion"`
}

func (workflow Workflow) GetJobRuns(repo string, token string) ([]Job, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/actions/runs/%d/jobs", repo, workflow.ID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error during creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error during requesting: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Api failure: %s (Status: %d), CallURL: %s", string(body), resp.StatusCode, url)
	}

	var result JobRunResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("Error while Decoding: %w", err)
	}

	if result.TotalCount == 0 {
		fmt.Println("⚠️  Found no jobs for this specific workflow")
	}

	return result.Jobs, nil
}

func (job Job) GetSteps() ([]Step, error){
	var steps []Step
	for _, step := range job.Steps{
		if step.Name != "Set up job" && step.Name != "Complete job"{
			steps = append(steps, step)
		}
	}
	return steps, nil
}
