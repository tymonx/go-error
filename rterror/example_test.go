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

	"gitlab.com/tymonx/go-error/rterror"
)

func ExampleRuntimeError_withoutArguments() {
	err := rterror.New("Error message")

	fmt.Println(err)
	// Output: gitlab.com/tymonx/go-error/rterror_test:example_test.go:25:ExampleRuntimeError_withoutArguments(): Error message
}

func ExampleRuntimeError_withArguments() {
	err := rterror.New("Error message {p1} - {p0}", 3, "foo")

	fmt.Println(err)
	// Output: gitlab.com/tymonx/go-error/rterror_test:example_test.go:32:ExampleRuntimeError_withArguments(): Error message foo - 3
}

func ExampleRuntimeError_setFormat() {
	err := rterror.New("Error message {p1} - {p0}", 5, "bar").SetFormat("#{.Package | base}.{.FunctionBase}: '{.String}'")

	fmt.Println(err)
	// Output: #rterror_test.ExampleRuntimeError_setFormat: 'Error message bar - 5'
}

func ExampleRuntimeError_unwrap() {
	wrapped := rterror.New("wrapped error").SetFormat("{.Message}")

	err := rterror.New("Error message", 5).Wrap(wrapped)

	fmt.Println(errors.Is(err, wrapped))
	fmt.Println(err)
	// Output:
	// true
	// gitlab.com/tymonx/go-error/rterror_test:example_test.go:48:ExampleRuntimeError_unwrap(): Error message 5
	// `--wrapped error
}

func ExampleNewSkipCaller() {
	MyNewError := func(message string, arguments ...interface{}) *rterror.RuntimeError {
		return rterror.NewSkipCaller(1, message, arguments...)
	}

	err := MyNewError("Error message {p1}", "caller", "skip")

	fmt.Println(err)
	// Output: gitlab.com/tymonx/go-error/rterror_test:example_test.go:63:ExampleNewSkipCaller(): Error message skip caller
}
