package config

import (
	"embed"
	"fmt"
	"strconv"
)

//go:embed files/*
var files embed.FS

// FileMap holds the 'desired filename' as key and the embedded file as value
// embed package does not allow '.' files to be compiled into binary executable
var FileMap = map[string]string{
	".flake8":        "flake8",
	"License":        "License",
	"pyproject.toml": "pyproject.toml",
	".dockerignore":  "dockerignore",
	"Dockerfile":     "Dockerfile",
}

// GetEmbeds reads embedded files in binary and returns their content as string
func GetEmbeds(filename string) (string, error) {
	var data string
	name, ok := FileMap[filename]
	if ok {
		d, err := files.ReadFile("files/" + name)
		if err != nil {
			return string(d), err
		}
		data = string(d)
	} else {
		fmt.Println("Unable to load embedded file " + strconv.Quote(filename) + ".")
	}
	return data, nil
}
