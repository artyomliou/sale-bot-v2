package utils

import (
	"fmt"
	"log"
	"os"
)

func NewModuleLogger(module string) *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("[%s] ", module), log.LstdFlags)
}
