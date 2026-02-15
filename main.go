package main

import (
	"os"

	"com.dotvinci.tm/cmd/help"
	"com.dotvinci.tm/cmd/lta"
	"com.dotvinci.tm/internal/common/logger"
)

func main() {
	if len(os.Args) < 2 {
		logger.Fatal("usage: tm <command> or type 'tm help' to show all commands")
	}

	switch os.Args[1] {
	case "help":
		help.HelpCommand()
	case "lta":
		lta.LtaCommand()
	default:
		logger.Fatal("command not found")
	}
}
