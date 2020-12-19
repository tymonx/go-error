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
	"path/filepath"
	"runtime"
	"strings"

	"gitlab.com/tymonx/go-formatter/formatter"
)

// These constants define default values for runtime error.
const (
	SkipCall = 1
)

// DefaultFormat defines default error message format.
var DefaultFormat = `{bold}{cyan}{.FileBase}{reset}:{bold}{magenta}{.Line}{reset}:` + // nolint: gochecknoglobals
	`{bold}{blue | bright}{.FunctionBase}(){reset}: {.String}`

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

	runtime.Callers(skip+2*SkipCall, r.pc[:])

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
func (r *RuntimeError) Line() (value int) {
	if function := runtime.FuncForPC(r.pc[0]); function != nil {
		_, value = function.FileLine(r.pc[0])
	}

	return value
}

// File returns file absolute path.
func (r *RuntimeError) File() (name string) {
	if function := runtime.FuncForPC(r.pc[0]); function != nil {
		name, _ = function.FileLine(r.pc[0])
	}

	return name
}

// FileBase returns file base path.
func (r *RuntimeError) FileBase() string {
	return filepath.Base(r.File())
}

// Function returns function full name.
func (r *RuntimeError) Function() (name string) {
	if function := runtime.FuncForPC(r.pc[0]); function != nil {
		name = function.Name()
	}

	return name
}

// FunctionBase returns function base name.
func (r *RuntimeError) FunctionBase() string {
	return strings.TrimPrefix(filepath.Ext(r.Function()), ".")
}

// Package returns full package path.
func (r *RuntimeError) Package() string {
	function := r.Function()
	return strings.TrimSuffix(function, filepath.Ext(function))
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
	formatted, err := r.formatter.Format(r._message, r._arguments...)

	if err != nil {
		// Failback
		formatted = r._message
	}

	return formatted
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
func (r *RuntimeError) Error() string {
	formatted, err := formatter.Format(r.format, r)

	if err != nil {
		// Failback
		formatted = r._message
	}

	return formatted
}

// Unwrap returns wrapped error.
func (r *RuntimeError) Unwrap() error {
	for _, argument := range r._arguments {
		if err, ok := argument.(error); ok {
			return err
		}
	}

	return nil
}
