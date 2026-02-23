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

	// Commented strings for testing purposes.
	// mc := engine.NewMultipleConverter([]string{"test_firestore_conversion.json", "test_generic.json"}, []string{"test_firestore_conversion_res.json", "test_generic_res.json"})
	// mc.Run()
}

