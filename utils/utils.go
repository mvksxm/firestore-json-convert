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
	
	// Check if paths provided are not dirs
	pathInfo, err := os.Stat(path)
	if err == nil && pathInfo.IsDir() {
		return false, "path provided is an existing directory"
	}

	// In case, if validation errored for input path - return error.
	if err != nil && isInput {
		return false, err.Error()
	}

	// For the output path, we just need to make sure that the parent dir of the path provided is valid. 
	if _, err := os.Stat(filepath.Dir(path)); err != nil && !isInput {
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