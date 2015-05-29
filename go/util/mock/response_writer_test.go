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

package mock

import (
	"os"
	"testing"
)

var response MockResponseWriter

func TestHeaders(t *testing.T) {
	header := response.Header()
	header.Set("foo", "bar")

	if header.Get("foo") != "bar" {
		t.Error("settings a header failed")
	}
}

func TestWrite(t *testing.T) {
	var data [10]byte

	written, err := response.Write(data[0:])

	if written != 10 {
		t.Error("byte count after write differs")
	}

	if err != nil {
		t.Error("write test failed")
	}
}

func TestWriteHeader(t *testing.T) {
	response.WriteHeader(666)

	if response.Code != 666 {
		t.Error("failed to set response code")
	}
}

func TestMain(m *testing.M) {
	response = CreateMockResponseWriter()

	os.Exit(m.Run())
}
