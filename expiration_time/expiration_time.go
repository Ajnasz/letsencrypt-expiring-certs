package expiration_time

import (
	"errors"
	"time"
)

func GetDefaultExpireTime() time.Time {
	return time.Now().AddDate(0, 0, 14).Truncate(time.Hour)
}

var dateFormats = []string{
	time.UnixDate,
	time.RFC3339,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123Z,
}

func getUserDefinedExpireTime(expireTime string) (time.Time, error) {
	for _, format := range dateFormats {
		expire, err := time.Parse(format, expireTime)
		if err == nil {
			return expire, nil
		}
	}

	return time.Time{}, errors.New("Invalid date format")
}

func GetExpireTime(expireTime string) (time.Time, error) {
	if expireTime == "" {
		return GetDefaultExpireTime(), nil
	} else {
		return getUserDefinedExpireTime(expireTime)
	}
}
