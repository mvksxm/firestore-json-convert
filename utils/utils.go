package utils

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/mvksxm/firestore-json-convert/models"
)


func ValidatePaths(paths []string, isInput bool, vc chan <- models.StampedPath, wg *sync.WaitGroup) {
	
	defer wg.Done()

	for i, path := range paths {
		if validated, err := ValidatePath(path, isInput); !validated {
			sp := models.StampedPath{
				Id: i,
				Path: path,
				Error: err,	
			}
			vc <- sp
		}
	}
}

func ValidatePath(path string, isInput bool) (bool, string) {
	
	// TODO:
	// Add an initial check to make sure that paths provided are generally valid UNIX paths.
	// Add a check for making sure that path for the input file is not a dir.

	pathToCheck := ""
	if isInput {
		pathToCheck = path
	} else {
		pathToCheck = filepath.Dir(path)
	}

	if _, err := os.Stat(pathToCheck); err != nil {
		return false, err.Error()
	}

	return true, ""
}

func ValidatePayloads(payloads []string) bool {
	return false
}

func ValidatePayload(payload string) bool {
	return false
}