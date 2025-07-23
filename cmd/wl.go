package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/Youdontknowme720/Cimon/gitlab"
	"github.com/Youdontknowme720/Cimon/github"
)
var repoUrl string
var limit int

var workflowCmd = &cobra.Command{
	Use:   "wf",
	Short: "Shows workflows from Github",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := cmd.Flags().GetString("tokenGithub")
		if err != nil {
			return err
		}

		if token == "" {
			config, err := utils.LoadConfig()
			if err != nil {
				return fmt.Errorf("Couldn't load config %w: ", err)
			}
			token = config.TokenGithub
		}

		if token == "" {
			return fmt.Errorf("TokenGitlab is empty. Please set it using auth command or as a flag.")
		}

		workFlowsResp, err := github.GetWorkflowStatus(repoUrl, limit, token)
		if err != nil{
			return fmt.Errorf("Error while fetching workflows")
		}

		for _, wfr := range workFlowsResp.WorkflowRuns{
			jobsRunRes, _ := github.GetJobRuns(repoUrl, wfr.ID, token)
			for _, job := range jobsRunRes.Jobs{
				job.DisplaySteps()
				fmt.Println("-----------------------------")
			}
		}
		return nil
	},
}

func init() {
	workflowCmd.Flags().StringVarP(&repoUrl, "repo", "r", "", "Github Repo Url requiered")
	workflowCmd.Flags().IntVarP(&limit, "limiter", "l", 5, "Limits the output shown for workflows")
	workflowCmd.Flags().StringP("tokenGithub", "t", "", "GitHub Accesstoken")
	pipelineCmd.MarkFlagRequired("repo")
	rootCmd.AddCommand(workflowCmd)
}
