package configuration

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/monzo/terrors"
)

func ReadApplicationConfig(jsonFilePath string) (*ApplicationConfig, error) {
	errorTags := map[string]string{
		"file": jsonFilePath,
	}

	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		return nil, terrors.Augment(err, "opening config", errorTags)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, terrors.Augment(err, "reading config", errorTags)
	}

	var applicationConfig ApplicationConfig

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err = json.Unmarshal(byteValue, &applicationConfig)
	if err != nil {
		return nil, terrors.Augment(err, "unmarshalling json config", errorTags)
	}
	return &applicationConfig, nil
}
