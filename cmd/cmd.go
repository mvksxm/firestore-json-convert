package cmd

import (
	"github.com/mvksxm/firestore-json-convert/cmd/commands"
	"github.com/spf13/cobra"
)


const (
	rootCmdDescription = "CLI tool for converting JSON files to the Firestore API compatible ones and the other way around."
)

var RootCmd = &cobra.Command{
	Use:   "fic",
	Short: rootCmdDescription,
	// Long: `Will be added later.`
}


func init() {

	// Register commands
	previewCmd := commands.NewPreviewCommand().GetCommand()
	generateCmd := commands.NewGenerateCommand().GetCommand()


	// Add commands to the root cmd
	RootCmd.AddCommand(previewCmd, generateCmd)
}