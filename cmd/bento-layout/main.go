// Command bento-layout creates and checks a standard folder layout on a drive.
package main

import (
	"os"

	"github.com/EduardKav/bento-layout/internal/cli"
)

func main() {
	os.Exit(cli.Execute(os.Args[1:], os.Stdout, os.Stderr))
}
