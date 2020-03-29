package utils

import "time"

var LocalFormat = "2006-01-02 15:04:05"

func TimeNowString() string {
	return time.Now().Format(LocalFormat)
}
