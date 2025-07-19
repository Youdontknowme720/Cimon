
package cmd

import (
	"fmt"
	"gitlab.com/ayan0k0uji-group/Cimon/utils"
	"github.com/spf13/cobra"
)

var projectID, token string
var cnt int

var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Zeigt Pipeline-Status von GitLab",
	Run: func(cmd *cobra.Command, args []string) {
		var limit int
		pipelines, err := utils.GetPiplineStatus(projectID, token)
		if err != nil {
			fmt.Println("Fehler beim Abrufen:", err)
			return
		}
		countSet:= cmd.Flags().Changed("counter")
		if countSet{
			if cnt > len(pipelines){
				limit = len(pipelines)
			}
			limit = cnt
			utils.DisplayPipelines(pipelines[:limit])
		}else{
			utils.DisplayPipelines(pipelines)
		}
	},
}

func init() {
	pipelineCmd.Flags().StringVarP(&projectID, "project", "p", "", "GitLab Project ID (required)")
	pipelineCmd.Flags().StringVarP(&token, "token", "t", "", "GitLab Private Token (required)")
	pipelineCmd.Flags().IntVar(&cnt, "counter", 1, "Show nth latest pipelines")
	pipelineCmd.MarkFlagRequired("project")
	pipelineCmd.MarkFlagRequired("token")
	rootCmd.AddCommand(pipelineCmd)
}
