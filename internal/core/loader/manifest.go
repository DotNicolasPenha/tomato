package loader

import (
	"fmt"

	"com.dotvinci.tm/internal/common/logger"
	"com.dotvinci.tm/internal/common/reader"
)

type Manifest struct {
	Port            *int    `json:"port"`
	NameApplication *string `json:"string"`
}

func LoadManifest(path string) Manifest {
	var manifestJson Manifest = reader.Json[Manifest](path)

	if manifestJson.Port == nil {
		logger.Error(fmt.Sprintf("field 'port' is required in manifest of path: %s", path))
	}
	if *manifestJson.Port <= 0 || *manifestJson.Port > 65535 {
		logger.Error(fmt.Sprintf("invalid port number to field 'port' of manifest: %s", path))
	}
	if manifestJson.NameApplication == nil || *manifestJson.NameApplication == "" {
		logger.Error(fmt.Sprintf("field 'nameApplication' is required in manifest of path: %s", path))
	}
	return manifestJson
}
