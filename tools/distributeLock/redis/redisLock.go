package redisLock

//基于redis实现的分布式锁工具
//因为基于redis连接池，并发数受连接池最大连接数限制(连接数仅250左右，问题待定位...)

import (
	"context"
	"errors"
	"github.com/bluesky1024/goMblog/tools/distributeLock"
	"github.com/gomodule/redigo/redis"
	"time"
)

type redisLockSrv struct {
	RedisPool *redis.Pool
}

type redisLock struct {
	LockId   string
	TimeWait time.Duration
	TimeExpire int64
	TimeHoldLock time.Duration
}

var (
	RedisLockSrv *redisLockSrv
	lockKey      = "thisIsKeyForTrans"
)

func InitRedisLockSrv(addr string) error {
	pool := &redis.Pool{
		MaxActive:   5000,
		MaxIdle:     1000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			//if _, err := c.Do("AUTH", password); err != nil {
			//	c.Close()
			//	return nil, err
			//}
			//if _, err := c.Do("SELECT", db); err != nil {
			//	c.Close()
			//	return nil, err
			//}
			return c, nil
		},
	}
	RedisLockSrv = &redisLockSrv{
		RedisPool: pool,
	}
	return nil
}

func LoadSrv() distributeLock.LockSrv {
	return RedisLockSrv
}

func (s *redisLockSrv) TryGetLock(lockName string, timeWait time.Duration,timeHoldLock time.Duration) (distributeLock.DisLock, error) {
	conn := RedisLockSrv.RedisPool.Get()
	defer conn.Close()

	lockId := lockName + lockKey
	//尝试SETNX
	timeExpire := time.Now().UnixNano() + int64(timeHoldLock)
	res, err := redis.Int64(conn.Do("SETNX", lockId, timeExpire))
	if err != nil {
		return nil, err
	}
	//获取成功
	if res != 0 {
		return &redisLock{
			LockId:   lockId,
			TimeWait: timeWait,
			TimeExpire:timeExpire,
			TimeHoldLock:timeHoldLock,
		}, nil
	}

	//获取失败,开始检测当前持有锁的那位仁兄是否超时了，同时设置自身等待时间
	ctx := context.Background()
	ctxTimer,cancel := context.WithTimeout(ctx,timeWait)
	defer cancel()

	script :=	"if redis.call('GET', KEYS[1]) == ARGV[1] " +
					"then " +
						"return redis.call('SET', KEYS[1], ARGV[2]) " +
					"else " +
						"return 'fail' " +
				"end "
	redisScrip := redis.NewScript(1,script)

	for{
		//获取锁的值
		timeExpireOthers,err := redis.Int64(conn.Do("GET",lockId))
		if err != nil {
			//fmt.Println("try get",err.Error())
			//若获取不到，表示当前持有人已经释放了，再次尝试setnx
			if err.Error() == "redigo: nil returned" {
				timeExpire = time.Now().UnixNano() + int64(timeHoldLock)
				res, err = redis.Int64(conn.Do("SETNX", lockId, timeExpire))
				if err != nil {
					return nil, err
				}
				//获取成功
				if res != 0 {
					return &redisLock{
						LockId:   lockId,
						TimeWait: timeWait,
						TimeExpire:timeExpire,
						TimeHoldLock:timeHoldLock,
					}, nil
				}
				continue
			}else{
				return nil,err
			}
		}
		//判断是否超时了
		if timeExpireOthers <= time.Now().UnixNano(){
			//尝试释放锁,并由自身获取锁
			timeExpire = time.Now().UnixNano() + int64(timeHoldLock)
			getRes,err := redis.String(redisScrip.Do(conn,lockId,timeExpireOthers,timeExpire))
			if err != nil {
				return nil,err
			}
			if getRes == "OK" {
				return &redisLock{
					LockId:lockId,
					TimeWait:timeWait,
					TimeExpire:timeExpire,
					TimeHoldLock:timeHoldLock,
				},nil
			}
		}

		select {
		case <-ctxTimer.Done():
			return nil,errors.New("get lock timeout")
		default:
		}
		time.Sleep(1 * time.Millisecond)
	}
}

//启动一个协程监控当前客户端，若还存活根据redisLock.TimeHoldLock定期给锁续费
func (l *redisLock) ExtendLock(ctx context.Context) (err error){
	go func(lock *redisLock,ctxCancel context.Context) {
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
		redisScrip := redis.NewScript(1,script)
		conn := RedisLockSrv.RedisPool.Get()
		defer conn.Close()
		for{
			time.Sleep(l.TimeHoldLock*2/3)
			//检测是否还持有锁，若持有则续费并继续循环,若不再持有则退出该协程
			newTimeExpire := time.Now().UnixNano() + int64(l.TimeHoldLock)
			delRes,err := redis.String(redisScrip.Do(conn,l.LockId,l.TimeExpire,newTimeExpire))
			if err != nil {
				return
			}
			if delRes != "OK" {
				return
			}
			//更新锁的过期时间
			l.TimeExpire = newTimeExpire

			select{
				case <-ctxCancel.Done():
					return
			default:
			}
		}
	}(l,ctx)
	return nil
}

func (l *redisLock) UnLock() (err error){
	conn := RedisLockSrv.RedisPool.Get()
	defer conn.Close()

	script :=	"if redis.call('GET', KEYS[1]) == KEYS[2] " +
		"then " +
		"return redis.call('DEL', KEYS[1])" +
		"else " +
		"return 0 " +
		"end "
	redisScrip := redis.NewScript(2,script)

	delRes,err := redis.String(redisScrip.Do(conn,l.LockId,l.TimeExpire))
	if err != nil {
		return  errors.New("release lock fail:" + err.Error())
	}
	if delRes != "OK" {
		return errors.New("release lock fail")
	}
	return nil
}