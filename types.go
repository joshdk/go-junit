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
	Tests     int           `json:"tests" yaml:"tests"`
	Failures  int           `json:"failures" yaml:"failures"`
	Errors    int           `json:"errors" yaml:"errors"`
	Successes int           `json:"successes" yaml:"successes"`
	Duration  time.Duration `json:"duration" yaml:"duration"`
}

type Suite struct {
	Name string `json:"name" yaml:"name"`

	Package string `json:"package" yaml:"package"`

	Properties map[string]string `json:"properties,omitempty" yaml:"properties"`

	Tests []Test `json:"tests,omitempty" yaml:"tests"`

	Totals Totals `json:"totals" yaml:"totals"`

	SystemOut string `json:"stdout,omitempty"`

	SystemErr string `json:"stderr,omitempty"`
}

type Test struct {
	Name string `json:"name" yaml:"name"`

	Classname string `json:"classname" yaml:"classname"`

	Duration time.Duration `json:"duration" yaml:"duration"`

	Status status
	Error  error
}

type Error struct {
	Message string `json:"message,omitempty" yaml:"message"`

	Type string `json:"type,omitempty" yaml:"type"`

	Body string `json:"body,omitempty" yaml:"body"`
}

func (err Error) Error() string {
	return err.Body
}

type Failure struct {
	Message string `json:"message,omitempty" yaml:"message"`

	Type string `json:"type,omitempty" yaml:"type"`

	Body string `json:"body,omitempty" yaml:"body"`
}

func (err Failure) Error() string {
	return err.Body
}
