package chron

import (
	"math"
	"math/rand"
	"time"
)

func ThisUnixNano() int64 {
	return time.Now().UnixNano() + rand.Int63n(1000)
}

func FromUnixNano(t int64) Chron {
	secs := int64(math.Floor(float64(t) / 1e9))
	nsecs := t - (int64(secs) * 1e9)
	return TimeOf(time.Unix(secs, nsecs))
}
