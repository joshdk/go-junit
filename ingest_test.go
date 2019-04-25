// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.

package junit

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestInjest(t *testing.T) {
	tests := []struct {
		title    string
		input    []byte
		expected string
	}{
		{
			title: "xml input",
			input: []byte(`<testsuite errors="0" failures="1" file="Foo.java"><testcase name="unit tests" file="Foo.java"/></testsuite>`),
			expected: `
[
  {
    "name": "",
    "package": "",
    "tests": [
      {
        "name": "unit tests",
        "classname": "",
        "duration": 0,
        "status": "passed",
        "error": null,
        "Properties": {
          "file": "Foo.java",
          "name": "unit tests"
        }
      }
    ],
    "totals": {
      "tests": 1,
      "passed": 1,
      "skipped": 0,
      "failed": 0,
      "error": 0,
      "duration": 0
    }
  }
]`,
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("#%d - %s", index+1, test.title)

		t.Run(name, func(t *testing.T) {
			actual, err := Ingest(test.input)
			if err != nil {
				t.Fatal(err)
			}
			if actual == nil {
				t.Fatalf("No suites found!")
			}

			actualJSON, err := json.MarshalIndent(actual, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, strings.TrimSpace(test.expected), string(actualJSON))
		})
	}
}
