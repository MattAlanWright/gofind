package main

import (
	"errors"
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

func getRootDirFromArg(directoryArg string) (string, error) {
	rootDirectory, err := filepath.Abs(directoryArg)
	if err != nil {
		log.Fatalf("Failed to get absolute path of %v\n", directoryArg)
		return "", err
	}
	if !isDirectory(rootDirectory) {
		log.Fatalf("No such directory %v\n", rootDirectory)
		return "", errors.New("Invalid directory")
	}
	return rootDirectory, nil
}

func main() {

	var verbose bool
	app := &cli.App{
		Name:      "gofind",
		Usage:     "Search for files containing a pattern",
		UsageText: "gofind <pattern> [directory]",
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
				cli.ShowAppHelp(cCtx)
				os.Exit(1)
			}

			var err error
			pattern := cCtx.Args().First()
			rootDirectory, err := getRootDirFromArg(cCtx.Args().Get(1))
			if err != nil {
				return err
			}

			gofind(pattern, rootDirectory, verbose)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
