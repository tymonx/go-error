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
	"runtime"

	"gitlab.com/tymonx/go-formatter/formatter"
)

// These constants define default values for runtime error.
const (
	DefaultFormat = "{.File | base}:{.Line}:{.Package | base}.{.Function}(): {.Message}"
)

// RuntimeError defines a runtime error with message string formatted using
// "replacement fields" surrounded by curly braces {} format strings from
// the Go Formatter library. It contains line number, file path and function name
// from where a runtime error was called.
type RuntimeError struct {
	base base
}

// New creates a new runtime error object with message string formatted using
// "replacement fields" surrounded by curly braces {} format strings, line number,
// file path and function name from where the New() function was called.
func New(message string, arguments ...interface{}) *RuntimeError {
	return NewSkipCaller(1, message, arguments...)
}

// NewSkipCaller creates a new runtime error object with message string formatted using
// "replacement fields" surrounded by curly braces {} format strings, line number,
// file path and function name from where the NewSkipCaller() function was called.
// The argument skip is the number of stack frames to ascend,
// with 0 identifying the caller of NewSkipCaller.
func NewSkipCaller(skip int, message string, arguments ...interface{}) *RuntimeError {
	pc, _, _, _ := runtime.Caller(skip + 1) // nolint: dogsled

	return &RuntimeError{
		base: base{
			pc:        pc,
			format:    DefaultFormat,
			message:   message,
			formatter: formatter.New(),
			arguments: arguments,
		},
	}
}

// Line returns line number.
func (r *RuntimeError) Line() int {
	return r.base.Line()
}

// File returns full file path.
func (r *RuntimeError) File() string {
	return r.base.File()
}

// Function returns function name.
func (r *RuntimeError) Function() string {
	return r.base.Function()
}

// Package returns full package path.
func (r *RuntimeError) Package() string {
	return r.base.Package()
}

// ProgramCounter returns program counter.
func (r *RuntimeError) ProgramCounter() uintptr {
	return r.base.ProgramCounter()
}

// Arguments returns arguments.
func (r *RuntimeError) Arguments() []interface{} {
	return r.base.Arguments()
}

// SetFormat sets error message format string for formatter.
func (r *RuntimeError) SetFormat(format string) *RuntimeError {
	r.base.format = format
	return r
}

// GetFormat returns error message format string for formatter.
func (r *RuntimeError) GetFormat() string {
	return r.base.format
}

// ResetFormat resets error message format string for formatter to default value.
func (r *RuntimeError) ResetFormat() *RuntimeError {
	r.base.format = DefaultFormat
	return r
}

// SetFormatter sets formatter.
func (r *RuntimeError) SetFormatter(f *formatter.Formatter) *RuntimeError {
	r.base.formatter = f
	return r
}

// GetFormatter returns formatter.
func (r *RuntimeError) GetFormatter() *formatter.Formatter {
	return r.base.formatter
}

// Message returns formatted error message string.
func (r *RuntimeError) Message() string {
	return r.base.Message()
}

// Error returns formatted error message string.
func (r *RuntimeError) Error() string {
	return r.base.formatMessage()
}

// Unwrap returns wrapped error.
func (r *RuntimeError) Unwrap() error {
	for _, argument := range r.base.arguments {
		if err, ok := argument.(error); ok {
			return err
		}
	}

	return nil
}
