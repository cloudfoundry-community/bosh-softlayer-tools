package common

import "fmt"

type consoleUi struct {
}

func NewConsoleUi() UI {
	return consoleUi{}
}

func (ui consoleUi) Printf(msg string, args ...interface{}) (int, error) {
	return fmt.Printf(msg, args...)
}

func (ui consoleUi) Println(args ...interface{}) (int, error) {
	return fmt.Println(args...)
}
