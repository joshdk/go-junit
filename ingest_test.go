// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package junit

import (
	"fmt"
	"testing"
	"time"
)

func TestExamplesInTheWild(t *testing.T) {
	tests := []struct {
		title    string
		filename string
		origin   string
		check    func(*testing.T, []Suite)
	}{
		{
			title:    "catchsoftware example",
			filename: "testdata/catchsoftware.xml",
			origin:   "https://help.catchsoftware.com/display/ET/JUnit+Format",
			check: func(t *testing.T, suites []Suite) {
				assertLen(t, suites, 2)
				assertLen(t, suites[0].Tests, 0)
				assertLen(t, suites[1].Tests, 3)
				assertError(t, suites[1].Tests[0].Error, "Assertion failed")
			},
		},
		{
			title:    "cubic example",
			filename: "testdata/cubic.xml",
			origin:   "https://llg.cubic.org/docs/junit/",
			check: func(t *testing.T, suites []Suite) {
				assertLen(t, suites, 1)
				assertLen(t, suites[0].Tests, 1)
				assertEqual(t, "STDOUT text", suites[0].SystemOut)
				assertEqual(t, "STDERR text", suites[0].SystemErr)
				assertEqual(t, "STDOUT text", suites[0].Tests[0].SystemOut)
				assertEqual(t, "STDERR text", suites[0].Tests[0].SystemErr)
			},
		},
		{
			title:    "go-junit-report example",
			filename: "testdata/go-junit-report.xml",
			origin:   "https://github.com/jstemmer/go-junit-report/blob/master/testdata/06-report.xml",
			check: func(t *testing.T, suites []Suite) {
				assertLen(t, suites, 2)
				assertLen(t, suites[0].Tests, 2)
				assertLen(t, suites[1].Tests, 2)
				assertEqual(t, "1.0", suites[0].Properties["go.version"])
				assertEqual(t, "1.0", suites[1].Properties["go.version"])
				assertError(t, suites[1].Tests[0].Error, "file_test.go:11: Error message\nfile_test.go:11: Longer\n\terror\n\tmessage.")
			},
		},
		{
			title:    "go-junit-report skipped example",
			filename: "testdata/go-junit-report-skipped.xml",
			origin:   "https://github.com/jstemmer/go-junit-report/blob/master/testdata/03-report.xml",
			check: func(t *testing.T, suites []Suite) {
				assertLen(t, suites, 1)
				assertLen(t, suites[0].Tests, 2)
				assertEqual(t, "package/name", suites[0].Name)
				assertEqual(t, "TestOne", suites[0].Tests[0].Name)
				assertEqual(t, "file_test.go:11: Skip message", suites[0].Tests[0].Message)
			},
		},
		{
			title:    "ibm example",
			filename: "testdata/ibm.xml",
			origin:   "https://www.ibm.com/support/knowledgecenter/en/SSQ2R2_14.2.0/com.ibm.rsar.analysis.codereview.cobol.doc/topics/cac_useresults_junit.html",
			check: func(t *testing.T, suites []Suite) {
				assertLen(t, suites, 1)
				assertLen(t, suites[0].Tests, 1)
				assertError(t, suites[0].Tests[0].Error, "\nWARNING: Use a program name that matches the source file name\nCategory: COBOL Code Review – Naming Conventions\nFile: /project/PROGRAM.cbl\nLine: 2\n      ")
			},
		},
		{
			title:    "jenkinsci example",
			filename: "testdata/jenkinsci.xml",
			origin:   "https://github.com/jenkinsci/junit-plugin/blob/master/src/test/resources/hudson/tasks/junit/junit-report-1463.xml",
			check: func(t *testing.T, suites []Suite) {
				assertLen(t, suites, 1)
				assertLen(t, suites[0].Tests, 6)
				assertEqual(t, "\n", suites[0].Properties["line.separator"])
				assertEqual(t, `\`, suites[0].Properties["file.separator"])
			},
		},
		{
			title:    "nose2 example",
			filename: "testdata/nose2.xml",
			origin:   "https://nose2.readthedocs.io/en/latest/plugins/junitxml.html",
			check: func(t *testing.T, suites []Suite) {
				assertLen(t, suites, 1)
				assertLen(t, suites[0].Tests, 25)
				assertError(t, suites[0].Tests[22].Error, "Traceback (most recent call last):\n  File \"nose2/tests/functional/support/scenario/tests_in_package/pkg1/test/test_things.py\", line 13, in test_typeerr\n    raise TypeError(\"oops\")\nTypeError: oops\n")
			},
		},
		{
			title:    "python junit-xml example",
			filename: "testdata/python-junit-xml.xml",
			origin:   "https://pypi.org/project/junit-xml/",
			check: func(t *testing.T, suites []Suite) {
				assertLen(t, suites, 1)
				assertLen(t, suites[0].Tests, 1)
				assertEqual(t, "\n                I am stdout!\n            ", suites[0].Tests[0].SystemOut)
				assertEqual(t, "\n                I am stderr!\n            ", suites[0].Tests[0].SystemErr)
				assertEqual(t, "some.property.value", suites[0].Tests[0].Properties["some.property.name"])
				assertEqual(t, "some.class.name", suites[0].Tests[0].Properties["classname"])
			},
		},
		{
			title:    "surefire example",
			filename: "testdata/surefire.xml",
			origin:   "https://gist.github.com/rwbergstrom/6f0193b1a12dca9d358e6043ee6abba4",
			check: func(t *testing.T, suites []Suite) {
				assertLen(t, suites, 1)
				assertLen(t, suites[0].Tests, 1)
				assertEqual(t, "\n", suites[0].Properties["line.separator"])
				assertEqual(t, "Hello, World\n", suites[0].Tests[0].SystemOut)
				assertEqual(t, "I'm an error!\n", suites[0].Tests[0].SystemErr)

				testcase := Test{
					Name:      "testStdoutStderr",
					Classname: "com.example.FooTest",
					Duration:  1234560 * time.Millisecond,
					Status:    StatusFailed,
					Error: Error{
						Type: "java.lang.AssertionError",
						Body: "java.lang.AssertionError\n\tat com.example.FooTest.testStdoutStderr(FooTest.java:13)\n",
					},
					Properties: map[string]string{
						"classname": "com.example.FooTest",
						"name":      "testStdoutStderr",
						"time":      "1,234.56",
					},
					SystemOut: "Hello, World\n",
					SystemErr: "I'm an error!\n",
				}

				assertEqual(t, testcase, suites[0].Tests[0])
			},
		},
		{
			title:    "fastlane example",
			filename: "testdata/fastlane-trainer.xml",
			check: func(t *testing.T, suites []Suite) {
				assertLen(t, suites, 1)
				assertLen(t, suites[0].Tests, 4)

				testcase := Test{
					Name:      "testSomething()",
					Classname: "TestClassSample",
					Duration:  342 * time.Millisecond,
					Status:    StatusFailed,
					Message:   "XCTAssertTrue failed",
					Error: Error{
						Message: "XCTAssertTrue failed",
						Body:    "\n            ",
					},
					Properties: map[string]string{
						"classname": "TestClassSample",
						"name":      "testSomething()",
						"time":      "0.342",
					},
				}

				assertEqual(t, testcase, suites[0].Tests[2])
				assertError(t, suites[0].Tests[2].Error, "XCTAssertTrue failed")
				assertError(t, suites[0].Tests[3].Error, "NullPointerException")
			},
		},
		{
			title:    "phpunit example",
			filename: "testdata/phpunit.xml",
			check: func(t *testing.T, suites []Suite) {
				assertLen(t, suites, 1)
				assertLen(t, suites[0].Tests, 0)
				assertLen(t, suites[0].Suites, 1)

				suite := suites[0].Suites[0]
				assertLen(t, suite.Tests, 1)
				assertLen(t, suite.Suites, 2)

				assertEqual(t, "SampleTest", suite.Name)
				assertEqual(t, "/untitled/tests/SampleTest.php", suite.Properties["file"])

				testcase := Test{
					Name:      "testA",
					Classname: "SampleTest",
					Duration:  5917 * time.Microsecond,
					Status:    StatusPassed,
					Properties: map[string]string{
						"assertions": "1",
						"class":      "SampleTest",
						"classname":  "SampleTest",
						"file":       "/untitled/tests/SampleTest.php",
						"line":       "7",
						"name":       "testA",
						"time":       "0.005917",
					},
				}

				assertEqual(t, testcase, suite.Tests[0])

				assertLen(t, suite.Suites[1].Suites, 0)
				assertLen(t, suite.Suites[1].Tests, 3)
				assertEqual(t, "testC with data set #0", suite.Suites[1].Tests[0].Name)

				// checking recursive aggregation
				suites[0].Aggregate()
				actualTotals := suites[0].Totals
				expectedTotals := Totals{
					Tests:    7,
					Passed:   4,
					Skipped:  0,
					Failed:   3,
					Error:    0,
					Duration: 8489 * time.Microsecond,
				}
				assertEqual(t, expectedTotals, actualTotals)
			},
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("#%d - %s", index+1, test.title)

		t.Run(name, func(t *testing.T) {
			suites, err := IngestFile(test.filename)
			assertNoError(t, err)
			test.check(t, suites)
		})
	}
}
