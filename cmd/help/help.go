package help

import (
	"fmt"
)

func HelpCommand() {
	commands := []string{"lta", "help"}

	for i := 0; i < len(commands); i++ {
		fmt.Printf("[%d] - %s \n", i, commands[i])
	}
}
