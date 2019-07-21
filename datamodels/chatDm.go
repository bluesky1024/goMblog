package datamodels

import (
	"time"
)

type BarrageInfo struct {
	Uid        int64
	message    string
	CreateTime time.Time
	VideoTime  time.Time
}
