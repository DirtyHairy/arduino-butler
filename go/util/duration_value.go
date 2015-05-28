package util

import (
	"time"
)

type DurationValue time.Duration

func (d *DurationValue) String() string {
	return (*time.Duration)(d).String()
}

func (d *DurationValue) Set(value string) error {
	duration, err := time.ParseDuration(value)

	if err != nil {
		return err
	}

	*d = DurationValue(duration)

	return nil
}

func (d *DurationValue) Value() time.Duration {
	return time.Duration(*d)
}
