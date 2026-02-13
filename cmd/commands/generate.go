package commands

import (
	"github.com/mvksxm/firestore-json-convert/engine"
	"github.com/spf13/cobra"
)

type GenerateCommand struct {
	BaseCommand
	outputPath string
}

func (gc *GenerateCommand) run(_ *cobra.Command, _ []string) {

	fileArr := gc.generateArrays()
	var outputArr []string = nil

	if gc.outputPath != "" {
		outputArr = []string {gc.outputPath}
	}

	c := engine.NewMultipleConverter(fileArr, outputArr)

	c.Run()
} 

func (gc *GenerateCommand) Init() {

	gc.BaseCommand.Init(
		"generate",
		generateCmdDescription,
		gc.run,
	)

	gc.command.Flags().StringVarP(&gc.outputPath, "output", "o", "", "Specify output file path.")
}

func NewGenerateCommand() *GenerateCommand {
	pc := new(GenerateCommand)
	pc.Init()
	return pc
}
