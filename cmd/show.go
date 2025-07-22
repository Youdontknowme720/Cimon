package cmd

import (
	"github.com/spf13/cobra"
	"github.com/Youdontknowme720/Cimon/ui"
)

var ShowCmd = &cobra.Command{
	Use: "show",
	Short: "Shows a tview canvas",
	Run: func(cmd *cobra.Command, args []string){
		ui.StartView()
	},
}

func init(){
	rootCmd.AddCommand(ShowCmd)
}