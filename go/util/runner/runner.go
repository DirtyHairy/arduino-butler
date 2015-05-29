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

package runner

type T struct {
	commandChannel chan cmdWrapper
	controlChannel chan int

	Execute func(interface{}) error
}

type cmdWrapper struct {
	cmd           interface{}
	resultChannel chan error
}

func (runner *T) Start() error {
	if runner.commandChannel != nil {
		return Error("already started")
	}

	runner.commandChannel = make(chan cmdWrapper)
	runner.controlChannel = make(chan int)

	go func() {
		for {
			cmd, ok := <-runner.commandChannel

			if !ok {
				break
			}

			cmd.resultChannel <- runner.Execute(cmd.cmd)
		}

		runner.controlChannel <- 1
	}()

	return nil
}

func (runner *T) Stop() error {
	if runner.commandChannel == nil {
		return Error("already stopped")
	}

	close(runner.commandChannel)

	_ = <-runner.controlChannel

	runner.commandChannel = nil
	runner.controlChannel = nil

	return nil
}

func (runner *T) Dispatch(cmd interface{}) error {
	if runner.commandChannel == nil {
		return Error("not running")
	}

	resultChannel := make(chan error)

	runner.commandChannel <- cmdWrapper{cmd, resultChannel}

	return <-resultChannel
}
