package logger

import (
	"fmt"
	"os"
)

func Info(text string) {
	fmt.Printf("[INFO] %s \n", text)
}
func Error(text string) {
	fmt.Printf("[ERROR] %s \n", text)
}
func Ok(text string) {
	fmt.Printf("[OK] %s \n", text)
}
func Fatal(text string) {
	fmt.Printf("[FATAL] %s \n", text)
	os.Exit(0)
}
