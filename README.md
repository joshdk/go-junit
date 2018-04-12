[![License](https://img.shields.io/github/license/joshdk/go-junit.svg)](https://opensource.org/licenses/MIT)
[![GoDoc](https://godoc.org/github.com/joshdk/go-junit?status.svg)](https://godoc.org/github.com/joshdk/go-junit)
[![Go Report Card](https://goreportcard.com/badge/github.com/joshdk/go-junit)](https://goreportcard.com/report/github.com/joshdk/go-junit)
[![CircleCI](https://circleci.com/gh/joshdk/go-junit.svg?&style=shield)](https://circleci.com/gh/joshdk/go-junit/tree/master)

# Go Junit

🐜 Go library for ingesting JUnit XML reports

## Installing

You can fetch this library by running the following

```bash
go get -u github.com/joshdk/go-junit
```

## Usage

### Data Ingestion

This library has a number of ingestion methods for convenient.

The simplest of which parses raw JUnit XML data.

```go
xml := []byte(`
    <?xml version="1.0" encoding="UTF-8"?>
    <testsuites>
        <testsuite name="JUnitXmlReporter" errors="0" tests="0" failures="0" time="0" timestamp="2013-05-24T10:23:58" />
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
`)

suites, err := junit.Ingest(xml)
if err != nil {
    log.Fatalf("failed to ingest JUnit xml %v", err)
}
```

You can then inspect the contents of the ingestes suites.

```go
for _, suite := range suites {
    fmt.Println(suite.Name)
    for _, test := range suite.Tests {
        fmt.Printf("  %s\n", test.Name)
        if test.Error != nil {
            fmt.Printf("    %s: %s\n", test.Status, test.Error.Error())
        } else {
            fmt.Printf("    %s\n", test.Status)
        }
    }
}
```

And observe some output like this.

```
JUnitXmlReporter
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
if err != nil {
    log.Fatalf("failed to ingest JUnit xml %v", err)
}
```

Or a list of multiple files.

```go
suites, err := junit.IngestFiles([]string{
    "test-reports/report-1.xml",
    "test-reports/report-2.xml",
})
if err != nil {
    log.Fatalf("failed to ingest JUnit xml %v", err)
}
```

Or any `.xml` files inside of a directory.

```go
suites, err := junit.IngestDir("test-reports/")
if err != nil {
    log.Fatalf("failed to ingest JUnit xml %v", err)
}
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

## License

This library is distributed under the [MIT License](https://opensource.org/licenses/MIT), see [LICENSE.txt](https://github.com/joshdk/go-junit/blob/master/LICENSE.txt) for more information.
