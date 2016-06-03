package common

import (
	"github.com/olekukonko/tablewriter"
	"os"
)

type defaultPrinter struct {
	Ui      UI
	Verbose bool
}

func NewDefaultPrinter(ui UI, verbose bool) Printer {
	return defaultPrinter{
		Ui:      ui,
		Verbose: verbose,
	}
}

func (p defaultPrinter) Printf(msg string, args ...interface{}) (int, error) {
	if p.Verbose {
		return p.Ui.Printf(msg, args)
	}

	return 0, nil
}

func (p defaultPrinter) Println(args ...interface{}) (int, error) {
	if p.Verbose {
		return p.Ui.Println(args)
	}

	return 0, nil
}

func (p defaultPrinter) PrintTable(table *tablewriter.Table) (int, error) {
	table.Render()

	return 0, nil
}

func (p defaultPrinter) PrintfInfo(msg string, args ...interface{}) (int, error) {
	return p.Ui.PrintfInfo(msg, args)
}

func (p defaultPrinter) PrintlnInfo(args ...interface{}) (int, error) {
	return p.Ui.PrintlnInfo(args)
}

func (p defaultPrinter) Scanln(args ...interface{}) (int, error) {
	return p.Ui.Scanln(args...)
}

func (p defaultPrinter) NewTableWriter() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	return table
}
