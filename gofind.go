package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const verbose = false

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
		fmt.Println("Must provide pattern argument")
		return
	}

	pattern := os.Args[1]
	var rootDirectory string

	var err error
	switch numArgs {
	case 2:
		// TODO: Error handling here
		rootDirectory, err = os.Getwd()
		if err != nil {
			log.Fatal("Failed to get current directory")
			return
		}

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
		if err == fs.SkipDir {
			if verbose {
				fmt.Printf("Skipping directory %v\n", path)
			}
			return err
		}

		if err == fs.SkipAll {
			if verbose {
				fmt.Printf("Skipping ALL at path %v\n", path)
			}
			return err
		}

		if err != nil {
			if verbose {
				fmt.Printf("Error in WalkFunc - %v\n", err)
			}
			return nil
		}

		if !info.IsDir() && strings.Contains(path, pattern) {
			fmt.Println(path)
		}
		return nil
	})
}
