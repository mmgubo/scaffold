package main

import (
	"flag"
	"fmt"
	"os"

	"scaffold/scaffolder"
)

func main() {
	module := flag.String("module", "", "Go module path (default: project name)")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "scaffold — project scaffolding tool")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, "  scaffold [options] <type> <name>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Types:")
		fmt.Fprintln(os.Stderr, "  cli      Command-line application")
		fmt.Fprintln(os.Stderr, "  web      Web application")
		fmt.Fprintln(os.Stderr, "  library  Reusable Go library")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Examples:")
		fmt.Fprintln(os.Stderr, "  scaffold cli mytool")
		fmt.Fprintln(os.Stderr, "  scaffold web myapp --module github.com/alice/myapp")
		fmt.Fprintln(os.Stderr, "  scaffold library mylib")
	}
	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	projectType := flag.Arg(0)
	projectName := flag.Arg(1)

	modulePath := *module
	if modulePath == "" {
		modulePath = projectName
	}

	if err := scaffolder.Create(projectType, projectName, modulePath); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
