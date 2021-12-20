package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestIsDir(t *testing.T) {
	if !isDir("test-fixtures") {
		t.Errorf("test-fixtures should be a directory")
	}

	if isDir("main.go") {
		t.Errorf("main.go should not be a directory")
	}
}

var expectedFormatSpaces = `{
  "icecream": [
    "chocolate",
    "strawberry",
    "vanilla"
  ]
}
`

var expectedFormatTabs = `{
	"icecream": [
		"chocolate",
		"strawberry",
		"vanilla"
	]
}
`

var expectedFormatCompact = `{"icecream":["chocolate","strawberry","vanilla"]}`

type FormatTest struct {
	InputFile string
	Indent    string
	Expected  string
}

func TestFormat(t *testing.T) {
	inputs := map[string]FormatTest{
		"spaces": {
			InputFile: filepath.Join("test-fixtures", "icecream.json"),
			Indent:    "  ",
			Expected:  expectedFormatSpaces,
		},
		"tabs": {
			InputFile: filepath.Join("test-fixtures", "icecream.json"),
			Indent:    "\t",
			Expected:  expectedFormatTabs,
		},
		"compact": {
			InputFile: filepath.Join("test-fixtures", "icecream.json"),
			Indent:    formatCompact,
			Expected:  expectedFormatCompact,
		},
	}

	for testCase, input := range inputs {
		actual, err := format(input.InputFile, input.Indent)
		if err != nil {
			t.Errorf("Test case %q error: %s", testCase, err)
			continue
		}

		actualString := actual.String()
		if actualString != input.Expected {
			t.Errorf("Test case %q doesn't match\n--Expected--\n%s\n--Actual--\n%s", testCase, input.Expected, actualString)
		}
	}
}

func TestIndentString(t *testing.T) {
	actual := indentString("sts ")
	expected := " \t  "
	if actual != expected {
		t.Errorf("Expected %q, found %q", expected, actual)
	}

	actual2 := indentString("cststs")
	expected2 := "c"
	if actual2 != expected2 {
		t.Errorf("Expected %q, found %q", expected, actual)
	}
}

func TestListJSONFiles(t *testing.T) {
	expectedFiles1 := []string{filepath.Join("test-fixtures", "icecream.json")}

	actualFiles1, err := listJSONFiles("test-fixtures", false)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expectedFiles1, actualFiles1) {
		t.Errorf("No match. Expected %v, found %v", expectedFiles1, actualFiles1)
	}
}

func TestListJSONFilesRecursive(t *testing.T) {
	expectedFiles2 := []string{
		filepath.Join("test-fixtures", "dogs", "corgi.json"),
		filepath.Join("test-fixtures", "dogs", "dachshund.json"),
		filepath.Join("test-fixtures", "icecream.json"),
	}

	actualFiles2, err := listJSONFiles("test-fixtures", true)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expectedFiles2, actualFiles2) {
		t.Errorf("No match. Expected %v, found %v", expectedFiles2, actualFiles2)
	}
}

func TestFormatFile(t *testing.T) {
	data, err := ioutil.ReadFile(filepath.Join("test-fixtures", "icecream.json"))
	if err != nil {
		t.Fatal(err)
	}

	file, err := ioutil.TempFile("", "jsonf-test")
	if err != nil {
		t.Fatal(err)
	}

	filename := file.Name()

	file.Write(data)
	file.Close()

	defer os.Remove(filename)

	if err := formatFile(filename, "ss", true); err != nil {
		t.Fatal(err)
	}

	actual, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	actualString := string(actual)

	if actualString != expectedFormatSpaces {
		t.Errorf("Output does not match:\n--Expected--\n%s\n--Actual--\n%s", expectedFormatSpaces, actualString)
	}
}
