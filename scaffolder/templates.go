package scaffolder

// ─── common ──────────────────────────────────────────────────────────────────

const tmplGitignore = `# Binaries
[[.Name]]
*.exe
*.dll
*.so
*.dylib

# Test binary, built with go test -c
*.test

# Coverage output
*.out
coverage.html

# Dependency directories
vendor/

# IDE / editor
.idea/
.vscode/
*.swp
*.swo
*~

# OS artefacts
.DS_Store
Thumbs.db
`

const tmplGoMod = `module [[.Module]]

go 1.21
`

// ─── cli ─────────────────────────────────────────────────────────────────────

const tmplCLIReadme = `# [[.Name]]

A command-line tool built with Go.

## Installation

` + "```bash" + `
go install [[.Module]]@latest
` + "```" + `

Or build from source:

` + "```bash" + `
git clone <repo-url>
cd [[.Name]]
go build -o [[.Name]] .
` + "```" + `

## Usage

` + "```" + `
[[.Name]] [options]

Options:
  -config string   path to config file (default "config.json")
  -verbose         enable verbose output
` + "```" + `

## Configuration

Copy and edit ` + "`config.json`" + ` to customise the tool.

## License

MIT © [[.Year]]
`

const tmplCLIMain = `package main

import (
	"fmt"
	"os"

	"[[.Module]]/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
`

const tmplCLICmdRoot = `package cmd

import (
	"flag"
	"fmt"
	"os"

	"[[.Module]]/internal/config"
)

// Execute is the entry point called from main.
func Execute() error {
	fs := flag.NewFlagSet("[[.Name]]", flag.ContinueOnError)
	verbose := fs.Bool("verbose", false, "enable verbose output")
	cfgPath := fs.String("config", "config.json", "path to config file")

	if err := fs.Parse(os.Args[1:]); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}

	cfg, err := config.Load(*cfgPath)
	if err != nil && *verbose {
		fmt.Fprintf(os.Stderr, "warning: could not load config (%v), using defaults\n", err)
	}

	if *verbose {
		fmt.Printf("starting [[.Name]] (env=%s)\n", cfg.Env)
	}

	// TODO: implement your CLI logic here.
	fmt.Println("Hello from [[.Name]]!")
	_ = cfg
	return nil
}
`

const tmplCLIConfigGo = `package config

import (
	"encoding/json"
	"os"
)

// Config holds application configuration.
type Config struct {
	Env      string ` + "`" + `json:"env"` + "`" + `
	LogLevel string ` + "`" + `json:"log_level"` + "`" + `
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		Env:      "development",
		LogLevel: "info",
	}
}

// Load reads configuration from a JSON file.
// If the file cannot be opened or decoded, defaults are returned alongside the error.
func Load(path string) (*Config, error) {
	cfg := Default()
	f, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
`

const tmplCLIConfigJSON = `{
  "env": "development",
  "log_level": "info"
}
`

func cliFiles(_ ProjectData) []fileSpec {
	return []fileSpec{
		{path: "README.md", content: tmplCLIReadme},
		{path: ".gitignore", content: tmplGitignore},
		{path: "go.mod", content: tmplGoMod},
		{path: "main.go", content: tmplCLIMain},
		{path: "cmd/root.go", content: tmplCLICmdRoot},
		{path: "internal/config/config.go", content: tmplCLIConfigGo},
		{path: "config.json", content: tmplCLIConfigJSON},
	}
}

// ─── web ─────────────────────────────────────────────────────────────────────

const tmplWebReadme = `# [[.Name]]

A web application built with Go's standard library.

## Running

` + "```bash" + `
go run .
` + "```" + `

Open http://localhost:8080.

## Configuration

Edit ` + "`config.json`" + ` or point to a different file with ` + "`-config`" + `:

` + "```bash" + `
go run . -config /etc/[[.Name]]/config.json
` + "```" + `

## Project Structure

` + "```" + `
[[.Name]]/
├── main.go                 Entry point
├── handlers/
│   └── handlers.go         HTTP request handlers
├── internal/
│   └── config/
│       └── config.go       Configuration loader
└── config.json             Default configuration
` + "```" + `

## License

MIT © [[.Year]]
`

const tmplWebMain = `package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"[[.Module]]/handlers"
	"[[.Module]]/internal/config"
)

func main() {
	cfgPath := flag.String("config", "config.json", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: using defaults (%v)\n", err)
		cfg = config.Default()
	}

	mux := http.NewServeMux()
	handlers.Register(mux)

	log.Printf("[[.Name]] listening on %s (env=%s)", cfg.Addr, cfg.Env)
	if err := http.ListenAndServe(cfg.Addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
`

const tmplWebHandlers = `package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Register mounts all application routes onto mux.
func Register(mux *http.ServeMux) {
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/healthz", handleHealth)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "Hello from [[.Name]]!")
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
`

const tmplWebConfigGo = `package config

import (
	"encoding/json"
	"os"
)

// Config holds application configuration.
type Config struct {
	Env  string ` + "`" + `json:"env"` + "`" + `
	Addr string ` + "`" + `json:"addr"` + "`" + `
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		Env:  "development",
		Addr: ":8080",
	}
}

// Load reads configuration from a JSON file.
// If the file cannot be opened or decoded, defaults are returned alongside the error.
func Load(path string) (*Config, error) {
	cfg := Default()
	f, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
`

const tmplWebConfigJSON = `{
  "env": "development",
  "addr": ":8080"
}
`

func webFiles(_ ProjectData) []fileSpec {
	return []fileSpec{
		{path: "README.md", content: tmplWebReadme},
		{path: ".gitignore", content: tmplGitignore},
		{path: "go.mod", content: tmplGoMod},
		{path: "main.go", content: tmplWebMain},
		{path: "handlers/handlers.go", content: tmplWebHandlers},
		{path: "internal/config/config.go", content: tmplWebConfigGo},
		{path: "config.json", content: tmplWebConfigJSON},
	}
}

// ─── library ─────────────────────────────────────────────────────────────────

const tmplLibReadme = `# [[.Name]]

> TODO: one-line description.

## Installation

` + "```bash" + `
go get [[.Module]]
` + "```" + `

## Usage

` + "```go" + `
import "[[.Module]]"

msg := [[.PackageName]].Greet("World")
fmt.Println(msg) // Hello, World!
` + "```" + `

## License

MIT © [[.Year]]
`

const tmplLibDoc = `// Package [[.PackageName]] provides ...
//
// TODO: describe what this package does.
package [[.PackageName]]
`

const tmplLibGo = `package [[.PackageName]]

// Greet returns a greeting message for the given name.
// TODO: replace with your library's actual API.
func Greet(name string) string {
	return "Hello, " + name + "!"
}
`

const tmplLibTestGo = `package [[.PackageName]]

import "testing"

func TestGreet(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"World", "Hello, World!"},
		{"Go", "Hello, Go!"},
	}
	for _, tt := range tests {
		got := Greet(tt.input)
		if got != tt.want {
			t.Errorf("Greet(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
`

func libraryFiles(data ProjectData) []fileSpec {
	return []fileSpec{
		{path: "README.md", content: tmplLibReadme},
		{path: ".gitignore", content: tmplGitignore},
		{path: "go.mod", content: tmplGoMod},
		{path: "doc.go", content: tmplLibDoc},
		{path: data.PackageName + ".go", content: tmplLibGo},
		{path: data.PackageName + "_test.go", content: tmplLibTestGo},
	}
}
