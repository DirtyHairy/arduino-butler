package util

import (
    "testing"
    "time"
)

func TestValidDuration(t *testing.T) {
    var d DurationValue

    err := d.Set("1h")

    if err != nil {
        t.Errorf("error during set: %v", err)
    }

    if d.String() != "1h0m0s" {
        t.Errorf("reserialization failed; should habe been 1h0m0s, got %s", d.String())
    }

    if d.Value() != 1 * time.Hour {
        t.Errorf("duration should be 1h, got %s", d.Value())
    }
}

func TestInvalidDuration(t *testing.T) {
    var d DurationValue

    if d.Set("+++") == nil {
        t.Error("setting value to an invalid duration should error")
    }
}
