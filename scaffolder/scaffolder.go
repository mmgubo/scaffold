package scaffolder

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// ProjectData is passed to every file template.
type ProjectData struct {
	Name        string // raw project name, used for the directory
	Module      string // Go module path
	PackageName string // sanitized Go package identifier
	Year        int
}

type fileSpec struct {
	path    string // relative path inside project root (forward slashes)
	content string // template string using [[ ]] delimiters
}

// Create scaffolds a new project of the given type in a subdirectory named name.
func Create(projectType, name, module string) error {
	valid := map[string]bool{"cli": true, "web": true, "library": true}
	if !valid[projectType] {
		return fmt.Errorf("unknown project type %q — valid types: cli, web, library", projectType)
	}
	if _, err := os.Stat(name); err == nil {
		return fmt.Errorf("directory %q already exists", name)
	}

	data := ProjectData{
		Name:        name,
		Module:      module,
		PackageName: sanitizePackageName(name),
		Year:        time.Now().Year(),
	}

	files := getFiles(projectType, data)
	for _, f := range files {
		fullPath := filepath.Join(name, filepath.FromSlash(f.path))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return fmt.Errorf("mkdir %s: %w", filepath.Dir(fullPath), err)
		}
		rendered, err := render(f.content, data)
		if err != nil {
			return fmt.Errorf("render %s: %w", f.path, err)
		}
		if err := os.WriteFile(fullPath, []byte(rendered), 0644); err != nil {
			return fmt.Errorf("write %s: %w", fullPath, err)
		}
		fmt.Printf("  create  %s\n", filepath.ToSlash(filepath.Join(name, f.path)))
	}

	fmt.Printf("\nCreated %s project %q\n", projectType, name)
	fmt.Printf("Next steps:\n")
	fmt.Printf("  cd %s\n", name)
	switch projectType {
	case "cli":
		fmt.Printf("  go run . --help\n")
	case "web":
		fmt.Printf("  go run .\n")
	case "library":
		fmt.Printf("  go test ./...\n")
	}
	return nil
}

func render(tmplStr string, data ProjectData) (string, error) {
	t, err := template.New("").Delims("[[", "]]").Parse(tmplStr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func getFiles(projectType string, data ProjectData) []fileSpec {
	switch projectType {
	case "cli":
		return cliFiles(data)
	case "web":
		return webFiles(data)
	case "library":
		return libraryFiles(data)
	}
	return nil
}

// sanitizePackageName converts an arbitrary string to a valid Go package identifier.
func sanitizePackageName(name string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(name) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			b.WriteRune(r)
		}
	}
	result := b.String()
	if result == "" {
		return "pkg"
	}
	if result[0] >= '0' && result[0] <= '9' {
		result = "p" + result
	}
	return result
}
