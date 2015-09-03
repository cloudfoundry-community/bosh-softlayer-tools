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

	LightStemcellTypeFlag string
	LightStemcellPathFlag string

	VersionFlag              string
	StemcellInfoFilenameFlag string
	InfrastructureFlag       string
	HypervisorFlag           string
	OsNameFlag               string

	NamePatternFlag   string
	LastValidDateFlag string
	ShipItTagFlag     string
}
