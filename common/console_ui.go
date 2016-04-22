package common

import "fmt"

type consoleUi struct {
	verbose bool
}

func NewConsoleUi(verbose bool) UI {
	return consoleUi{
		verbose: verbose,
	}
}

func (ui consoleUi) Printf(msg string, args ...interface{}) (int, error) {
	if !ui.verbose {
		return 0, nil
	}

	return fmt.Printf(msg, args...)
}

func (ui consoleUi) Println(args ...interface{}) (int, error) {
	if !ui.verbose {
		return 0, nil
	}

	return fmt.Println(args...)
}
