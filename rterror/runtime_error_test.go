// Copyright 2020 Tymoteusz Blazejczyk
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rterror_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/tymonx/go-error/rterror"
	"gitlab.com/tymonx/go-formatter/formatter"
)

func ExampleRuntimeError_withoutArguments() {
	err := rterror.New("Error message")

	fmt.Println(err)
	// Output: runtime_error_test.go:28:ExampleRuntimeError_withoutArguments(): Error message
}

func ExampleRuntimeError_withArguments() {
	err := rterror.New("Error message {p1} - {p0}", 3, "foo")

	fmt.Println(err)
	// Output: runtime_error_test.go:35:ExampleRuntimeError_withArguments(): Error message foo - 3
}

func ExampleRuntimeError_setFormat() {
	err := rterror.New("Error message {p1} - {p0}", 5, "bar").SetFormat("#{.Package | base}.{.FunctionBase}: '{.Message}'")

	fmt.Println(err)
	// Output: #rterror_test.ExampleRuntimeError_setFormat: 'Error message bar - 5'
}

func ExampleRuntimeError_unwrap() {
	wrapped := rterror.New("wrapped error").SetFormat("{.Message}")

	err := rterror.New("Error message", 5, wrapped)

	fmt.Println(errors.Is(err, wrapped))
	fmt.Println(err)
	// Output:
	// true
	// runtime_error_test.go:51:ExampleRuntimeError_unwrap(): Error message 5 wrapped error
}

func ExampleNewSkipCaller() {
	MyNewError := func(message string, arguments ...interface{}) *rterror.RuntimeError {
		return rterror.NewSkipCaller(1, message, arguments...)
	}

	err := MyNewError("Error message {p1}", "caller", "skip")

	fmt.Println(err)
	// Output: runtime_error_test.go:65:ExampleNewSkipCaller(): Error message skip caller
}

func TestRuntimeError(test *testing.T) {
	err := rterror.New("Error message", 5)

	assert.NotNil(test, err)
	assert.Error(test, err)
	assert.Equal(test, "runtime_error_test.go:72:TestRuntimeError(): Error message 5", err.Error())
}

func TestRuntimeLine(test *testing.T) {
	assert.NotZero(test, rterror.New("Error message").Line())
}

func TestRuntimeFile(test *testing.T) {
	assert.Contains(test, rterror.New("Error message").File(), "go-error/rterror/runtime_error_test.go")
}

func TestRuntimeFunction(test *testing.T) {
	assert.Equal(test, "gitlab.com/tymonx/go-error/rterror_test.TestRuntimeFunction", rterror.New("Error message").Function())
}

func TestRuntimeFunctionBase(test *testing.T) {
	assert.Equal(test, "TestRuntimeFunctionBase", rterror.New("Error message").FunctionBase())
}

func TestRuntimePackage(test *testing.T) {
	assert.Equal(test, "gitlab.com/tymonx/go-error/rterror_test", rterror.New("Error message").Package())
}

func TestRuntimeArguments(test *testing.T) {
	want := []interface{}{"5", 3, true, nil, 4.5}

	assert.Equal(test, want, rterror.New("Error message", want...).Arguments())
}

func TestRuntimeSetFormat(test *testing.T) {
	assert.Equal(test, "-> Error message", rterror.New("Error message").SetFormat("-> {.Message}").Error())
}

func TestRuntimeGetFormat(test *testing.T) {
	assert.Equal(test, rterror.DefaultFormat, rterror.New("Error message").GetFormat())
}

func TestRuntimeResetFormat(test *testing.T) {
	assert.Equal(test, rterror.DefaultFormat, rterror.New("Error message").SetFormat("X").ResetFormat().GetFormat())
}

func TestRuntimeGetFormatter(test *testing.T) {
	want := formatter.New()

	assert.NotNil(test, rterror.New("").GetFormatter())
	assert.Equal(test, want, rterror.New("").SetFormatter(want).GetFormatter())
}

func TestRuntimeUnwrap(test *testing.T) {
	err := rterror.New("error")

	assert.Equal(test, err, rterror.New("Error message", 5, err, 3).Unwrap())
}

func TestRuntimeUnwrapNil(test *testing.T) {
	assert.Nil(test, rterror.New("Error message").Unwrap())
}

func TestRuntimeMessageFailback(test *testing.T) {
	want := "Error message {invalid}"

	assert.Equal(test, want, rterror.New(want, 3, "foo").Message())
}

func TestRuntimeErrorFailback(test *testing.T) {
	want := "Error message {p1} {p0}"

	assert.Equal(test, want, rterror.New(want, 3, "foo").SetFormat("{invalid}").Error())
}

func TestRuntimeArgumets(test *testing.T) {
	err := rterror.New("error", 3, "test", 0.3, true)

	assert.Len(test, err.Arguments(), 4)
}
