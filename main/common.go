package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
)

func DisplayCrashDialog(errorMessage string) {
	formattedString := `
Something completely unexpected happened. This is a bug in %s.
Please file this bug : https://github.com/cloudfoundry-community/bosh-softlayer-tools/issues
Tell us that you ran this command:

	%s

this error occurred:

	%s

and this stack trace:

%s
	`

	stackTrace := "\t" + strings.Replace(string(debug.Stack()), "\n", "\n\t", -1)
	println(fmt.Sprintf(formattedString, "bosh-softlayer-tools", strings.Join(os.Args, " "), errorMessage, stackTrace))
}

func HandlePanic() {
	err := recover()

	if err != nil {
		switch err := err.(type) {
		case error:
			DisplayCrashDialog(err.Error())
		case string:
			DisplayCrashDialog(err)
		default:
			DisplayCrashDialog("An unexpected type of error")
		}
	}

	if err != nil {
		os.Exit(1)
	}
}
