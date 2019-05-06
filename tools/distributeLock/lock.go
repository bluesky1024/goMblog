package distributeLock

import (
	"context"
	"time"
)

type LockSrv interface {
	TryGetLock(lockName string, timeWait time.Duration, timeHoldLock time.Duration) (lock DisLock, err error)
}

type DisLock interface {
	ExtendLock(ctx context.Context) (err error)
	UnLock() (err error)
}
