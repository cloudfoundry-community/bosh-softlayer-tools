package cmds

type Options struct {
	Help bool `short:"h" long:"help" description:"Show help information"`

	Verbose bool `short:"v" long:"verbose" description:"Show verbose debug information"`

	DryRun bool `long:"dry-run" description:"Runs command in a dry-run fashion (i.e., fake)"`

	Latest uint `long:"latest" description:"The latest task number to use"`

	Packages       string `long:"packages" description:"List SL packages"`
	PackageOptions string `long:"package-options" description:"List SL package options"`

	Username string `long:"username" short:"u" description:"the username for login in"`
	Password string `long:"password" short:"p" description:"the password for login in"`

	Target string `long:"target" short:"t" description:"the target URL"`

	Deployment string `long:"Deployment" short:"d" description:"The deployment file"`
}

type Command interface {
	Name() string
	Description() string
	Usage() string
	Options() Options

	Validate() (bool, error)
	Execute(args []string) (int, error)
}
