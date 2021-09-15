package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"pyinit/config"
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

	-a     Create all available files
	-d     Create Dockerfile with basic Python (3.8) template
	-g	   Create .gitignore file with default language options (macos, windows, python)
	-f	   Create .flake8 file with default settings
	-l	   Create License file (MIT)
	-p	   Create pyproject.toml file with black formatter default settings for Python 3.8
Arguments:
	TARGETS: Space separated list of gitignore.io language options.	[optional]
Examples:
$ pyinit --help
$ pyinit -f -g -l -p -d go python java
$ pyinit -a macos python`
)

var (
	helpFlag       bool
	versionFlag    bool
	listFlag       bool
	gitignoreFlag  bool
	flake8Flag     bool
	licenseFlag    bool
	tomlFlag       bool
	dockerfileFlag bool
	allFlag        bool
)

func main() {
	flag.BoolVar(&helpFlag, "help", false, "Help information")
	flag.BoolVar(&versionFlag, "version", false, "Version number")
	flag.BoolVar(&listFlag, "list", false, "Gitignore API language options list")
	flag.BoolVar(&allFlag, "a", false, "Create all files")
	flag.BoolVar(&dockerfileFlag, "d", false, "Create Dockerfile")
	flag.BoolVar(&flake8Flag, "f", false, "Create .flake8")
	flag.BoolVar(&gitignoreFlag, "g", false, "Create .gitignore")
	flag.BoolVar(&licenseFlag, "l", false, "Create License")
	flag.BoolVar(&tomlFlag, "p", false, "Create pyproject.toml")

	flag.Usage = func() {
		fmt.Println(helpMessage)
		os.Exit(1)
	}

	flag.Parse()

	run()
}

func run() {
	if flag.NArg() == 0 && flag.NFlag() == 0 && !(helpFlag || listFlag || versionFlag) {
		fmt.Println(helpMessage)
		os.Exit(1)
	}

	switch {
	case helpFlag:
		fmt.Println(helpMessage)
		os.Exit(0)

	case versionFlag:
		fmt.Println(versionMessage)
		os.Exit(0)

	case listFlag:
		printApiList(ignoreURL)
		os.Exit(0)

	default:
		if allFlag {
			createAllFiles()
			ignoreList := flag.Args()
			makeIgnoreFile(ignoreList)
			os.Exit(0)
		}
		if licenseFlag {
			create("License")
		}
		if flake8Flag {
			create(".flake8")
		}
		if tomlFlag {
			create("pyproject.toml")
		}
		if dockerfileFlag {
			create("Dockerfile")
		}
		if gitignoreFlag {
			ignoreList := flag.Args()
			makeIgnoreFile(ignoreList)
		}
	}
}

// getIgnore calls the gitignore API and returns response
func getIgnore(targets []string, url string) ([]byte, error) {
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

// buildIgnoreOptions builds a comma separated string of desired language options for gitignore api
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

// removeDuplicateStrings removes duplicate string values from string slices
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

// getList returns the gitignore API response for 'list'
func getList(url string) ([]byte, error) {
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

// printApiList retrieves list of available language options from gitignore api and prints to stdout
func printApiList(url string) {
	data, err := getList(url)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
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

// customIgnoreOptions adds custom ignore options to .gitignore file
func customIgnoreOptions(data []byte) []byte {
	customOptions := []string{".idea", ".vscode"}
	for _, c := range customOptions {
		b := []byte("\n" + c + "/\n")
		data = append(data, b...)
	}
	return data
}

// makeIgnoreFile retrieves gitignore content from api, adds custom entries and creates the actual file
func makeIgnoreFile(targets []string) {
	data, err := getIgnore(targets, ignoreURL)
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
	data, err := config.GetEmbeds(name)
	if err != nil {
		fmt.Println("Unable to load embedded file " + strconv.Quote(name) + ".")
	}
	createDataFile(data, name)
}

// createAllFiles creates all available files in current directory regardless of cli flags
func createAllFiles() {
	for file, _ := range config.FileMap {
		data, err := config.GetEmbeds(file)
		if err != nil {
			fmt.Println("Unable to load embedded file " + strconv.Quote(file) + ".")
		}
		createDataFile(data, file)
	}
}
