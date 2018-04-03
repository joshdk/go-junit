package junit

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReparent(t *testing.T) {
	tests := []struct {
		title    string
		input    []byte
		expected string
	}{
		{
			title:    "nil input",
			expected: "<fake-root></fake-root>",
		},
		{
			title:    "empty input",
			input:    []byte(""),
			expected: "<fake-root></fake-root>",
		},
		{
			title:    "xml input",
			input:    []byte(`<testcase name="unit tests" />`),
			expected: `<fake-root><testcase name="unit tests" /></fake-root>`,
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("#%d - %s", index+1, test.title)

		t.Run(name, func(t *testing.T) {
			actual := reparentXML(test.input)

			assert.Equal(t, test.expected, string(actual))
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		title    string
		input    []byte
		expected []xmlNode
	}{
		{
			title: "nil input",
		},
		{
			title: "empty input",
			input: []byte(``),
		},
		{
			title: "plaintext input",
			input: []byte(`This is some data that does not look like xml.`),
		},
		{
			title: "json input",
			input: []byte(`{"This is some data": "that looks like json"}`),
		},
		{
			title: "single xml node",
			input: []byte(`<this-is-a-tag/>`),
			expected: []xmlNode{
				{
					XMLName: xml.Name{
						Local: "this-is-a-tag",
					},
				},
			},
		},
		{
			title: "multiple xml nodes",
			input: []byte(`
				<this-is-a-tag/>
				<this-is-also-a-tag/>
			`),
			expected: []xmlNode{
				{
					XMLName: xml.Name{
						Local: "this-is-a-tag",
					},
				},
				{
					XMLName: xml.Name{
						Local: "this-is-also-a-tag",
					},
				},
			},
		},
		{
			title: "single xml node with content",
			input: []byte(`<this-is-a-tag>This is some content.</this-is-a-tag>`),
			expected: []xmlNode{
				{
					XMLName: xml.Name{
						Local: "this-is-a-tag",
					},
					Content: []byte("This is some content."),
				},
			},
		},
		{
			title: "single xml node with encoded content",
			input: []byte(`<this-is-a-tag>&lt;sender&gt;John Smith&lt;/sender&gt;</this-is-a-tag>`),
			expected: []xmlNode{
				{
					XMLName: xml.Name{
						Local: "this-is-a-tag",
					},
					Content: []byte("<sender>John Smith</sender>"),
				},
			},
		},
		{
			title: "single xml node with cdata content",
			input: []byte(`<this-is-a-tag><![CDATA[<sender>John Smith</sender>]]></this-is-a-tag>`),
			expected: []xmlNode{
				{
					XMLName: xml.Name{
						Local: "this-is-a-tag",
					},
					Content: []byte("<sender>John Smith</sender>"),
				},
			},
		},
		{
			title: "single xml node with attributes",
			input: []byte(`<this-is-a-tag name="my name" status="passed"></this-is-a-tag>`),
			expected: []xmlNode{
				{
					XMLName: xml.Name{
						Local: "this-is-a-tag",
					},
					Attrs: map[string]string{
						"name":   "my name",
						"status": "passed",
					},
				},
			},
		},
		{
			title: "single xml node with encoded attributes",
			input: []byte(`<this-is-a-tag name="&lt;sender&gt;John Smith&lt;/sender&gt;"></this-is-a-tag>`),
			expected: []xmlNode{
				{
					XMLName: xml.Name{
						Local: "this-is-a-tag",
					},
					Attrs: map[string]string{
						"name": "<sender>John Smith</sender>",
					},
				},
			},
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("#%d - %s", index+1, test.title)

		t.Run(name, func(t *testing.T) {
			actual, err := parse(test.input)

			require.Nil(t, err)

			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestExtract(t *testing.T) {
	tests := []struct {
		title    string
		input    []byte
		expected []byte
		err      string
	}{
		{
			title: "nil content",
		},
		{
			title: "empty content",
			input: []byte(""),
		},
		{
			title:    "simple content",
			input:    []byte("hello world"),
			expected: []byte("hello world"),
		},
		{
			title:    "complex content",
			input:    []byte("No bugs 🐜"),
			expected: []byte("No bugs 🐜"),
		},
		{
			title:    "simple encoded content",
			input:    []byte("I &lt;/3 XML"),
			expected: []byte("I </3 XML"),
		},
		{
			title:    "complex encoded content",
			input:    []byte(`&lt;[[&apos;/\&quot;]]&gt;`),
			expected: []byte(`<[['/\"]]>`),
		},
		{
			title: "empty cdata content",
			input: []byte("<![CDATA[]]>"),
		},
		{
			title:    "simple cdata content",
			input:    []byte("<![CDATA[hello world]]>"),
			expected: []byte("hello world"),
		},
		{
			title:    "complex cdata content",
			input:    []byte("<![CDATA[I </3 XML]]>"),
			expected: []byte("I </3 XML"),
		},
		{
			title:    "complex encoded cdata content",
			input:    []byte("<![CDATA[I &lt;/3 XML]]>"),
			expected: []byte("I &lt;/3 XML"),
		},
		{
			title:    "encoded content then cdata content",
			input:    []byte("I want to say that <![CDATA[I </3 XML]]>"),
			expected: []byte("I want to say that I </3 XML"),
		},
		{
			title:    "cdata content then encoded content",
			input:    []byte("<![CDATA[I </3 XML]]> a lot"),
			expected: []byte("I </3 XML a lot"),
		},
		{
			title:    "mixture of encoded and cdata content",
			input:    []byte("I &lt;/3 XML <![CDATA[a lot]]>. 🐜 You probably <![CDATA[</3 XML]]> too."),
			expected: []byte("I </3 XML a lot. 🐜 You probably </3 XML too."),
		},
		{
			title: "unmatched cdata start tag",
			input: []byte("<![CDATA["),
			err:   "unmatched CDATA start tag",
		},
		{
			title: "unmatched cdata end tag",
			input: []byte("]]>"),
			err:   "unmatched CDATA end tag",
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("#%d - %s", index+1, test.title)

		t.Run(name, func(t *testing.T) {
			actual, err := extractContent(test.input)

			checkError(t, test.err, err)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func checkError(t *testing.T, expected string, actual error) {
	t.Helper()

	switch {
	case expected == "" && actual == nil:
		return
	case expected == "" && actual != nil:
		t.Fatalf("expected no error but got %q", actual.Error())
	case expected != "" && actual == nil:
		t.Fatalf("expected %q but got nil", expected)
	case expected == actual.Error():
		return
	default:
		t.Fatalf("expected %q but got %q", expected, actual.Error())
	}
}
