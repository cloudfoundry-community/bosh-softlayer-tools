package common

type Printer interface {
	Println(args ...interface{}) (int, error)
	Printf(msg string, args ...interface{}) (int, error)
}

type UI interface {
	Printer
}
