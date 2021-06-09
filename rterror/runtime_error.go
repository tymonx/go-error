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

package rterror

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"gitlab.com/tymonx/go-formatter/formatter"
)

// These constants define default values for runtime error.
const (
	SkipCall      = 1
	DefaultIndent = "`--"
	DefaultFormat = `{cyan | bright}{.Package}{reset}:{bold}{cyan}{.FileBase}{reset}:{bold}{magenta}{.Line}{reset}:` +
		`{bold}{blue | bright}{.FunctionBase}(){reset}: {.String}`

	indentSize = len(DefaultIndent)
)

// RuntimeError defines a runtime error with message string formatted using
// "replacement fields" surrounded by curly braces {} format strings from
// the Go Formatter library. It contains line number, file path and function name
// from where a runtime error was called.
type RuntimeError struct {
	pc         [1]uintptr
	_message   string
	format     string
	formatter  *formatter.Formatter
	_arguments []interface{}
	err        error
}

// New creates a new runtime error object with message string formatted using
// "replacement fields" surrounded by curly braces {} format strings, line number,
// file path and function name from where the New() function was called.
func New(message string, arguments ...interface{}) *RuntimeError {
	return NewSkipCaller(SkipCall, message, arguments...)
}

// NewSkipCaller creates a new runtime error object with message string formatted using
// "replacement fields" surrounded by curly braces {} format strings, line number,
// file path and function name from where the NewSkipCaller() function was called.
// The argument skip is the number of stack frames to ascend,
// with 0 identifying the caller of NewSkipCaller.
func NewSkipCaller(skip int, message string, arguments ...interface{}) *RuntimeError {
	r := &RuntimeError{
		format:     DefaultFormat,
		formatter:  formatter.New(),
		_message:   message,
		_arguments: arguments,
	}

	runtime.Callers((SkipCall + SkipCall + skip), r.pc[:])

	return r
}

// Message returns unformatted error message.
func (r *RuntimeError) Message() string {
	return r._message
}

// Arguments returns error arguments.
func (r *RuntimeError) Arguments() []interface{} {
	return r._arguments
}

// Line returns line number.
func (r *RuntimeError) Line() int {
	return r.frame().Line
}

// File returns file absolute path.
func (r *RuntimeError) File() string {
	return r.frame().File
}

// FileBase returns file base path.
func (r *RuntimeError) FileBase() string {
	return filepath.Base(r.File())
}

// Function returns function full name.
func (r *RuntimeError) Function() string {
	return r.frame().Function
}

// FunctionBase returns function base name.
func (r *RuntimeError) FunctionBase() string {
	function := r.Function()

	if index := strings.LastIndexByte(function, '/'); index != -1 {
		function = function[index+1:]
	}

	if index := strings.IndexByte(function, '.'); index != -1 {
		function = function[index+1:]
	}

	return function
}

// Package returns full package path.
func (r *RuntimeError) Package() string {
	_package := r.Function()
	function := r.FunctionBase()

	return _package[:len(_package)-len(function)-1]
}

// PackageBase returns package name.
func (r *RuntimeError) PackageBase() string {
	_package := r.Package()

	if index := strings.LastIndexByte(_package, '/'); index != -1 {
		_package = _package[index+1:]
	}

	return _package
}

// SetFormat sets error message format string for formatter.
func (r *RuntimeError) SetFormat(format string) *RuntimeError {
	r.format = format
	return r
}

// GetFormat returns error message format string for formatter.
func (r *RuntimeError) GetFormat() string {
	return r.format
}

// ResetFormat resets error message format string for formatter to default value.
func (r *RuntimeError) ResetFormat() *RuntimeError {
	r.format = DefaultFormat
	return r
}

// SetFormatter sets formatter.
func (r *RuntimeError) SetFormatter(f *formatter.Formatter) *RuntimeError {
	r.formatter = f
	return r
}

// GetFormatter returns formatter.
func (r *RuntimeError) GetFormatter() *formatter.Formatter {
	return r.formatter
}

// String returns formatted error message string.
func (r *RuntimeError) String() string {
	if formatted, err := r.formatter.Format(r._message, r._arguments...); err == nil {
		return formatted
	}

	return r._message // Failback
}

// MarshalText encodes runtime error to text.
func (r *RuntimeError) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

// MarshalJSON encodes runtime error to JSON.
func (r *RuntimeError) MarshalJSON() ([]byte, error) {
	return json.Marshal(&marshal{
		Line:      r.Line(),
		File:      r.File(),
		Function:  r.Function(),
		Message:   r._message,
		Arguments: r._arguments,
	})
}

// Error returns formatted error message string.
//
// With wrapped errors it returns:
//
//  <error>
//  `--<error>
//     `--<error>
//        `--<error>
func (r *RuntimeError) Error() (result string) {
	var builder strings.Builder

	fmt.Fprint(&builder, r.TopError())

	for level, err := 0, r.err; err != nil; level, err = (level + 1), errors.Unwrap(err) {
		var message string

		if e, ok := err.(*RuntimeError); ok {
			message = e.TopError()
		} else {
			message = err.Error()
		}

		indent := indentSize * level

		builder.Grow(1 + indent + indentSize + len(message))

		fmt.Fprintln(&builder)

		for i := 0; i < indent; i++ {
			fmt.Fprint(&builder, " ") // indention
		}

		fmt.Fprint(&builder, DefaultIndent, message)
	}

	return builder.String()
}

// TopError returns top error message without any wrapped error messages.
//
// With wrapped errors it simple returns:
//
//  <error>
func (r *RuntimeError) TopError() string {
	if formatted, err := formatter.Format(r.format, r); err == nil {
		return formatted
	}

	return r._message // Failback
}

// Wrap wraps provided error into runtime error.
func (r *RuntimeError) Wrap(err error) *RuntimeError {
	r.err = err
	return r
}

// Unwrap returns wrapped error.
func (r *RuntimeError) Unwrap() error {
	return r.err
}

func (r *RuntimeError) frame() *runtime.Frame {
	frame, _ := runtime.CallersFrames(r.pc[:]).Next()
	return &frame
}
