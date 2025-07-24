package cmd

import (
	"fmt"
	"github.com/Youdontknowme720/Cimon/gitlab"
	"github.com/Youdontknowme720/Cimon/ui"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var repoUrlShow string
var limitShow int

var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Shows a tview canvas",
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
		app := tview.NewApplication()
		err = ui.StartView(app, repoUrlShow, token)
		if err != nil {
			return fmt.Errorf("UI failed during loading process: %w", err)
		}
		return nil
	},
}

func init() {
	ShowCmd.Flags().StringVarP(&repoUrlShow, "project", "p", "", "GitLab Project ID (required)")
	ShowCmd.Flags().StringP("tokenGithub", "t", "", "GitLab Private TokenGitlab (optional)")
	workflowCmd.Flags().IntVarP(&limitShow, "limiterShow", "o", 10, "Limits the output shown for workflows")
	rootCmd.AddCommand(ShowCmd)
}

