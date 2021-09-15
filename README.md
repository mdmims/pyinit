## Pyinit
Helpful CLI application that can generate and create applicable default settings for the following files that are useful for Python projects.
Optionally, arguments can be passed to specify the exact language options desired to ignore for the `.gitignore` file. This is achieved 
through the [gitignore api](https://docs.gitignore.io/use/api).
- .gitignore
   - default language options: macos, windows, python
- .flake8
  - defaults to max line length (160) and max-complexity (10)
- pyproject.toml
  - includes `black` formatter settings for max line length (160) and Python Version 3.8
- License
  - Generates MIT License file
- Dockerfile
  - `python:3.8-slim-buster` basic template
- .dockerignore
  - includes various defaults

## Installation
Install Go 
```
https://golang.org/doc/install
```
Clone repository 
```
git clone git@gitlab.verizon.com:data_science_cary_nc/pyinit.git
```
Compile binary within project
```
go build
```
Install binary into `$GOPATH/pkg`
```
go install
```

## Usage
```
pyinit --help
```