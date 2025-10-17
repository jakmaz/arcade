package cmd

import (
	"fmt"

	"github.com/jakmaz/arcade/internal/core"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all games",
	Long:  "List all games",
	Run: func(cmd *cobra.Command, args []string) {
		listGames()
	},
}

func listGames() {
	games := core.AvailableGames()

	fmt.Println("Available Games:")
	fmt.Println()

	for _, game := range games {
		fmt.Printf("  %-12s %s\n", game.ID, game.Description)
	}
}
