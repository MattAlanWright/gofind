package main

import (
	"errors"
	"fmt"
	"io/fs"
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

func gofind(pattern string, rootDirectory string) {
	matches := make([]string, 0, 20)
	filepath.Walk(rootDirectory, func(path string, info fs.FileInfo, err error) error {
		if err == fs.SkipDir || err == fs.SkipAll {
			return err
		}
		if err != nil {
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
		return "", err
	}
	if !isDirectory(rootDirectory) {
		errString := fmt.Sprintf("No such directory %v\n", rootDirectory)
		return "", errors.New(errString)
	}
	return rootDirectory, nil
}

func main() {
	app := &cli.App{
		Name:      "gofind",
		Usage:     "Search for files containing a pattern",
		UsageText: "gofind <pattern> [directory]",
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

			gofind(pattern, rootDirectory)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
