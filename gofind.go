package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func main() {

	// TODO: Use CLI package
	var numArgs = len(os.Args)
	if numArgs == 1 {
		fmt.Println("Why no argament?")
		return
	}

	pattern := os.Args[1]
	var rootDirectory string
	switch numArgs {
	case 2:
		// TODO: Error handling here
		rootDirectory, _ = os.Getwd()
	case 3:
		rootDirectory = os.Args[2]
	}

	if !isDirectory(rootDirectory) {
		fmt.Printf("'%v' is not a directory\n", rootDirectory)
		return
	}

	// TODO: Add 'verbose' mode for stuff like this
	fmt.Printf("Looking for '%v' in '%v'\n", pattern, rootDirectory)

	filepath.Walk(rootDirectory, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.Contains(path, pattern) {
			fmt.Println(path)
		}
		return nil
	})
}
