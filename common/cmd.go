package common

type Options struct {
	CommandFlag string

	HelpFlag     bool
	LongHelpFlag bool

	DryRunFlag bool

	NameFlag      string
	NoteFlag      string
	OsRefCodeFlag string
	UriFlag       string
}
