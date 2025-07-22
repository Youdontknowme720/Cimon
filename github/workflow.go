package github

import (
    "encoding/json"
    "fmt"
    "net/http"
)
type WorkflowRunsResponse struct {
	TotalCount   int            `json:"total_count"`
	WorkflowRuns []Workflow `json:"workflow_runs"`
}
type Workflow struct {
    ID int `json:"id"`
	DisplayTitle string `json:"display_title"`
	Status string `json:"status"`
	Conclusion string `json:"conclusion"`
	HtmlUrl string `json:"html_url"`
}

func GetWorkflowStatus(repo string, limit int, token string) (WorkflowRunsResponse, error){
	url := fmt.Sprintf("https://api.github.com/repos/%s/actions/runs?per_page=%d", repo, limit)

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

	var result WorkflowRunsResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil{
		panic(err)
	}

	if result.TotalCount == 0 {
		fmt.Println("Found no workflows")
	}
	return result, nil
}