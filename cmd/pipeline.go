
package cmd

import (
	"fmt"
	"gitlab.com/ayan0k0uji-group/Cimon/utils"
	"github.com/spf13/cobra"
)

var projectID, token string

var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Zeigt Pipeline-Status von GitLab",
	Run: func(cmd *cobra.Command, args []string) {
		pipelines, err := utils.GetPiplineStatus(projectID, token)
		if err != nil {
			fmt.Println("Fehler beim Abrufen:", err)
			return
		}
		for _, p := range pipelines {
			fmt.Printf("ID: %d | Status: %s | Created: %s | URL: %s\n", p.ID, p.Status, p.Created, p.WebURL)
		}
	},
}

func init() {
	pipelineCmd.Flags().StringVarP(&projectID, "project", "p", "", "GitLab Project ID (required)")
	pipelineCmd.Flags().StringVarP(&token, "token", "t", "", "GitLab Private Token (required)")
	pipelineCmd.MarkFlagRequired("project")
	pipelineCmd.MarkFlagRequired("token")
	rootCmd.AddCommand(pipelineCmd)
}
