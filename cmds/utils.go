package cmds

func EqualOptions(opts1 Options, opts2 Options) bool {
	return opts1.Help == opts2.Help &&
		opts1.Verbose == opts2.Verbose &&
		opts1.DryRun == opts2.DryRun &&
		opts1.Latest == opts2.Latest &&
		opts1.Packages == opts2.Packages &&
		opts1.PackageOptions == opts2.PackageOptions
}
