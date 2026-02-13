package commands

import (

	"github.com/mvksxm/firestore-json-convert/engine"
	"github.com/spf13/cobra"
)

type PreviewCommand struct {
	BaseCommand
}

func (pc *PreviewCommand) run(_ *cobra.Command, _ []string) {
	fileArr := pc.generateArrays()
	c := engine.NewMultipleConverterPreview(fileArr)
	c.Run()
} 

func (pc *PreviewCommand) Init() {
	pc.BaseCommand.Init(
		"preview",
		previewCmdDescription,
		pc.run,
	)
}

func NewPreviewCommand() *PreviewCommand {
	pc := new(PreviewCommand)
	pc.Init()
	return pc
}

