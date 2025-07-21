package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
    "gitlab.com/ayan0k0uji-group/Cimon/gitlab"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Save your private gitlab api token in the config file.",
	RunE: func(cmd *cobra.Command, args []string) error {
		config := &utils.Config{}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Please insert your GitLab token: ")
		token, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("couldn't read user input: %w", err)
		}

		config.Token = strings.TrimSpace(token)

		if err := utils.SaveConfig(config); err != nil {
			return fmt.Errorf("couldn't save token to configuration: %w", err)
		}

		fmt.Println("Token saved to config file.")
		return nil
	},
}


func init() {
	rootCmd.AddCommand(authCmd)
}
