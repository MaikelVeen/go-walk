package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-walk <command>",
	Short: "walk is a CLI to interact with geopositional data",
}

func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(NewVisualiseCommand().Command)
}
