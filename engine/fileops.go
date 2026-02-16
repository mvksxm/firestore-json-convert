package engine

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)


type FileIO struct { 
	inputPath string
	outputPath string
}

func (fo *FileIO) ReadInput() (map[string]interface{}, error) {
	payload := make(map[string]interface{})

	content, err := os.ReadFile(fo.inputPath)
	if err != nil {
		slog.Warn(fmt.Sprintf("There was an issue with reading an input file - %s. It will be skipped, for now.", fo.inputPath))
		return nil, err
	}

	err = json.Unmarshal(content, &payload)
	if err != nil {
		slog.Warn(fmt.Sprintf("Provided input file - %s contains an invalid json structure!. It will be skipped, for now.", fo.inputPath))
		return nil, err
	}

	return payload, nil
}

func (fo *FileIO) WriteOutput(payload map[string]interface{}) error {

	byteArr, err := json.Marshal(payload)
	if err != nil {
		slog.Warn(fmt.Sprintf("There was an issue with converting a payload of the file %s to a byte array. Err - %s", fo.inputPath, err.Error()))
		return err
	}
	
	err = os.WriteFile(fo.outputPath, byteArr, 0777)
	if err != nil {
		slog.Warn(fmt.Sprintf("There was an issue with writing a payload to the output file - %s. Err - %s", fo.outputPath, err.Error()))
		return err

	}
	return nil
}

func NewFileIO(inputPath string, outputPath string) *FileIO {
	return &FileIO{
		inputPath: inputPath,
		outputPath: outputPath,
	}
}


