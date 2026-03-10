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
			if test.wantErr && err == nil {
				t.Fatalf("Expected an error but got nil")
			}
			if !test.wantErr && err != nil {
				t.Fatalf("Failed with unexpected error: %v", err)
			}
			if bodyTree == nil {
				t.Fatalf("Expected a body tree but got nil")
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

func TestTransformSkippedNodeErrors(t *testing.T) {
	testCases := []struct {
		name              string
		input             string
		wantOutput        string
		wantErrorContains []string
	}{
		{
			name:       "unsupported tag",
			input:      `<body><p>Before</p><mystery><p>Skipped</p></mystery><p>After</p></body>`,
			wantOutput: `{"type":"root","body":{"type":"body","children":[{"type":"paragraph","children":[{"type":"text","value":"Before"}]},{"type":"paragraph","children":[{"type":"text","value":"After"}]}],"version":1}}`,
			wantErrorContains: []string{
				"skipped unsupported element <mystery>",
			},
		},
		{
			name:       "unknown div transformer result",
			input:      `<body><p>Before</p><div class="mystery"><p>Skipped</p></div><p>After</p></body>`,
			wantOutput: `{"type":"root","body":{"type":"body","children":[{"type":"paragraph","children":[{"type":"text","value":"Before"}]},{"type":"paragraph","children":[{"type":"text","value":"After"}]}],"version":1}}`,
			wantErrorContains: []string{
				"skipped unsupported element <div>",
			},
		},
		{
			name:       "invalid child node",
			input:      `<body><ul><p>Skipped</p><li><p>Kept</p></li></ul></body>`,
			wantOutput: `{"type":"root","body":{"type":"body","children":[{"type":"list","children":[{"type":"list-item","children":[{"type":"paragraph","children":[{"type":"text","value":"Kept"}]}]}],"ordered":false}],"version":1}}`,
			wantErrorContains: []string{
				`skipped invalid child node "paragraph" under "list"`,
			},
		},
		{
			name:       "invalid text node",
			input:      `<body><ul>Skipped<li><p>Kept</p></li></ul></body>`,
			wantOutput: `{"type":"root","body":{"type":"body","children":[{"type":"list","children":[{"type":"list-item","children":[{"type":"paragraph","children":[{"type":"text","value":"Kept"}]}]}],"ordered":false}],"version":1}}`,
			wantErrorContains: []string{
				`skipped invalid text node under "list"`,
			},
		},
		{
			name:       "multiple skipped errors are collated",
			input:      `<body><mystery/><div class="mystery"/><ul>Skipped<p>Also skipped</p><li><p>Kept</p></li></ul></body>`,
			wantOutput: `{"type":"root","body":{"type":"body","children":[{"type":"list","children":[{"type":"list-item","children":[{"type":"paragraph","children":[{"type":"text","value":"Kept"}]}]}],"ordered":false}],"version":1}}`,
			wantErrorContains: []string{
				"skipped unsupported element <mystery>",
				"skipped unsupported element <div>",
				`skipped invalid text node under "list"`,
				`skipped invalid child node "paragraph" under "list"`,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			bodyTree, err := Transform(test.input)
			if err == nil {
				t.Fatal("Expected an error but got nil")
			}
			if bodyTree == nil {
				t.Fatal("Expected a body tree but got nil")
			}

			for _, wantErr := range test.wantErrorContains {
				if !strings.Contains(err.Error(), wantErr) {
					t.Fatalf("expected error to contain %q, got %v", wantErr, err)
				}
			}

			assertJSONMatch(t, bodyTree, test.wantOutput)
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
				wantErr: skippedErrorFixtures[caseName],
			})
		}
	}

	return testCases
}

func assertJSONMatch(t *testing.T, v any, want string) {
	t.Helper()

	got, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Failed to marshal value: %v", err)
	}

	opts := jsondiff.DefaultJSONOptions()
	diffQuery, _ := jsondiff.Compare(got, []byte(strings.TrimSpace(want)), &opts)
	if diffQuery != jsondiff.FullMatch {
		t.Fatalf("got: %s\n\n want: %s\n", got, strings.TrimSpace(want))
	}
}

var skippedErrorFixtures = map[string]bool{
	"simple-body-invalid-child":     true,
	"simple-body-section-innumbers": true,
	"simple-body-unknown-div":       true,
	"simple-body-unknown-tag":       true,
}
