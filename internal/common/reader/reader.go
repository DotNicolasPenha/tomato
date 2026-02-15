package reader

import (
	"encoding/json"
	"fmt"
	"os"

	"com.dotvinci.tm/internal/common/logger"
)

func Json[T any](path string) *T {
	var jsonToReturn T

	bytes, err := os.ReadFile(path)
	if err != nil {
		logger.Fatal(fmt.Sprintf("error reading json at path: %s\n", path))
	}

	if err := json.Unmarshal(bytes, &jsonToReturn); err != nil {
		logger.Fatal(fmt.Sprintf("error unmarshaling json: %s\n", err))
	}

	return &jsonToReturn
}
