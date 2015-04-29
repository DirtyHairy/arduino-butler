package ip

import (
	"errors"
	"flag"
	"fmt"
	"regexp"
	"strconv"
)

type T interface {
	flag.Value
}

type implementation string

func (ip implementation) String() string {
	return string(ip)
}

func (ip *implementation) Set(value string) error {
	if value == "" {
		return nil
	}

	rx, _ := regexp.Compile("^(\\d{1,3})\\.(\\d{1,3})\\.(\\d{1,3})\\.(\\d{1,3})$")
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

	*ip = implementation(value)

	return nil
}

func Create() T {
	ip := implementation("")

	return &ip
}
