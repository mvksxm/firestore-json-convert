package main

import (
	"log"

	"github.com/mvksxm/firestore-json-convert/cmd"
	// "github.com/mvksxm/firestore-json-convert/engine"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalln(err.Error())
	}

	// mc := engine.NewMultipleConverter([]string{"???"}, []string{"???"})
	// mc.Run()
}

