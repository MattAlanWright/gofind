package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func gofind(pattern string, rootDirectory string, verbose bool) {
	if !isDirectory(rootDirectory) {
		fmt.Printf("'%v' is not a directory\n", rootDirectory)
		return
	}

	if verbose {
		fmt.Printf("Looking for '%v' in '%v'\n", pattern, rootDirectory)
	}

	matches := make([]string, 0, 20)
	filepath.Walk(rootDirectory, func(path string, info fs.FileInfo, err error) error {
		if err == fs.SkipDir {
			if verbose {
				fmt.Printf("Skipping directory %v\n", path)
			}
			return err
		}

		if err == fs.SkipAll {
			if verbose {
				fmt.Printf("Skipping all at path %v\n", path)
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
			matches = append(matches, path)
		}
		return nil
	})

	for _, match := range matches {
		fmt.Println(match)
	}
}

func main() {

	var verbose bool
	app := &cli.App{
		Name:  "gofind",
		Usage: "Find files against a search pattern",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "verbose",
				Usage:       "Print extra information while running",
				Aliases:     []string{"v"},
				Destination: &verbose,
			},
		},
		Action: func(cCtx *cli.Context) error {

			if cCtx.NArg() == 0 {
				fmt.Println("Must provide pattern argument")
				return nil
			}

			pattern := cCtx.Args().First()
			rootDirectory := cCtx.Args().Get(2)
			var err error
			if len(rootDirectory) == 0 {
				rootDirectory, err = os.Getwd()
				if err != nil {
					log.Fatal("Failed to get current directory")
					return err
				}
			}

			gofind(pattern, rootDirectory, verbose)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
