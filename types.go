package junit

import (
	"time"
)

type status int

const (
	// StatusPassed represents a JUnit testcase that was run, and did not
	// result in an error or a failure.
	StatusPassed status = iota

	// StatusSkipped represents a JUnit testcase that was intentionally
	// skipped.
	StatusSkipped

	// StatusFailed represents a JUnit testcase that was run, but resulted in
	// a failure. Failures are violations of declared test expectations,
	// such as a failed assertion.
	StatusFailed

	// StatusError represents a JUnit testcase that was run, but resulted in
	// an error. Errors are unexpected violations of the test itself, such as
	// an uncaught exception.
	StatusError
)

type Totals struct {
	Tests    int           `json:"tests" yaml:"tests"`
	Passed   int           `json:"passed" yaml:"passed"`
	Skipped  int           `json:"skipped" yaml:"skipped"`
	Failed   int           `json:"failed" yaml:"failed"`
	Error    int           `json:"error" yaml:"error"`
	Duration time.Duration `json:"duration" yaml:"duration"`
}

type Suite struct {
	Name string `json:"name" yaml:"name"`

	Package string `json:"package" yaml:"package"`

	Properties map[string]string `json:"properties,omitempty" yaml:"properties"`

	Tests []Test `json:"tests,omitempty" yaml:"tests"`

	SystemOut string `json:"stdout,omitempty"`

	SystemErr string `json:"stderr,omitempty"`

	Totals Totals `json:"totals" yaml:"totals"`
}

// Aggregate calculates result sums across all tests.
func (s *Suite) Aggregate() {
	totals := Totals{Tests: len(s.Tests)}

	for _, test := range s.Tests {
		totals.Duration += test.Duration
		switch test.Status {
		case StatusPassed:
			totals.Passed++
		case StatusSkipped:
			totals.Skipped++
		case StatusFailed:
			totals.Failed++
		case StatusError:
			totals.Error++
		}
	}

	s.Totals = totals
}

type Test struct {
	Name string `json:"name" yaml:"name"`

	Classname string `json:"classname" yaml:"classname"`

	Duration time.Duration `json:"duration" yaml:"duration"`

	Status status

	Error error
}

type Error struct {
	Message string `json:"message,omitempty" yaml:"message"`

	Type string `json:"type,omitempty" yaml:"type"`

	Body string `json:"body,omitempty" yaml:"body"`
}

func (err Error) Error() string {
	return err.Body
}
