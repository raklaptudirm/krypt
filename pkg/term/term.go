package term

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/term"
)

func Error(a ...interface{}) (int, error) {
	return fmt.Fprint(os.Stderr, a...)
}

func Errorf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(os.Stderr, format, a...)
}

func Errorln(a ...interface{}) (int, error) {
	return fmt.Fprintln(os.Stderr, a...)
}

func Pass(format string, a ...interface{}) (pw []byte, err error) {
	fmt.Printf(format, a...)
	pw, err = term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	return
}
