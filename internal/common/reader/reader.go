package reader

import (
	"encoding/json"
	"fmt"
	"os"

	"com.dotvinci.tm/internal/common/envx"
	"com.dotvinci.tm/internal/common/logger"
)

func Json[T any](path string) *T {
	var jsonToReturn T

	bytes, err := os.ReadFile(path)
	if err != nil {
		logger.Fatal(fmt.Sprintf("error reading json at path: %s\n", path))
	}

	var raw any
	if err := json.Unmarshal(bytes, &raw); err != nil {
		logger.Fatal(fmt.Sprintf("error unmarshaling json: %s\n", err))
	}

	resolved, err := envx.Resolve(raw)
	if err != nil {
		logger.Fatal(fmt.Sprintf("error resolving env values at path %s: %s\n", path, err))
	}

	resolvedBytes, err := json.Marshal(resolved)
	if err != nil {
		logger.Fatal(fmt.Sprintf("error marshaling resolved json at path: %s\n", path))
	}

	if err := json.Unmarshal(resolvedBytes, &jsonToReturn); err != nil {
		logger.Fatal(fmt.Sprintf("error unmarshaling resolved json: %s\n", err))
	}

	return &jsonToReturn
}
