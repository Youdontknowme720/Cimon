package cmd

import (
	"fmt"
	"github.com/Youdontknowme720/Cimon/github"
	"github.com/Youdontknowme720/Cimon/gitlab"
	"github.com/spf13/cobra"
)

var (
	repoUrl string
	limit   int
)

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
		if err != nil {
			return fmt.Errorf("Error while fetching workflows")
		}

		for _, wfr := range workFlowsResp.WorkflowRuns {
			jobs, _ := wfr.GetJobRuns(repoUrl, token)
			for _, job := range jobs {
				job.GetSteps()
			}
			fmt.Println("-----------------------------")
		}
		return nil
	},
}

func init() {
	workflowCmd.Flags().StringVarP(&repoUrl, "repo", "r", "", "Github Repo Url requiered")
	workflowCmd.Flags().IntVarP(&limit, "limiter", "l", 10, "Limits the output shown for workflows")
	workflowCmd.Flags().StringP("tokenGithub", "t", "", "GitHub Accesstoken")
	pipelineCmd.MarkFlagRequired("repo")
	rootCmd.AddCommand(workflowCmd)
}
