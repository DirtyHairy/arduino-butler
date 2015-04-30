package ip

import (
	"errors"
	"fmt"
	"regexp"
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

	rx, err := regexp.Compile("^(\\d{1,3})\\.(\\d{1,3})\\.(\\d{1,3})\\.(\\d{1,3})$")

	if err != nil {
		panic(err)
	}

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
