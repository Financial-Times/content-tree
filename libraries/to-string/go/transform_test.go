package tostring

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTransform(t *testing.T) {
	schemas := []Schema{BodyTree, TransitTree}
	for _, s := range schemas {
		for _, test := range getTestCases(t, s) {
			t.Run(test.name, func(t *testing.T) {
				got, err := Transform(test.input, s)

				if err != nil && !test.wantErr {
					t.Errorf("Failed with unexpected error: %v", err)
				}

				if err != nil && test.wantErr {
					return
				}

				if got != strings.TrimSpace(test.output) {
					t.Errorf("got: <%v>\n want: <%v>\n", got, test.output)
				}
			})
		}
	}
}

type TestCase struct {
	name    string
	input   json.RawMessage
	output  string
	wantErr bool
}

func getTestCases(t *testing.T, s Schema) []TestCase {
	t.Helper()

	inputPath := fmt.Sprintf("../../../tests/%s-to-string/input", s)
	outputPath := fmt.Sprintf("../../../tests/%s-to-string/output", s)

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
				name:    fmt.Sprintf("%s-%s", s, caseName),
				input:   input,
				output:  "",
				wantErr: true,
			})
		} else {
			output, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read file %s: %s", outputFile, err)
			}

			testCases = append(testCases, TestCase{
				name:    fmt.Sprintf("%s-%s", s, caseName),
				input:   input,
				output:  string(output),
				wantErr: false,
			})
		}
	}

	return testCases
}
