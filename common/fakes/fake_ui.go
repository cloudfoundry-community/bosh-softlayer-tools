package fakes

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

type FakeUi struct {
	Output string

	Msg  string
	Args []interface{}

	PrintfRc  int
	PrintfErr error

	PrintlnRc  int
	PrintlnErr error

	ScanlnRc  int
	ScanlnErr error
	UserInput string
}

func NewFakeUi() *FakeUi {
	return &FakeUi{
		Output: "",
		Msg:    "",
		Args:   make([]interface{}, 0),
	}
}

func (fakeUi *FakeUi) Printf(msg string, args ...interface{}) (int, error) {
	fakeUi.Msg = msg
	fakeUi.Args = args

	fakeUi.Output = fmt.Sprintf(msg, args...)

	return fakeUi.PrintfRc, fakeUi.PrintfErr
}

func (fakeUi *FakeUi) Println(args ...interface{}) (int, error) {
	fakeUi.Args = args

	fakeUi.Output = fmt.Sprintln(args...)

	return fakeUi.PrintlnRc, fakeUi.PrintlnErr
}

func (fakeUi *FakeUi) PrintTable(table *tablewriter.Table) (int, error) {
	table.Render()

	return 0, nil
}

func (fakeUi *FakeUi) PrintfInfo(msg string, args ...interface{}) (int, error) {
	fakeUi.Msg = msg
	fakeUi.Args = args

	fakeUi.Output = fmt.Sprintf(msg, args...)

	return fakeUi.PrintfRc, fakeUi.PrintfErr
}

func (fakeUi *FakeUi) PrintlnInfo(args ...interface{}) (int, error) {
	fakeUi.Args = args

	fakeUi.Output = fmt.Sprintln(args...)

	return fakeUi.PrintlnRc, fakeUi.PrintlnErr
}

func (fakeUi *FakeUi) Scanln(args ...interface{}) (int, error) {
	fakeUi.Args = args

	_, err := fmt.Sscanln(fakeUi.UserInput, args[0])
	if err != nil {
		return 1, err
	}

	return fakeUi.ScanlnRc, fakeUi.ScanlnErr
}

func (fakeUi *FakeUi) NewTableWriter() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	return table
}
