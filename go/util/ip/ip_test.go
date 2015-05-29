/**
 * The MIT License (MIT)
 *
 * Copyright (c) 2015 Christian Speckner
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package ip

import (
	"testing"
)

func assertInvalid(address string, t *testing.T) {
	ip := Create()

	if err := ip.Set(address); err == nil {
		t.Errorf("'%s' should be rejected", address)
	}
}

func assertValid(address string, t *testing.T) {
	ip := Create()

	if err := ip.Set(address); err != nil {
		t.Errorf("'%s' should be accepted, got %v instead", address, err)
	}
}

func TestInvalidIp(t *testing.T) {
	testcases := []string{
		"foo",
		"foo.0.0.1",
		"127.foo.0.1",
		"127.0.foo.1",
		"127.0.0.foo",
		"127.0.0",
		"127x0.0.1",
		"127.0x0.1",
		"127.0.0x1",
		"256.0.0.1",
		"0.256.0.1",
		"0.0.256.1",
		"0.0.0.256",
	}

	for _, testcase := range testcases {
		assertInvalid(testcase, t)
	}
}

func TestValidIp1(t *testing.T) {
	testcases := []string{
		"127.0.0.1",
		"255.0.0.1",
		"0.255.0.1",
		"0.0.255.1",
		"0.0.0.255",
		"",
	}

	for _, testcase := range testcases {
		assertValid(testcase, t)
	}
}

func TestSerialization(t *testing.T) {
	ip := Create()
	ip.Set("127.0.0.1")

	if ip.String() != "127.0.0.1" {
		t.Error("Serialization should return the IP")
	}
}
