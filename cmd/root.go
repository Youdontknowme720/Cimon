package cmd

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "deincli",
    Short: "Dein CLI Tool",
}

// Execute wird in main.go aufgerufen
func Execute() error {
    return rootCmd.Execute()
}

func init() {
    // Hier Subcommands registrieren
    rootCmd.AddCommand(pipelineCmd)
}
