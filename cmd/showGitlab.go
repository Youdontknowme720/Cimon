package cmd

import (
	"fmt"
	utils "github.com/Youdontknowme720/Cimon/gitlab"
	"github.com/Youdontknowme720/Cimon/ui"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var GitlabTree = &cobra.Command{
	Use:   "gtree",
	Short: "Shows pipelines in a tree format",
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
		app := tview.NewApplication()
		ui.InitTree(app, projectID, token)
		return nil
	},
}

func init() {
	GitlabTree.Flags().StringVarP(&projectID, "project", "p", "", "GitLab Project ID (required)")
	GitlabTree.Flags().StringP("tokenGitlab", "t", "", "GitLab Private TokenGitlab (optional)")
	rootCmd.AddCommand(GitlabTree)
}
