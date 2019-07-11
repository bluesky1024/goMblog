package redisLock

//基于redis实现的分布式锁工具
//保证锁的可重入性，首先需要确定同一个协程（mac地址+进程Id+线程Id+协程Id）
//可重入计划流产，go语言设计者认为可重入锁是个失败的设计，所以coroutineId无法获取

import (
	"context"
	"errors"
	"github.com/bluesky1024/goMblog/tools/distributeLock"
	"github.com/go-redis/redis"
	"reflect"
	"strconv"
	"time"
)

type redisLockSrv struct {
	RedisPool *redis.ClusterClient
}

type redisLock struct {
	LockId       string
	TimeExpire   int64
	TimeHoldLock time.Duration
}

var (
	RedisLockSrv *redisLockSrv
	lockKey      = "lockKey"
)

func InitRedisLockSrv(addrs []string) error {
	redisCluster := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addrs,
	})
	pong, err := redisCluster.Ping().Result()
	if err != nil || pong != "PONG" {
		return err
	}
	RedisLockSrv = &redisLockSrv{
		RedisPool: redisCluster,
	}
	return nil
}

func LoadSrv() distributeLock.LockSrv {
	return RedisLockSrv
}

func (s *redisLockSrv) TryGetLock(lockName string, timeWait time.Duration, timeHoldLock time.Duration) (distributeLock.DisLock, error) {

	lockId := lockName + lockKey
	//尝试SETNX
	timeExpire := time.Now().UnixNano() + int64(timeHoldLock)
	res, err := s.RedisPool.SetNX(lockId, timeExpire, 0).Result()
	if err != nil {
		return nil, err
	}
	//获取成功
	if res {
		return &redisLock{
			LockId:       lockId,
			TimeExpire:   timeExpire,
			TimeHoldLock: timeHoldLock,
		}, nil
	}

	//获取失败,开始检测当前持有锁的那位仁兄是否超时了，同时设置自身等待时间
	ctx := context.Background()
	ctxTimer, cancel := context.WithTimeout(ctx, timeWait)
	defer cancel()

	script := "if redis.call('GET', KEYS[1]) == ARGV[1] " +
		"then " +
		"return redis.call('SET', KEYS[1], ARGV[2]) " +
		"else " +
		"return 'fail' " +
		"end "

	for {
		//获取锁的值
		timeExpireOthersStr, err := s.RedisPool.Get(lockId).Result()
		if err != nil {
			//若获取不到，表示当前持有人已经释放了，再次尝试setnx
			if err.Error() == "redis: nil" {
				timeExpire = time.Now().UnixNano() + int64(timeHoldLock)
				res, err = s.RedisPool.SetNX(lockId, timeExpire, 0).Result()
				if err != nil {
					return nil, err
				}
				//获取成功
				if res {
					return &redisLock{
						LockId:       lockId,
						TimeExpire:   timeExpire,
						TimeHoldLock: timeHoldLock,
					}, nil
				}
				continue
			} else {
				return nil, err
			}
		}
		timeExpireOthers, _ := strconv.ParseInt(timeExpireOthersStr, 10, 64)
		//判断是否超时了
		if timeExpireOthers <= time.Now().UnixNano() {
			//尝试释放锁,并由自身获取锁
			timeExpire = time.Now().UnixNano() + int64(timeHoldLock)
			keys := []string{lockId}
			args := []interface{}{timeExpireOthers, timeExpire}
			getRes, err := s.RedisPool.Eval(script, keys, args...).Result()
			getRes = getRes.(string)
			if err != nil {
				return nil, err
			}
			if getRes == "OK" {
				return &redisLock{
					LockId:       lockId,
					TimeExpire:   timeExpire,
					TimeHoldLock: timeHoldLock,
				}, nil
			}
		}

		select {
		case <-ctxTimer.Done():
			return nil, errors.New("get lock timeout")
		default:
		}
		time.Sleep(1 * time.Millisecond)
	}
}

//启动一个协程监控当前客户端，若还存活根据redisLock.TimeHoldLock定期给锁续费
func (l *redisLock) ExtendLock(ctx context.Context) (err error) {
	go func(lock *redisLock, ctxCancel context.Context) {
		script := `local getRes=redis.call('GET', KEYS[1])
				if getRes
    				then
        				if getRes==ARGV[1]
            				then
                				return redis.call('SET',KEYS[1],ARGV[2])
							else
                				return "invalid"
        				end
    				else
        				return "not exist"
				end`
		for {
			time.Sleep(l.TimeHoldLock * 2 / 3)
			//检测是否还持有锁，若持有则续费并继续循环,若不再持有则退出该协程
			newTimeExpire := time.Now().UnixNano() + int64(l.TimeHoldLock)
			keys := []string{l.LockId}
			args := []interface{}{l.TimeExpire, newTimeExpire}
			delRes, err := RedisLockSrv.RedisPool.Eval(script, keys, args...).Result()
			delRes = delRes.(string)
			if err != nil {
				return
			}
			if delRes != "OK" {
				return
			}
			//更新锁的过期时间
			l.TimeExpire = newTimeExpire

			select {
			case <-ctxCancel.Done():
				return
			default:
			}
		}
	}(l, ctx)
	return nil
}

func (l *redisLock) UnLock() (err error) {
	script := "if redis.call('GET', KEYS[1]) == ARGV[1] " +
		"then " +
		"return redis.call('DEL', KEYS[1])" +
		"else " +
		"return 'fail' " +
		"end "
	keys := []string{l.LockId}
	args := []interface{}{l.TimeExpire}
	delRes, err := RedisLockSrv.RedisPool.Eval(script, keys, args...).Result()
	if err != nil {
		return errors.New("release lock fail:" + err.Error())
	}
	if reflect.TypeOf(delRes).Name() != "int64" || delRes.(int64) != 1 {
		return errors.New("release lock fail")
	}
	return nil
}
