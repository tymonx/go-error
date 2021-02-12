// Copyright 2021 Tymoteusz Blazejczyk
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
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/tymonx/go-error/rterror"
)

func TestIsTemporaryTrue(test *testing.T) {
	assert.True(test, rterror.IsTemporary(rterror.New("temporary", syscall.EAGAIN)))
}

func TestIsTemporaryFalse(test *testing.T) {
	assert.False(test, rterror.IsTemporary(rterror.New("temporary")))
}

func TestIsTimeoutTrue(test *testing.T) {
	assert.True(test, rterror.IsTimeout(rterror.New("timeout", syscall.ETIMEDOUT)))
}

func TestIsTimeoutFalse(test *testing.T) {
	assert.False(test, rterror.IsTimeout(rterror.New("timeout")))
}
