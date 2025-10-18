package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/jakmaz/arcade/internal/theme"
	"github.com/jakmaz/arcade/internal/ui/styles"
	"github.com/spf13/cobra"
)

var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Theme management commands",
	Long:  `Commands for managing arcade themes`,
}

var listThemesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available themes",
	Long:  `List all available themes in the arcade`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize themes
		if err := theme.Initialize(); err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing themes: %v\n", err)
			os.Exit(1)
		}

		themes := theme.ListThemes()
		currentTheme := theme.GetCurrentTheme()

		lipgloss.Println("Available themes:")
		for _, name := range themes {
			if currentTheme != nil && name == currentTheme.Name() {
				lipgloss.Printf("* %s (current)\n", name)
			} else {
				lipgloss.Printf("  %s\n", name)
			}
		}
	},
}

var setThemeCmd = &cobra.Command{
	Use:   "set [theme-name]",
	Short: "Set the current theme",
	Long:  `Set the current theme for the arcade`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		themeName := args[0]

		// Initialize themes
		if err := theme.Initialize(); err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing themes: %v\n", err)
			os.Exit(1)
		}

		if err := theme.SetCurrentTheme(themeName); err != nil {
			fmt.Fprintf(os.Stderr, "Error setting theme: %v\n", err)
			os.Exit(1)
		}

		// Refresh styles to use new theme
		styles.RefreshStyles()

		lipgloss.Printf("Theme set to: %s\n", themeName)
	},
}

var previewThemeCmd = &cobra.Command{
	Use:   "preview [theme-name]",
	Short: "Preview a theme",
	Long:  `Preview how a theme looks with sample colors`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		themeName := args[0]

		// Initialize themes
		if err := theme.Initialize(); err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing themes: %v\n", err)
			os.Exit(1)
		}

		// Get the theme
		themeObj, exists := theme.GetTheme(themeName)
		if !exists {
			fmt.Fprintf(os.Stderr, "Theme '%s' not found\n", themeName)
			os.Exit(1)
		}

		// Create styles with this theme
		themeStyles := theme.NewStylesWithTheme(themeObj)

		lipgloss.Printf("Theme: %s\n\n", themeObj.Name())

		// Show UI colors
		lipgloss.Println("UI Colors:")
		lipgloss.Printf("  Primary:   %s\n", themeStyles.TitleStyle().Render("Sample Primary Text"))
		lipgloss.Printf("  Secondary: %s\n", themeStyles.MenuItemStyle().Render("Sample Secondary Text"))
		lipgloss.Printf("  Accent:    %s\n", themeStyles.SelectedItemStyle().Render("Sample Accent Text"))
		lipgloss.Printf("  Success:   %s\n", themeStyles.SuccessStyle().Render("Sample Success Text"))
		lipgloss.Printf("  Warning:   %s\n", themeStyles.WarningStyle().Render("Sample Warning Text"))
		lipgloss.Printf("  Error:     %s\n", themeStyles.ErrorStyle().Render("Sample Error Text"))

		lipgloss.Println("\nGame Colors:")
		lipgloss.Printf("  Player 1:    %s\n", themeStyles.Player1Style().Render("●"))
		lipgloss.Printf("  Player 2:    %s\n", themeStyles.Player2Style().Render("●"))
		lipgloss.Printf("  Snake Body:  %s\n", themeStyles.SnakeStyle().Render("●"))
		lipgloss.Printf("  Snake Head:  %s\n", themeStyles.SnakeHeadStyle().Render("◉"))
		lipgloss.Printf("  Food:        %s\n", themeStyles.FoodStyle().Render("◆"))

		lipgloss.Println("\nChess Pieces:")
		lipgloss.Printf("  White: %s\n", themeStyles.WhitePieceStyle().Render("♔ ♕ ♖ ♗ ♘ ♙"))
		lipgloss.Printf("  Black: %s\n", themeStyles.BlackPieceStyle().Render("♚ ♛ ♜ ♝ ♞ ♟"))

		lipgloss.Println("\nTetris Pieces:")
		lipgloss.Printf("  I: %s  ", themeStyles.TetrisPieceStyle("I").Render("████"))
		lipgloss.Printf("O: %s  ", themeStyles.TetrisPieceStyle("O").Render("██"))
		lipgloss.Printf("T: %s  ", themeStyles.TetrisPieceStyle("T").Render("███"))
		lipgloss.Printf("S: %s\n", themeStyles.TetrisPieceStyle("S").Render("██"))
		lipgloss.Printf("  Z: %s  ", themeStyles.TetrisPieceStyle("Z").Render("██"))
		lipgloss.Printf("J: %s  ", themeStyles.TetrisPieceStyle("J").Render("███"))
		lipgloss.Printf("L: %s\n", themeStyles.TetrisPieceStyle("L").Render("███"))
	},
}

func init() {
	themeCmd.AddCommand(listThemesCmd)
	themeCmd.AddCommand(setThemeCmd)
	themeCmd.AddCommand(previewThemeCmd)
	rootCmd.AddCommand(themeCmd)
}
