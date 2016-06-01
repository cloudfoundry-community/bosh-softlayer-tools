package common

import (
	"github.com/olekukonko/tablewriter"
)

type Printer interface {
	Println(args ...interface{}) (int, error)
	Printf(msg string, args ...interface{}) (int, error)
	PrintTable(table *tablewriter.Table) (int, error)
	PrintfInfo(msg string, args ...interface{}) (int, error)
	PrintlnInfo(args ...interface{}) (int, error)
	Scanln(args ...interface{}) (int, error)
	NewTableWriter() *tablewriter.Table
}

type UI interface {
	Printer
}
