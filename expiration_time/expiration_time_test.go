package expiration_time

import (
	"testing"
	"time"
)

func TestGetDefaultExpireTime(t *testing.T) {
	actual := GetDefaultExpireTime()

	now := time.Now()

	actual = actual.Truncate(time.Hour)

	expected := now.Truncate(time.Hour).AddDate(0, 0, 14)

	if actual != expected {
		t.Fatal("Default expire time should be 2 weeks", actual, expected)
	}
}
