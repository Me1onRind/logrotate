package ticker

import (
	"time"
)

func CalRotateTimeDuration(now time.Time, duration time.Duration) time.Duration {
	nowUnixNao := now.UnixNano()
	NanoSecond := duration.Nanoseconds()
	nextRotateTime := NanoSecond - (nowUnixNao % NanoSecond)
	return time.Duration(nextRotateTime)
}
