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

package controls

type mockBackend struct {
	id          string
	state       bool
	running     bool
	host        string
	typ         string
	lastToggled uint

	startError  error
	stopError   error
	toggleError error
}

func (backend *mockBackend) Id() string {
	return backend.id
}

func (backend *mockBackend) setId(id string) {
	backend.id = id
}

func (backend *mockBackend) Start() error {
	if backend.startError != nil {
		return backend.startError
	}

	backend.running = true
	return nil
}

func (backend *mockBackend) Stop() error {
	if backend.stopError != nil {
		return backend.stopError
	}

	backend.running = false
	return nil
}

func (backend *mockBackend) Marshal() MarshalledBackend {
	return MarshalledBackend{
		Id:   backend.id,
		Type: backend.typ,
		Host: backend.host,
	}
}

func (backend *mockBackend) Toggle(switchIdx uint, state bool) error {
	if backend.toggleError != nil {
		return backend.toggleError
	}

	backend.lastToggled = switchIdx
	backend.state = state

	return nil
}
