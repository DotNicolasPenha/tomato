package loader

import (
	"fmt"

	"com.dotvinci.tm/internal/common/logger"
	"com.dotvinci.tm/internal/common/reader"
)

type Manifest struct {
	Port            *int    `json:"port"`
	NameApplication *string `json:"nameApplication"`
}

func LoadManifest(path string) *Manifest {
	manifestJson := reader.Json[Manifest](path)

	var errs []string

	if manifestJson.Port == nil {
		errs = append(errs, fmt.Sprintf("field 'port' is required in manifest of path: %s", path))
	} else if *manifestJson.Port <= 0 || *manifestJson.Port > 65535 {
		errs = append(errs, fmt.Sprintf("invalid port number in field 'port' of manifest: %s", path))
	}


	if manifestJson.NameApplication == nil || *manifestJson.NameApplication == "" {
		errs = append(errs, fmt.Sprintf("field 'nameApplication' is required in manifest of path: %s", path))
	}

	if len(errs) > 0 {
		for _, e := range errs {
			logger.Error(e)
		}
		logger.Fatal("The manifest cannot be loaded because it has invalid fields.")
	}

	return manifestJson
}
