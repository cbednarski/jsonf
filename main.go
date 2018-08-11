package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	formatCompact = "c"
	formatString = "s"
	formatTab = "t"
)

const helpText = `jsonf - A simple JSON formatter

When invoked on a single file, the file will be reformatted and the output sent
to stdout. You may rewrite the original file(s) using the -w flag.

When invoked on a directory, the program will only operate on files ending in
.json. To operate on multiple files without the json extension, specify each one
as an additional argument.

When processing multiple files or one or more directories, the program will
output the filename of each file before it is processed, and will attempt to
continue when it encouters errors.

Usage:

  jsonf [options] filename or directory ...

Example:

  jsonf -i sss myfile.json    Indent using 3 spaces and write to stdout
  jsonf -w -i t myfile.json   Indent using tabs and rewrite the original files
  jsonf -w .                  Rewrite all .json files in the current directory
  jsonf -r file1 file2 dir    Rewrite file1, file2, and all .json files under dir

Options:

  -w  Overwrite files in place
  -r  Recurse into subdirectories
  -i  Set the indentation string (defaults to 2 spaces). You can use 's' and 't'
      to replace literal space and tab characters
  -c  Compact (minify) rather than indent. -c wins over -i if both are specified
  -h  Show help

Copyright 2018 Chris Bednarski <chris@cbednarski.com> - MIT License
Portions copyright 2009 the Go Authors - BSD License https://golang.org/LICENSE

Report issues to https://github.com/cbednarski/jsonf
`

var (
	exitCode = 0
	stdout = os.Stdout
	stderr = os.Stderr
)

func printError(err error) {
	stderr.WriteString(fmt.Sprintf("%s\n", err))
	exitCode = 1
}

func format(filename string, indent string) (*bytes.Buffer, error) {
	output := &bytes.Buffer{}

	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	input = bytes.TrimSpace(input)

	if indent == formatCompact {
		if err := json.Compact(output, input); err != nil {
			return nil, err
		}
	} else {
		if err := json.Indent(output, input, "", indent); err != nil {
			return nil, err
		}
		output.WriteString("\n")
	}

	return output, nil
}

func formatFile(filename string, indent string, replace bool) error {
	data, err := format(filename, indentString(indent))
	if err != nil {
		return err
	}

	if replace {
		info, err := os.Stat(filename)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(filename, data.Bytes(), info.Mode()); err != nil {
			return err
		}
	} else {
		data.WriteTo(stdout)
	}

	return nil
}

func indentString(input string) string {
	// Special case for compaction
	if strings.Contains(input, formatCompact) {
		return formatCompact
	}

	input = strings.Replace(input, formatString, " ", -1)
	input = strings.Replace(input, formatTab, "\t", -1)

	return input
}

func listJSONFiles(path string, recurse bool) ([]string, error) {
	files := []string{}
	if isDir(path) {
		filepath.Walk(path, func(innerPath string, info os.FileInfo, err error) error{
			// We always want to talk the current directory, but not any
			// subdirectories, unless recurse is true.
			if !recurse && info.IsDir() && innerPath != path {
				return filepath.SkipDir
			}

			if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
				files = append(files, innerPath)
			}

			return nil
		})
	} else {
		files = []string{path}
	}

	return files, nil
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func header(filename string) {
	fmt.Printf("--- %s\n", filename)
}

func formatPath(path, indent string, replace, recurse, headers bool) error {
	if isDir(path) {
		files, err := listJSONFiles(flag.Args()[0], recurse)
		if err != nil {
			return err
		}
		errors := 0
		for _, f := range files {
			header(f)
			if err := formatFile(f, indent, replace); err != nil {
				printError(err)
			}
		}
		if errors > 0 {
			return fmt.Errorf("Encountered %d errors in %s", errors, path)
		}
	} else {
		if headers {
			header(path)
		}
		if err := formatFile(path, indent, replace); err != nil {
			return err
		}
	}
	return nil
}

func wrappedMain() error {
	replace := flag.Bool("w", false, "overwrite files in place")
	recurse := flag.Bool("r", false, "recurse into subdirectories")
	compact := flag.Bool("c", false, "compact instead of intent")
	indent := flag.String("i", "  ", `indent string, e.g. "  "`)
	help := flag.Bool("h", false, "show help")

	flag.Parse()

	if *help {
		fmt.Println(helpText)
		os.Exit(0)
	}

	if *compact {
		*indent = formatCompact
	}

	if len(flag.Args()) < 1 {
		return errors.New(helpText)
	}

	headers := false
	if len(flag.Args()) > 1 {
		headers = true
	}

	for _, path := range flag.Args() {
		if err := formatPath(path, *indent, *replace, *recurse, headers); err != nil {
			printError(err)
		}
	}

	return nil
}

func main() {
	if err := wrappedMain(); err != nil {
		stderr.WriteString(fmt.Sprintf("%s\n", err))
		os.Exit(1)
	}
	os.Exit(exitCode)
}
