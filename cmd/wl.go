package cmd

import (
	"fmt"
	"gitlab.com/ayan0k0uji-group/Cimon/github"

	"github.com/spf13/cobra"
    "gitlab.com/ayan0k0uji-group/Cimon/gitlab"
)
var repoUrl string
var limit int

var workflowCmd = &cobra.Command{
	Use:   "wf",
	Short: "Zeigt Pipeline-Status von GitLab",
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

		workFlows, err := github.GetWorkflowStatus(repoUrl, limit, token)
		if err != nil{
			return fmt.Errorf("Error while fetching workflows")
		}
		fmt.Println(workFlows)
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
