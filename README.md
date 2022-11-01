[![License][license-badge]][license-link]
[![Actions][github-actions-badge]][github-actions-link]
[![Go Report Card][go-report-badge]][go-report-link]
[![Go Dev][godev-badge]][godev-link]
[![Releases][github-release-badge]][github-release-link]

# Go JUnit

🐜 Go library for ingesting JUnit XML reports

## Installing

You can fetch this library by running the following

```bash
go get -u github.com/joshdk/go-junit
```

## Usage

### Data Ingestion

This library has a number of ingestion methods for convenience.

The simplest of which parses raw JUnit XML data.

```xml
<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
    <testsuite name="JUnitXmlReporter.constructor" errors="0" skipped="1" tests="3" failures="1" time="0.006" timestamp="2013-05-24T10:23:58">
        <properties>
            <property name="java.vendor" value="Sun Microsystems Inc." />
            <property name="compiler.debug" value="on" />
            <property name="project.jdk.classpath" value="jdk.classpath.1.6" />
        </properties>
        <testcase classname="JUnitXmlReporter.constructor" name="should default path to an empty string" time="0.006">
            <failure message="test failure">Assertion failed</failure>
        </testcase>
        <testcase classname="JUnitXmlReporter.constructor" name="should default consolidate to true" time="0">
            <skipped />
        </testcase>
        <testcase classname="JUnitXmlReporter.constructor" name="should default useDotNotation to true" time="0" />
    </testsuite>
</testsuites>
```

```go
xml := []byte(`<?xml …`)

suites, err := junit.Ingest(xml)
```

You can then inspect the contents of the ingested suites.

```go
for _, suite := range suites {
    fmt.Println(suite.Name)
    for _, test := range suite.Tests {
        fmt.Printf("  %s\n", test.Name)
        if test.Error != nil {
            fmt.Printf("    %s: %v\n", test.Status, test.Error)
        } else {
            fmt.Printf("    %s\n", test.Status)
        }
    }
}
```

And observe some output like this.

```
JUnitXmlReporter.constructor
  should default path to an empty string
    failed: Assertion failed
  should default consolidate to true
    skipped
  should default useDotNotation to true
    passed
```

### More Examples

Additionally, you can ingest an entire file.

```go
suites, err := junit.IngestFile("test-reports/report.xml")
```

Or a list of multiple files.

```go
suites, err := junit.IngestFiles([]string{
    "test-reports/report-1.xml",
    "test-reports/report-2.xml",
})
```

Or any `.xml` files inside of a directory.

```go
suites, err := junit.IngestDir("test-reports/")
```

### Data Formats

Due to the lack of implementation consistency in software that generates JUnit XML files, this library needs to take a somewhat looser approach to ingestion. As a consequence, many different possible JUnit formats can easily be ingested.

A single top level `testsuite` tag, containing multiple `testcase` instances.

```xml
<testsuite>
    <testcase name="Test case 1" />
    <testcase name="Test case 2" />
</testsuite>
```

A single top level `testsuites` tag, containing multiple `testsuite` instances.

```xml
<testsuites>
    <testsuite>
        <testcase name="Test case 1" />
        <testcase name="Test case 2" />
    </testsuite>
</testsuites>
```

(Despite not technically being valid XML) Multiple top level `testsuite` tags, containing multiple `testcase` instances.

```xml
<testsuite>
    <testcase name="Test case 1" />
    <testcase name="Test case 2" />
</testsuite>
<testsuite>
    <testcase name="Test case 3" />
    <testcase name="Test case 4" />
</testsuite>
```

In all cases, omitting (or even duplicated) the XML declaration tag is allowed.

```xml
<?xml version="1.0" encoding="UTF-8"?>
```

## Contributing

Found a bug or want to make go-junit better? Please [open a pull request](https://github.com/joshdk/go-junit/compare)!

To make things easier, try out the following:

- Running `go test -v` will run the test suite to verify behavior.

- Running `golangci-lint run` will report any linting issues using [golangci/golangci-lint](https://github.com/golangci/golangci-lint/releases/tag/v1.50.1).

## License

This code is distributed under the [MIT License][license-link], see [LICENSE.txt][license-file] for more information.

<p align="center">
  Created by <a href="https://github.com/joshdk">Josh Komoroske</a> ☕
</p>

[github-actions-badge]:  https://github.com/joshdk/go-junit/workflows/Test/badge.svg?branch=master
[github-actions-link]:   https://github.com/joshdk/go-junit/actions
[github-release-badge]:  https://img.shields.io/github/release/joshdk/go-junit/all.svg
[github-release-link]:   https://github.com/joshdk/go-junit/releases
[go-report-badge]:       https://goreportcard.com/badge/github.com/joshdk/go-junit
[go-report-link]:        https://goreportcard.com/report/github.com/joshdk/go-junit
[godev-badge]:           https://pkg.go.dev/badge/github.com/joshdk/go-junit.svg
[godev-link]:            https://pkg.go.dev/github.com/joshdk/go-junit
[license-badge]:         https://img.shields.io/badge/license-MIT-green.svg
[license-file]:          https://github.com/joshdk/go-junit/blob/master/LICENSE.txt
[license-link]:          https://opensource.org/licenses/MIT
