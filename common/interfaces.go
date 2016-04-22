package common

import "github.com/olekukonko/tablewriter"

type Printer interface {
	Println(args ...interface{}) (int, error)
	Printf(msg string, args ...interface{}) (int, error)
	PrintTable(table *tablewriter.Table) (int, error)
}

type UI interface {
	Printer
}
