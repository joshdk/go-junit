// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.

package junit

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInjest(t *testing.T) {
	tests := []struct {
		title    string
		input    []byte
		expected []Suite
	}{
		{
			title: "xml input",
			input: []byte(`<testsuite errors="0" failures="1" file="Foo.java"><testcase name="unit tests" file="Foo.java"/></testsuite>`),
			expected: []Suite{
				{
					Tests: []Test{
						{
							Name:   "unit tests",
							Status: "passed",
							Properties: map[string]string{
								"file": "Foo.java",
								"name": "unit tests",
							},
						},
					},
					Totals: Totals{
						Tests:  1,
						Passed: 1,
					},
				},
			},
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("#%d - %s", index+1, test.title)

		t.Run(name, func(t *testing.T) {
			actual, err := Ingest(test.input)
			require.Nil(t, err)
			require.NotEmpty(t, actual)
			require.Equal(t, test.expected, actual)
		})
	}
}
