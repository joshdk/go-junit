// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package junit

import (
	"bytes"
	"reflect"
	"testing"
)

// assertEqual is a testing helper function which asserts that the given
// objects are equal.
func assertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("objects were not equal: \n"+
			"expected: %v\n"+
			"actual  : %v", expected, actual)
	}
}

// assertEqualBytes is a testing helper function which asserts that the given
// byte slices are equal.
func assertEqualBytes(t *testing.T, expected, actual []byte) {
	t.Helper()
	if !bytes.Equal(expected, actual) {
		t.Fatalf("byte slices were not equal: \n"+
			"expected: %s\n"+
			"actual  : %s", string(expected), string(actual))
	}
}

// assertLen is a testing helper function which asserts that the given object
// has the expected length.
func assertLen(t *testing.T, object interface{}, expected int) {
	t.Helper()
	actual := reflect.ValueOf(object).Len()
	if actual != expected {
		t.Fatalf("lengths were not equal: \n"+
			"expected: %d\n"+
			"actual  : %d", expected, actual)
	}
}

// assertError is a testing helper function which asserts that the given error
// is not nil and has the expected error text.
func assertError(t *testing.T, actual error, expected string) {
	t.Helper()
	if actual == nil || actual.Error() != expected {
		t.Fatalf("error was not equal: \n"+
			"expected: %v\n"+
			"actual  : %v", expected, actual)
	}
}

// assertNoError is a testing helper function which asserts that the given
// error is nil.
func assertNoError(t *testing.T, actual error) {
	t.Helper()
	if actual != nil {
		t.Fatalf("error was not nil: \n"+
			"expected: no error\n"+
			"actual  : %v", actual)
	}
}
