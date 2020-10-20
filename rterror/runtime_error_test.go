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
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/tymonx/go-error/rterror"
	"gitlab.com/tymonx/go-formatter/formatter"
)

func TestRuntimeError(test *testing.T) {
	err := rterror.New("Error message", 5)

	assert.NotNil(test, err)
	assert.Error(test, err)
	assert.Equal(test, "runtime_error_test.go:26:rterror_test.TestRuntimeError(): Error message 5", err.Error())
}

func TestRuntimeLine(test *testing.T) {
	assert.NotZero(test, rterror.New("Error message").Line())
}

func TestRuntimeProgramCounter(test *testing.T) {
	assert.NotZero(test, rterror.New("Error message").ProgramCounter())
}

func TestRuntimeFile(test *testing.T) {
	assert.Contains(test, rterror.New("Error message").File(), "go-error/rterror/runtime_error_test.go")
}

func TestRuntimeBaseFile(test *testing.T) {
	assert.Equal(test, "runtime_error_test.go", rterror.New("Error message").BaseFile())
}

func TestRuntimeFunction(test *testing.T) {
	assert.Contains(test, rterror.New("Error message").Function(), "go-error/rterror_test.TestRuntimeFunction")
}

func TestRuntimeBaseFunction(test *testing.T) {
	assert.Equal(test, "rterror_test.TestRuntimeBaseFunction", rterror.New("Error message").BaseFunction())
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
