// This file has code copied from the "randstr" Go module, which can be found at
// https://github.com/thanhpk/randsr. The original code is licensed under the MIT
// license, which a copy can be found at https://github.com/thanhpk/randstr/blob/master/LICENSE
// and is provided below:
//
// # The MIT License
//
// Copyright (c) 2010-2018 Google, Inc. http://angularjs.org
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package randstr provides basic functions for generating random bytes, string
package randstr

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
)

// HexChars holds a string containing all characters used in a hexadecimal value.
const HexChars = "0123456789abcdef"

// NewHex generates a new Hexadecimal string with length of n
//
// Example: 67aab2d956bd7cc621af22cfb169cba8
func NewHex(n int) (string, error) { return New(n, HexChars) }

// New generates a random string using only letters provided in the letters parameter.
//
// If the letters parameter is omitted, this function will use HexChars instead.
func New(n int, chars ...string) (string, error) {
	runes := []rune(HexChars)
	if len(chars) > 0 {
		runes = []rune(chars[0])
	}

	var b bytes.Buffer
	b.Grow(n)
	l := uint32(len(runes))
	for range n {
		by, err := Bytes(4)
		if err != nil {
			return "", err
		}
		b.WriteRune(runes[binary.BigEndian.Uint32(by)%l])
	}
	return b.String(), nil
}

// Bytes generates n random bytes
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}
