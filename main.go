package main

import (
	"os"

	"com.dotvinci.tm/cmd/help"
	"com.dotvinci.tm/cmd/lta"
	"com.dotvinci.tm/internal/common/logger"
)

func main() {
	if len(os.Args) < 1 {
		logger.Fatal("tm needs a one argument: 'tm <argument>'")
	}
	if os.Args[1] == "" {
		logger.Fatal("the command is empty, type 'tm help' to show commands")
	}
	switch os.Args[1] {
	 case "help":
		help.HelpCommand()
	 case "lta":
		lta.LtaCommand()
	}
}
