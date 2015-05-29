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
	"errors"
	"fmt"
	"github.com/DirtyHairy/arduino-butler/go/util"
	"strconv"
)

type T string

func (ip T) String() string {
	return string(ip)
}

func (ip *T) Set(value string) error {
	if value == "" {
		return nil
	}

	rx := util.CompileRegex("^(\\d{1,3})\\.(\\d{1,3})\\.(\\d{1,3})\\.(\\d{1,3})$")

	matches := rx.FindStringSubmatch(value)

	if matches == nil {
		return errors.New(fmt.Sprintf("'%s' is not an IP", value))
	}

	for _, component := range matches[1:] {
		_, err := strconv.ParseUint(component, 10, 8)

		if err != nil {
			return errors.New(fmt.Sprintf("'%s' is not a valid IP", value))
		}
	}

	*ip = T(value)

	return nil
}

func Create() T {
	ip := T("")

	return ip
}
