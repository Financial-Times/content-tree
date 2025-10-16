package tocontenttree

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nsf/jsondiff"
)

func TestTransform(t *testing.T) {
	for _, test := range getTestCases(t) {
		t.Run(test.name, func(t *testing.T) {
			bodyTree, err := Transform(test.input)
			if err != nil && !test.wantErr {
				t.Errorf("Failed with unexpected error: %v", err)
			}
			if err != nil && test.wantErr {
				return
			}

			want := strings.TrimSpace(test.output)
			got, err := json.Marshal(bodyTree)
			opts := jsondiff.DefaultJSONOptions()
			diffQuery, _ := jsondiff.Compare(got, []byte(want), &opts)
			if diffQuery != jsondiff.FullMatch {
				t.Errorf("got: %s\n\n want: %s\n", got, want)
			}
		})
	}
}

type TestCase struct {
	name    string
	input   string
	output  string
	wantErr bool
}

func getTestCases(t *testing.T) []TestCase {
	t.Helper()

	inputPath := "../../../tests/bodyxml-to-content-tree/input"
	outputPath := "../../../tests/bodyxml-to-content-tree/output"

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
		outputFile := filepath.Join(outputPath, caseName+".json")

		if _, err := os.Stat(outputFile); errors.Is(err, os.ErrNotExist) {
			testCases = append(testCases, TestCase{
				name:    caseName,
				input:   string(input),
				output:  "",
				wantErr: true,
			})
		} else {
			output, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read file %s: %s", outputFile, err)
			}

			testCases = append(testCases, TestCase{
				name:    caseName,
				input:   string(input),
				output:  string(output),
				wantErr: false,
			})
		}
	}

	return testCases
}
