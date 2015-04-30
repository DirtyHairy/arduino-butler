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
