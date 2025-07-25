package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/Youdontknowme720/Cimon/gitlab"
)

var projectID string
var cnt int

var pipelineCmd = &cobra.Command{
	Use:   "pl",
	Short: "Zeigt Pipeline-Status von GitLab",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := cmd.Flags().GetString("tokenGitlab")
		if err != nil {
			return err
		}

		if token == "" {
			config, err := utils.LoadConfig()
			if err != nil {
				return fmt.Errorf("Couldn't load config %w: ", err)
			}
			token = config.TokenGitlab
		}

		if token == "" {
			return fmt.Errorf("TokenGitlab is empty. Please set it using auth command or as a flag.")
		}

		var limit int
		pipelines, err := utils.GetPipelineStatus(projectID, token)
		if err != nil {
			return fmt.Errorf("Error while fetching pipelines: %w", err)
		}
		countSet := cmd.Flags().Changed("counter")
		if countSet {
			limit = cnt
			if cnt > len(pipelines) {
				limit = len(pipelines)
			}
			pipelines[:limit].DisplayPipelines()

		} else{
			pipelines.DisplayPipelines()
		}

		hasFailed, failedPipes := pipelines[:limit].HasFailedPipelines()

		if hasFailed{
			for _, pipe := range failedPipes{
				failedJobs, err := utils.GetJobDetails(projectID, pipe.ID, token)
				if err != nil{
					fmt.Printf("❌ FEHLER beim Abrufen der Jobs: %v\n", err)
				}
				failedJobs.DisplayJobs()
				logs, _ := failedJobs[0].GetJobsLog(projectID, token)
				fmt.Print(logs)
			}
		}
		return nil
	},
}

func init() {
	pipelineCmd.Flags().StringVarP(&projectID, "project", "p", "", "GitLab Project ID (required)")
	pipelineCmd.Flags().StringP("tokenGitlab", "t", "", "GitLab Private TokenGitlab (optional)")
	pipelineCmd.Flags().IntVar(&cnt, "counter", 1, "Show nth latest pipelines")
	pipelineCmd.MarkFlagRequired("project")
	rootCmd.AddCommand(pipelineCmd)
}
