# scaffold

A zero-dependency project scaffolding tool written in Go. Give it a project type and a name — it creates a ready-to-run directory structure with boilerplate files.

## Installation

```bash
git clone https://github.com/mmgubo/scaffold.git
cd scaffold
go build -o scaffold .
```

Or install directly:

```bash
go install github.com/mmgubo/scaffold@latest
```

## Usage

```
scaffold [options] <type> <name>

Types:
  cli      Command-line application
  web      Web application
  library  Reusable Go library

Options:
  -module string   Go module path (default: project name)
```

## Examples

```bash
scaffold cli mytool
scaffold web myapp --module github.com/alice/myapp
scaffold library mylib
```

## What gets generated

### `cli`

```
myapp/
├── README.md
├── .gitignore
├── go.mod
├── main.go
├── cmd/root.go
├── internal/config/config.go
└── config.json
```

### `web`

```
myapp/
├── README.md
├── .gitignore
├── go.mod
├── main.go
├── handlers/handlers.go
├── internal/config/config.go
└── config.json
```

### `library`

```
mylib/
├── README.md
├── .gitignore
├── go.mod
├── doc.go
├── mylib.go
└── mylib_test.go
```

## Features

- No external dependencies — pure Go stdlib
- Generated projects are also dependency-free and compile immediately
- `--module` flag for custom Go module paths (e.g. `github.com/user/project`)
- Sanitizes project names into valid Go package identifiers automatically

## License

MIT
