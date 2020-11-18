package comparator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

var (
	ErrReadingFile     = fmt.Errorf("An error occured when trying to read a file.")
	ErrInvalidJsonFile = fmt.Errorf("The json file is not valid !")
)

func getfileLoaderFunc(basepath string) LoadExternalResourceFunction {
	return func(path string) (map[string]interface{}, error) {
		filename := filepath.Clean(basepath + path + ".json")
		bytes, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("%w %s", ErrReadingFile, err.Error())
		}

		data := map[string]interface{}{}
		err = json.Unmarshal(bytes, &data)
		if err != nil {
			return nil, fmt.Errorf("%w %s", ErrInvalidJsonFile, err.Error())
		}

		return data, nil
	}
}
