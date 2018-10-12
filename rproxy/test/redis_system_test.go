package system_test

import (
	"testing"
)

// We just need to test Redis Protocol on our Proxy server.
// All other cases are handled in the HTTP tests.
func TestTrivialRedisGets(t *testing.T) {
	tt := []struct {
		Key string
		Val string
	}{
		{"5", "abc"},
		{"5930", "123"},
		{"6", "2"},
		{"iauau", "123"},
		{"dasd", "456"},
		{"2asd", "789"},
	}

	for _, tc := range tt {
		if err := redisSet(tc.Key, tc.Val, redisURI); err != nil {
			t.Errorf("unexpected err happened: %v", err)
		}

		rpVal, err := redisGet(tc.Key, rproxyTCPURI)
		if err != nil {
			t.Errorf("unexpected err happened: %v", err)
		}

		if rpVal != tc.Val {
			t.Errorf("Val was %v but we expected %v", rpVal, tc.Val)
		}
	}
}
