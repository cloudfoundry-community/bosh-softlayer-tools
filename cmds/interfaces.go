package cmds

type Options struct {
	Help bool `short:"h" long:"help" description:"Show help information"`

	Verbose bool `short:"v" long:"verbose" description:"Show verbose debug information"`

	DryRun bool `long:"dry-run" description:"Runs command in a dry-run fashion (i.e., fake)"`

	Latest uint `long:"latest" description:"The latest task number to use"`

	Packages       string `long:"packages" description:"List SL packages"`
	PackageOptions string `long:"package-options" description:"List SL package options"`
}

type Command interface {
	Name() string
	Description() string
	Usage() string
	Options() Options

	Validate() (bool, error)
	Execute(args []string) (int, error)
}
