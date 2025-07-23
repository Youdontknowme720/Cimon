package github

import (
	"encoding/json"
	"fmt"
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

func GetJobRuns(repo string, workflowID int, token string) (JobRunResponse, error){
	url := fmt.Sprintf("https://api.github.com/repos/%s/actions/runs/%d/jobs", repo, workflowID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil{
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil{
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200{
		panic(err)
	}

	var result JobRunResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil{
		panic(err)
	}

	if result.TotalCount == 0 {
		fmt.Println("Found no Jobs")
	}
	return result, nil
}

func (job Job) DisplaySteps(){
	for _, step := range job.Steps{
		if step.Name != "Set up job" && step.Name != "Complete job"{
			fmt.Printf("StepName: %s, Conclusion: %s, Status: %s\n",step.Name, step.Conclusion, step.Status)
		}
	}
}
