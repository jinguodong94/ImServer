package utils

import "time"

var (
	TimeUtils  = &timeUtils{}
)

type timeUtils struct {
}

func (t *timeUtils) GetCurrentTime() uint64 {
	return uint64(time.Now().Unix())
}


