package commands

import "github.com/spf13/cobra"

type BaseCommand struct {
	command cobra.Command
	// payload string
	file string
}

func (bc *BaseCommand) generateArrays() []string {

	// var payloadArr []string = nil
	var fileArr []string = nil

	// if bc.payload != "" {
	// 	payloadArr = []string {bc.payload}
	// }

	if bc.file != "" {
		fileArr = []string {bc.file}
	}

	return fileArr
}

func (bc *BaseCommand) Init(
	name string, 
	shortDesc string,
	// longDesc string, 
	runFunc func(cmd *cobra.Command, args []string),
) {

	// Populate basic args of the command.
	bc.command.Use = name
	bc.command.Short = shortDesc
	// bc.c.Long = longDesc	
	bc.command.Run = runFunc	

	// Global CLI args
	// bc.command.Flags().StringVarP(&bc.payload, "payload", "p", "", "Specify inline json payload to be converted.")
	bc.command.Flags().StringVarP(&bc.file, "file", "f", "", "Specify path to the file that contain json structure to be converted.")
}

func (bc *BaseCommand) GetCommand() *cobra.Command {
	return &bc.command
}
