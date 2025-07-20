package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func statusEmoji(status string) string {
	switch status {
	case "success":
		return "✅"
		case "failed":
		return "❌"
		case "running":
		return "🏃"
		case "pending":
		return "⏳"
		case "canceled":
		return "🚫"
		case "manual":
		return "✋"
		case "skipped":
		return "⤵️"
		default:
		return "❔"
	}
}


type Job struct{
	ID int `json:"id"`
	Name string `json:"name"`
	Stage string `json:"stage"`
	Status string `json:"status"`
	Duration float64 `json:"duration"`
	WebURL string `json:"web_url"`
}

type Jobs []Job

func GetJobDetails(projectID string, pipelineID int, accessToken string) (Jobs, error) {
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

func (job Job) GetJobsLog(projectID string, accessToken string) (string, error) {
	url := fmt.Sprintf("https://gitlab.com/api/v4/projects/%s/jobs/%d/trace", projectID, job.ID)
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

func (jobs Jobs) DisplayJobs() {
	for _, job := range jobs {
		emoji := statusEmoji(job.Status)
		fmt.Println("─────────────────────────────────────────────")
		fmt.Printf("%s  Job-ID:    %d\n", emoji, job.ID)
		fmt.Printf("👀  Name:      %s\n", job.Name)
		fmt.Printf("📦  Stage:     %s\n", job.Stage)
		fmt.Printf("📊  Status:    %s\n", job.Status)
		fmt.Printf("⏱️  Dauer:     %.2fs\n", job.Duration)
		fmt.Printf("🔗  URL:       %s\n", job.WebURL)
	}
	fmt.Println("─────────────────────────────────────────────")
}
