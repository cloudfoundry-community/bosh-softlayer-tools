package common

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

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

func (ui consoleUi) PrintTable(table *tablewriter.Table) (int, error) {
	table.Render()
	return 0, nil
}

func (ui consoleUi) PrintfInfo(msg string, args ...interface{}) (int, error) {
	return fmt.Printf(msg, args...)
}

func (ui consoleUi) PrintlnInfo(args ...interface{}) (int, error) {
	return fmt.Println(args...)
}

func (ui consoleUi) Scanln(args ...interface{}) (int, error) {
	return fmt.Scanln(args...)
}

func (ui consoleUi) NewTableWriter() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	return table
}
