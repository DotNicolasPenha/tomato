package lta

import (
	"fmt"
	"net/http"
	"os"

	"com.dotvinci.tm/internal/common/logger"
	"com.dotvinci.tm/internal/core/distros"
	"com.dotvinci.tm/internal/core/loader"
	"com.dotvinci.tm/internal/core/lta"
	"com.dotvinci.tm/internal/tmd"
)

func LtaCommand() {
	if len(os.Args) < 3 {
		logger.Fatal("lta needs a one argument, type 'tm lta help' to show all subcommands of lta")
	}
	if os.Args[2] == "" {
		logger.Fatal("the command is empty, type 'tm lta help' to show commands")
	}
	cwd, err := os.Getwd()
	if err != nil {
		logger.Fatal(err.Error())
	}
	switch os.Args[2] {
	case "help":
		commands := []string{
			"help - show lta commands",
			"init - initialize a lta",
		}

		for i := 0; i < len(commands); i++ {
			fmt.Printf("[%d] - %s \n", i, commands[i])
		}
	case "init":
		manifestPath := fmt.Sprintf("%s/manifest.json", cwd)
		manifest := loader.LoadManifest(manifestPath)
		lta := lta.Lta{
			Manifest: manifest,
			Distros:  []distros.Distro{},
			Mux:      http.NewServeMux(),
		}
		tmd.ImportsTMD()
		distro, err := distros.Find("tapi-1.0")
		if err != nil {
			logger.Fatal(err.Error())
		}
		lta.PlugDistro(distro)
		lta.ExecuteDistro("tapi-1.0")
		err = lta.Init()
		if err != nil {
			logger.Fatal(err.Error())
		}
	}
}
