package engine

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"github.com/mvksxm/firestore-json-convert/models"
	"github.com/mvksxm/firestore-json-convert/utils"
)

type MultipleConverter struct {
	isPreview bool
	inputPaths []string
	outputPaths []string
} 

func (mc *MultipleConverter) initValMap(valChannel chan models.StampedPath) map[int][]models.StampedPath {

	valMap := map[int][]models.StampedPath{}
	for i := range mc.inputPaths {
		valMap[i] = []models.StampedPath{}
	}
	
	for sp := range valChannel {
		paths := append(valMap[sp.Id], sp)
		valMap[sp.Id] = paths
	}

	return valMap
}

func (mc *MultipleConverter) checkPathDuplicates() error {

	dSet := map[string]interface{} {}

	for _, inputPath := range mc.inputPaths {
		dSet[inputPath] = nil
	}

	if len(dSet) != len(mc.inputPaths) {
		return errors.New("Input paths (CLI arg '-f') are not unique.")
	}

	clear(dSet)

	for _, outputPath := range mc.outputPaths {
		dSet[outputPath] = nil
	}

	if len(dSet) != len(mc.outputPaths) {
		return errors.New("Output paths (CLI arg '-o') are not unique.")
	}

	return nil
}

func (mc *MultipleConverter) validate() error {

	if mc.inputPaths == nil {
		return errors.New("Input paths (-f CLI argument) can't be empty!")
	}

	if !mc.isPreview && len(mc.outputPaths) == 0 {
		return errors.New("In case, if mode is 'generate' ('generate' CLI argument), output file paths should be specified!")
	}

	if dpError := mc.checkPathDuplicates(); dpError != nil {
		return dpError
	}

	if !mc.isPreview && len(mc.outputPaths) != len(mc.inputPaths) {
		return errors.New(
			`In generate mode ('generate' CLI argument), 
			amount of input paths (-f) should be equal to the amount of output paths (-o)`,
		)
	}

	 
	valChannel := make(chan models.StampedPath, len(mc.inputPaths) * 2)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go utils.ValidatePaths(mc.inputPaths, true, valChannel, wg)

	if !mc.isPreview {
		wg.Add(1)
		go utils.ValidatePaths(mc.outputPaths, false, valChannel, wg)
	} 

	wg.Wait()
	close(valChannel)

	// Init valMap
	valMap := mc.initValMap(valChannel)

	validInput := []string{}
	validOutput := []string{}

	// Iteration through the map to notify, which paths were not validated and populate valid paths, respectively.
	for idx, spArr := range valMap {

		iPath := mc.inputPaths[idx]

		var oPath string
		if !mc.isPreview {
			oPath = mc.outputPaths[idx]
		}

		if len(spArr) > 0 {
			for _, sp := range spArr {
				fmt.Printf(
					"Following path - %s is invalid and will be skipped alongside its respective pair. Invalidity reason - %s \n", 
					sp.Path,
					sp.Error,
				)
			}
		} else {
			// Valid block
			validInput = append(validInput, iPath)

			if !mc.isPreview {
				validOutput = append(validOutput, oPath)
			}
		}
	}

	if len(validInput) == 0 {
		return errors.New("All of the file path pairs provided are invalid! Exiting...")
	}

	mc.inputPaths = validInput
	mc.outputPaths = validOutput

	return nil
}  


func (mc *MultipleConverter) Run() {

	if err := mc.validate(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}


// type Converter struct {
// 	isPreview bool
// 	payload ...
// }


func NewMultipleConverter(
	inputPaths []string, 
	outputPaths []string,
) *MultipleConverter {
	
	return &MultipleConverter{
		isPreview: false,
		inputPaths: inputPaths,
		outputPaths: outputPaths, 
	}
}

func NewMultipleConverterPreview(
	inputPaths []string, 
) *MultipleConverter {
	
	return &MultipleConverter{
		isPreview: true,
		inputPaths: inputPaths,
		outputPaths: nil, 
	}
}