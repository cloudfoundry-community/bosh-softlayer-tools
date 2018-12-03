package common

type Options struct {
	CommandFlag string

	HelpFlag     bool
	LongHelpFlag bool

	DryRunFlag bool

	PublicFlag bool

	NameFlag       string
	NoteFlag       string
	PublicNameFlag string
	PublicNoteFlag string
	OsRefCodeFlag  string
	UriFlag        string

	LightStemcellTypeFlag string
	LightStemcellPathFlag string

	VersionFlag              string
	StemcellInfoFilenameFlag string
	InfrastructureFlag       string
	HypervisorFlag           string
	OsNameFlag               string

	StemcellFormatsFlag string
}
