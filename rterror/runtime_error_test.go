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
	"encoding/json"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/tymonx/go-error/rterror"
	"gitlab.com/tymonx/go-formatter/formatter"
)

func TestRuntimeError(test *testing.T) {
	err := rterror.New("Error message", 5)

	assert.NotNil(test, err)
	assert.Error(test, err)
	assert.Equal(test, "gitlab.com/tymonx/go-error/rterror_test:runtime_error_test.go:28:TestRuntimeError(): Error message 5", err.Error())
}

func TestRuntimeErrorLine(test *testing.T) {
	assert.NotZero(test, rterror.New("Error message").Line())
}

func TestRuntimeErrorFile(test *testing.T) {
	assert.Contains(test, rterror.New("Error message").File(), "go-error/rterror/runtime_error_test.go")
}

func TestRuntimeErrorFunction(test *testing.T) {
	assert.Equal(test, "gitlab.com/tymonx/go-error/rterror_test.TestRuntimeErrorFunction", rterror.New("Error message").Function())
}

func TestRuntimeErrorFunctionBase(test *testing.T) {
	assert.Equal(test, "TestRuntimeErrorFunctionBase", rterror.New("Error message").FunctionBase())
}

func TestRuntimeErrorPackage(test *testing.T) {
	assert.Equal(test, "gitlab.com/tymonx/go-error/rterror_test", rterror.New("Error message").Package())
}

func TestRuntimeErrorArguments(test *testing.T) {
	want := []interface{}{"5", 3, true, nil, 4.5}

	assert.Equal(test, want, rterror.New("Error message", want...).Arguments())
}

func TestRuntimeErrorSetFormat(test *testing.T) {
	assert.Equal(test, "-> Error message", rterror.New("Error message").SetFormat("-> {.Message}").Error())
}

func TestRuntimeErrorGetFormat(test *testing.T) {
	assert.Equal(test, rterror.DefaultFormat, rterror.New("Error message").GetFormat())
}

func TestRuntimeErrorResetFormat(test *testing.T) {
	assert.Equal(test, rterror.DefaultFormat, rterror.New("Error message").SetFormat("X").ResetFormat().GetFormat())
}

func TestRuntimeErrorGetFormatter(test *testing.T) {
	want := formatter.New()

	assert.NotNil(test, rterror.New("").GetFormatter())
	assert.Equal(test, want, rterror.New("").SetFormatter(want).GetFormatter())
}

func TestRuntimeErrorUnwrap(test *testing.T) {
	err := rterror.New("error")

	assert.Equal(test, err, rterror.New("Error message", 5, 3).Wrap(err).Unwrap())
}

func TestRuntimeErrorUnwrapNil(test *testing.T) {
	assert.Nil(test, rterror.New("Error message").Unwrap())
}

func TestRuntimeErrorMessageFailback(test *testing.T) {
	want := "Error message {invalid}"

	assert.Equal(test, want, rterror.New(want, 3, "foo").String())
}

func TestRuntimeErrorFailback(test *testing.T) {
	want := "Error message {p1} {p0}"

	assert.Equal(test, want, rterror.New(want, 3, "foo").SetFormat("{invalid}").Error())
}

func TestRuntimeErrorArgumets(test *testing.T) {
	err := rterror.New("error", 3, "test", 0.3, true)

	assert.Len(test, err.Arguments(), 4)
}

func TestRuntimeErrorMarshalJSON(test *testing.T) {
	e := rterror.New("error", 4, nil)

	data, err := json.Marshal(e)

	assert.NoError(test, err)
	assert.NotEmpty(test, data)
}

func TestRuntimeErrorMarshalText(test *testing.T) {
	e := rterror.New("error", 5, nil)

	data, err := e.MarshalText()

	assert.NoError(test, err)
	assert.NotEmpty(test, data)
}

func TestRuntimeErrorWrap(test *testing.T) {
	assert.NotNil(test, rterror.New("error").Wrap(rterror.New("wrap")).Unwrap())
}

func TestRuntimeErrorNestedWrap(test *testing.T) {
	assert.NotEmpty(test, rterror.New("A").Wrap(rterror.New("B").Wrap(rterror.New("C").Wrap(syscall.EAGAIN))).Error())
}

type Struct struct{}

func (*Struct) Error() *rterror.RuntimeError {
	return rterror.New("error")
}

func TestRuntimeErrorStructFunction(test *testing.T) {
	assert.Equal(test, "gitlab.com/tymonx/go-error/rterror_test.(*Struct).Error", new(Struct).Error().Function())
}

func TestRuntimeErrorStructFunctionBase(test *testing.T) {
	assert.Equal(test, "(*Struct).Error", new(Struct).Error().FunctionBase())
}

func TestRuntimeErrorStructPackage(test *testing.T) {
	assert.Equal(test, "gitlab.com/tymonx/go-error/rterror_test", new(Struct).Error().Package())
}

func TestRuntimeErrorStructPackageBase(test *testing.T) {
	assert.Equal(test, "rterror_test", new(Struct).Error().PackageBase())
}
