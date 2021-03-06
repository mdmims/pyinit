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

Help output:
```
Usage: pyinit [OPTIONS] [ARGS]...
CLI to generate gitignore files and other useful python files.
Options:
        --help     Display help message and exit.
        --version  Display version.
        --list     Display the valid gitignore.io API options.

        -a         Create all available files
        -d         Create Dockerfile with basic Python (3.8) template
        -g         Create .gitignore file with default language options (macos, windows, python)
        -f         Create .flake8 file with default settings
        -l         Create License file (MIT)
        -p         Create pyproject.toml file with black formatter default settings for Python 3.8
Arguments:
        TARGETS: Space separated list of gitignore.io language options. [optional]
Examples:
$ pyinit --help
$ pyinit -f -g -l -p -d go python java
$ pyinit -a macos python

```