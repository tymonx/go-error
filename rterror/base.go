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
	"path/filepath"
	"runtime"
	"strings"

	"gitlab.com/tymonx/go-formatter/formatter"
)

type base struct {
	pc        uintptr
	format    string
	message   string
	formatter *formatter.Formatter
	arguments []interface{}
}

// Line returns line number.
func (b *base) Line() int {
	_, line := runtime.FuncForPC(b.pc).FileLine(b.pc)
	return line
}

// File returns full file path.
func (b *base) File() string {
	file, _ := runtime.FuncForPC(b.pc).FileLine(b.pc)
	return file
}

// Function returns function name.
func (b *base) Function() string {
	return filepath.Ext(runtime.FuncForPC(b.pc).Name())[1:]
}

// Package returns full package path.
func (b *base) Package() string {
	name := runtime.FuncForPC(b.pc).Name()
	return strings.TrimSuffix(name, filepath.Ext(name))
}

// ProgramCounter returns program counter.
func (b *base) ProgramCounter() uintptr {
	return b.pc
}

// Arguments returns arguments.
func (b *base) Arguments() []interface{} {
	return b.arguments
}

// Message returns formatted error message string.
func (b *base) Message() string {
	formatted, err := b.formatter.Format(b.message, b.arguments...)

	if err != nil {
		// Failback
		formatted = b.message
	}

	return formatted
}

func (b *base) formatMessage() string {
	formatted, err := formatter.Format(b.format, b)

	if err != nil {
		// Failback
		formatted = b.message
	}

	return formatted
}
