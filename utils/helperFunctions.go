package utils

import "fmt"


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
func DisplayPipelines(pipelines []Pipeline) {
	for _, pipe := range pipelines{
		fmt.Printf("ID: %d | Status: %s | Created: %s | URL: %s\n",
			pipe.ID, pipe.Status, pipe.Created, pipe.WebURL)
	}
}

func DisplayJobs(jobs []Job) {
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

func HasFaildePipelines(pipeline []Pipeline) (bool, []Pipeline) {
	var hasFailed bool
	var failedPipeList []Pipeline

	for _, pipe := range pipeline{
		if pipe.Status == "failed" || pipe.Status == "canceled"{
			failedPipeList = append(failedPipeList, pipe)
			hasFailed = true
		}
	}

	return hasFailed, failedPipeList
}

