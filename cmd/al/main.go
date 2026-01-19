package main

import (
	"fmt"
	"io"
	"os"
)

// Version is overridden at build time.
var Version = "dev"

func main() {
	runMain(os.Args, os.Stdout, os.Stderr, os.Exit)
}

func execute(args []string, stdout io.Writer, stderr io.Writer) error {
	cmd := newRootCmd()
	cmd.Version = Version
	cmd.SetVersionTemplate("{{.Version}}\n")
	if len(args) > 1 {
		cmd.SetArgs(args[1:])
	}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	return cmd.Execute()
}

func runMain(args []string, stdout io.Writer, stderr io.Writer, exit func(int)) {
	if err := execute(args, stdout, stderr); err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		exit(1)
	}
}
