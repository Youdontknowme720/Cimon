package main

import (
    "fmt"
    "github.com/spf13/cobra"
)

func main() {
    var name string

    // Root Command
    var rootCmd = &cobra.Command{
        Use:   "app",
        Short: "Eine simple Cobra CLI",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Willkommen bei der Cobra CLI! Nutze 'app hello --name deinName'")
        },
    }

    // Hello Command
    var helloCmd = &cobra.Command{
        Use:   "hello",
        Short: "Sagt Hallo",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Printf("Hallo, %s!\n", name)
        },
    }

    // Flag für helloCmd
    helloCmd.Flags().StringVarP(&name, "name", "n", "Welt", "Name zum Grüßen")

    // Befehl anhängen
    rootCmd.AddCommand(helloCmd)

    // Ausführen
    rootCmd.Execute()
}
