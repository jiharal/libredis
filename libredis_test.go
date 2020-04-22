package libredis

import (
	"testing"
)

func TestLibredis(t *testing.T) {
	uri := Options{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		MaxIdle:  100,
		Enabled:  true,
	}
	conn := Connect(uri)
	defer conn.Close()
}
