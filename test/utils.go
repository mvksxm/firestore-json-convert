package test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	testPath = "./samples/test"
	resPath = "./samples/res"
)

func populateTestMap(
	testMap *map[string][]map[string]interface{},
	isTestPath bool,
	dirPath string,
	filesArr []os.DirEntry,
	middleVal string,
) {

	filePrefix := "res"
	if isTestPath {
		filePrefix = "test"
	}

	regexPattern := fmt.Sprintf(`^%s_%s_[1-9][0-9]*\.json$`, filePrefix, middleVal)

	for _, f := range filesArr {
		if matched, _ := regexp.MatchString(regexPattern, f.Name()); !matched {
			// fmt.Printf("File name -> %s does not match the -> %s regex pattern\n", f.Name(), regexPattern)
			continue
		}
		fileNameSplitted := strings.Split(strings.Split(f.Name(), ".")[0], "_")
		fileNameMiddle := fileNameSplitted[1]
		fileNameIdx := fileNameSplitted[2]

		if fileNameMiddle == middleVal {
			filePath := filepath.Join(dirPath, f.Name())
			payload := make(map[string]interface{})
			byteVal, err := os.ReadFile(filePath)
			if err != nil {
				log.Fatalf("Failed to read the following file -> %s", f.Name())
			}
			uErr := json.Unmarshal(byteVal, &payload)
			if uErr != nil {
				log.Fatalf("An issue occured, when unmarshalling the following file -> %s",filePath)
			}
			if payloadArr, ok := (*testMap)[fileNameIdx]; ok {
				payloadArr := append(payloadArr, payload)
				(*testMap)[fileNameIdx] = payloadArr
			} else {
				(*testMap)[fileNameIdx] = []map[string]interface{}{payload}
			}
		}
	}
}


func readPayloads(
	testPath string,
	resPath string,
	middleVal string,
	
) map[string][]map[string]interface{} {

	testPayloadsMap := make(map[string][]map[string]interface{})

	testFiles, err := os.ReadDir(testPath)
	if err != nil {
		log.Fatalf("Something went wrong, when reading test files from the dir -> %s. Err: %s", testPath, err.Error())
	}

	populateTestMap(
		&testPayloadsMap,
		true,
		testPath,
		testFiles,
		middleVal,
	)

	resFiles, err := os.ReadDir(resPath)
	if err != nil {
		log.Fatalf("Something went wrong, when reading res files from the dir -> %s. Err: %s", resPath, err.Error())
	}

	populateTestMap(
		&testPayloadsMap,
		false,
		resPath,
		resFiles,
		middleVal,
	)

	if len(testPayloadsMap) == 0 {
		log.Fatalln("Something went wrong. Map with test payloads is empty.")
	}

	filesWithoutPair := []string{}
	for k, v := range testPayloadsMap {
		if len(v) == 1 {
			filesWithoutPair = append(filesWithoutPair, k)
		}
	}

	if len(filesWithoutPair) > 0 {
		idxList := strings.Join(filesWithoutPair, ",")
		log.Fatalf("Files under the following indexes: '%s' do not have a respective pair.", idxList)
	}

	return testPayloadsMap
}


func getPayloads(isEncode bool)  map[string][]map[string]interface{}{
	middleVal := "decode"
	if isEncode {
		middleVal = "encode"
	}
	
	testPayloadsMap := readPayloads(testPath, resPath, middleVal)
	return testPayloadsMap
}