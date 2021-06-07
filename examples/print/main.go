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

package main

import (
	"fmt"

	"gitlab.com/tymonx/go-error/rterror"
)

func error1() error {
	return rterror.New("my error message 1").Wrap(error2())
}

func error2() error {
	return rterror.New("my error message 2").Wrap(error3())
}

func error3() error {
	return rterror.New("my error message 3").Wrap(error4())
}

func error4() error {
	return rterror.New("my error message 4")
}

func main() {
	fmt.Println(error1())
}
