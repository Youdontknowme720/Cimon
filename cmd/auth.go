package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/spf13/cobra"
	"github.com/Youdontknowme720/Cimon/gitlab"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Save your private gitlab api token in the config file.",
	RunE: func(cmd *cobra.Command, args []string) error {
		config := &utils.Config{}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Please insert your GitLab tokenGitlab: ")
		tokenGitlab, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("couldn't read user input: %w", err)
		}

		config.TokenGitlab = strings.TrimSpace(tokenGitlab)

		fmt.Print("Please insert your Github tokenGithub: ")
		tokenGithub, err := reader.ReadString('\n')

		config.TokenGithub = tokenGithub

		if err != nil {
			return fmt.Errorf("couldn't read user input: %w", err)
		}
		if err := utils.SaveConfig(config); err != nil {
			return fmt.Errorf("couldn't save tokenGitlab to configuration: %w", err)
		}

		fmt.Println("Tokens saved to config file.")
		return nil
	},
}


func init() {
	rootCmd.AddCommand(authCmd)
}
