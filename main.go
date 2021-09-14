package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"pyinit/config"
	"strings"
)

const (
	// ignoreURL is the base url for the gitignore API
	ignoreURL      string = "https://www.toptal.com/developers/gitignore/api"
	versionMessage string = "version: 0.0.1"
	listMessage    string = "To get a list of valid targets, run pyinit --list"

	helpMessage string = `
Usage: pyinit [OPTIONS] [ARGS]...
CLI to generate gitignore files and other useful python files.
Options:
	--help     Display help message and exit.
	--version  Display version.
	--list     Display the valid gitignore.io API options.

	-g	   Create .gitignore file with default language options (macos, windows, python)
	-f	   Create .flake8 file with default settings
	-l	   Create License file (MIT)
	-p	   Create pyproject.toml file with black formatter default settings for Python 3.8
Arguments:
	TARGETS: Space separated list of gitignore.io language options.	[optional]
Examples:
$ pyinit --help
$ pyinit -f -g -l -p go python java
`
)

var (
	helpFlag      bool
	versionFlag   bool
	listFlag      bool
	gitignoreFlag bool
	flake8Flag    bool
	licenseFlag   bool
	tomlFlag      bool
	allFlag       bool
)

func main() {
	flag.BoolVar(&helpFlag, "help", false, "Help information")
	flag.BoolVar(&versionFlag, "version", false, "Version number")
	flag.BoolVar(&listFlag, "list", false, "Gitignore API language options list")
	flag.BoolVar(&allFlag, "a", false, "-a")
	flag.BoolVar(&flake8Flag, "f", false, "-f")
	flag.BoolVar(&gitignoreFlag, "g", false, "-g")
	flag.BoolVar(&licenseFlag, "l", false, "-l")
	flag.BoolVar(&tomlFlag, "p", false, "-p")

	flag.Usage = func() {
		fmt.Println(helpMessage)
		os.Exit(1)
	}

	flag.Parse()

	run()
}

func run() {
	if flag.NArg() == 0 && flag.NFlag() == 0 && !(helpFlag || listFlag || versionFlag) {
		fmt.Fprintln(os.Stdout, helpMessage)
		os.Exit(1)
	}

	switch {
	case helpFlag:
		fmt.Fprintln(os.Stdout, helpMessage)
		os.Exit(0)

	case versionFlag:
		fmt.Fprintln(os.Stdout, versionMessage)
		os.Exit(0)

	case listFlag:
		printList(os.Stdout, listMessage)
		os.Exit(0)

	default:
		if licenseFlag {
			create("License")
		}
		if flake8Flag {
			create(".flake8")
		}
		if tomlFlag {
			create("pyproject.toml")
		}
		if gitignoreFlag {
			ignoreList := os.Args[1:]
			makeIgnoreFile(ignoreList, ignoreURL)
		}
		fmt.Fprintln(os.Stdout, helpMessage)
		os.Exit(0)
	}
}

// GetIgnore calls the gitignore API and returns response
func GetIgnore(targets []string, url string) ([]byte, error) {
	targetOptions := buildIgnoreOptions(targets)
	targetURL := strings.Join([]string{url, targetOptions}, "/")

	response, err := http.Get(targetURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return data, nil

}

func buildIgnoreOptions(targetOptions []string) string {
	// create map with default options
	s := make(map[string]bool)
	s["macos"] = true
	s["windows"] = true
	s["python"] = true

	// add default options to input
	for k := range s {
		targetOptions = append(targetOptions, k)
	}

	// dedupe the combined array
	targetOptions = removeDuplicateStrings(targetOptions)

	// build comma separated string of elements
	options := strings.Join(targetOptions, ",")

	return options
}

func removeDuplicateStrings(strSlice []string) []string {
	allKeys := make(map[string]bool)
	var list []string
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// GetList returns the gitignore API response for 'list'
func GetList(url string) ([]byte, error) {
	targetURL := strings.Join([]string{url, "list"}, "/")

	response, err := http.Get(targetURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// WriteToIgnoreFile creates .gitignore config using supplied config
func WriteToIgnoreFile(data []byte, filename string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	ignoreFilePath := filepath.Join(cwd, filename)

	file, err := os.OpenFile(ignoreFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}

	return nil
}

func printList(where io.Writer, url string) {
	data, err := GetList(url)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	fmt.Fprintln(where, string(data))
}

// customIgnoreOptions adds custom ignore options to .gitignore file
func customIgnoreOptions(data []byte) []byte {
	customOptions := []string{".idea", ".vscode"}
	for _, c := range customOptions {
		b := []byte("\n" + c + "/\n")
		data = append(data, b...)
	}
	return data
}

func makeIgnoreFile(targets []string, url string) {
	data, err := GetIgnore(targets, url)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	data = customIgnoreOptions(data)

	err = WriteToIgnoreFile(data, ".gitignore")
	if err != nil {
		fmt.Printf("Error: %s.\n", err)
	}
}

// createDataFile creates specified file in current directory
func createDataFile(data string, filename string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	ignoreFilePath := filepath.Join(cwd, filename)

	file, err := os.OpenFile(ignoreFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	defer file.Close()

	_, err = file.Write([]byte(data))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	if err := file.Sync(); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

// create reads available file mapping and creates designated file
func create(name string) {
	data := config.GetEmbeds(name)
	createDataFile(data, name)
}
