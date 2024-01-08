package stringifier

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStringify(t *testing.T) {
	for _, test := range getTestCases(t) {
		t.Run(test.name, func(t *testing.T) {
			got, err := Stringify(test.input)

			if err != nil && !test.err {
				t.Errorf("Failed with unexpected error: %v", err)
			}
			if err != nil && test.err {
				return
			}

			if got != strings.TrimSpace(test.output) {
				t.Errorf("got: <%v>\n wanted: <%v>\n", got, test.output)
			}
		})
	}
}

type TestCase struct {
	name   string
	input  json.RawMessage
	output string
	err    bool
}

func getTestCases(t *testing.T) []TestCase {
	t.Helper()
	inputPath := "../../tests/content-tree-to-string/input"
	outputPath := "../../tests/content-tree-to-string/output"

	entries, err := os.ReadDir(inputPath)
	if err != nil {
		t.Fatal(err)
	}

	testCases := make([]TestCase, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		inputFile := filepath.Join(inputPath, entry.Name())

		input, err := os.ReadFile(inputFile)
		if err != nil {
			t.Fatalf("Failed to read file %s: %s", inputFile, err)
		}

		caseName := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		outputFile := filepath.Join(outputPath, caseName+".text")

		if _, err := os.Stat(outputFile); errors.Is(err, os.ErrNotExist) {
			testCases = append(testCases, TestCase{
				name:   caseName,
				input:  input,
				output: "",
				err:    true,
			})
		} else {
			output, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read file %s: %s", outputFile, err)
			}

			testCases = append(testCases, TestCase{
				name:   caseName,
				input:  input,
				output: string(output),
				err:    false,
			})
		}
	}

	return testCases
}
