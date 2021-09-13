package main

import (
	//"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	// ignoreURL is the base url for the gitignore API
	ignoreURL      string = "https://www.toptal.com/developers/gitignore/api"
	versionMessage string = "goignore version: 0.2.4"
	listMessage    string = "To get a list of valid targets, run goignore --list"

	helpMessage string = `
Usage: goignore [OPTIONS] [ARGS]...
Handy CLI to generate great gitignore files.
Options:
	--version: Display goignore's version.
	--help: Show this help message and exit.
	--list: Show the valid gitignore.io targets.
Arguments:
	TARGETS: Space separated list of gitignore.io targets.	[required]
Examples:
$ goignore macos vscode go
$ goignore --list
`
)

var (
	helpFlag    bool
	versionFlag bool
	listFlag    bool
)

func main() {
	flag.BoolVar(&helpFlag, "help", false, "--help")
	flag.BoolVar(&versionFlag, "version", false, "--version")
	flag.BoolVar(&listFlag, "list", false, "--list")

	flag.Usage = func() {
		fmt.Println(helpMessage)
		os.Exit(1)
	}

	flag.Parse()

	run()

}

func run() {
	if flag.NArg() < 2 && !(helpFlag || listFlag || versionFlag) {
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
		printList(os.Stdout, ignoreURL)
		os.Exit(0)

	case os.Args[1] == "list":
		fmt.Println(listMessage)
		os.Exit(1)

	default:
		ignoreList := os.Args[1:]
		makeIgnoreFile(ignoreList, ignoreURL)
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

// WriteToIgnoreFile creates .gitignore data using supplied data
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
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error: %s.\n", err)
		os.Exit(1)
	}

	fmt.Println("Creating .gitignore file in " + strconv.Quote(cwd) + " directory...")
	data, err := GetIgnore(targets, url)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	data = customIgnoreOptions(data)

	err = WriteToIgnoreFile(data, ".gitignore")
	if err != nil {
		fmt.Printf("Error: %s.\n", err)
		os.Exit(1)
	} else {
		fmt.Println("Done!")
	}
}
