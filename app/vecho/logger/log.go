package logger

import (
	"fmt"
	"os"
)

func Infof(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	fmt.Println()
}

func Fatalf(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
