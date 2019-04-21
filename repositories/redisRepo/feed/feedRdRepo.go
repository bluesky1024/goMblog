package feedRdRepo

import (
	"errors"
	"fmt"
	"github.com/bluesky1024/goMblog/tools/idGenerate"
	"github.com/bluesky1024/goMblog/tools/logger"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

var logType string = "feedRdRepo"

type FeedRbRepository struct {
	RedisPoolM *redis.Pool
	RedisPoolS *redis.Pool
}

func NewFeedRdRepo(redisPoolM *redis.Pool, redisPoolS *redis.Pool) *FeedRbRepository {
	return &FeedRbRepository{
		RedisPoolM: redisPoolM,
		RedisPoolS: redisPoolS,
	}
}

func (f *FeedRbRepository) TestRedis() {
	conn := f.RedisPoolM.Get()
	defer conn.Close()

	_, err := conn.Do("SET", "my_go_key", "test_one")
	if err != nil {
		fmt.Println("redis set failed:", err)
	}

	username, err := redis.String(conn.Do("GET", "my_go_key"))
	if err != nil {
		fmt.Println("redis get failed:", err)
	} else {
		fmt.Printf("Get mykey: %v \n", username)
	}
}

func (f *FeedRbRepository) GetFeeds(uid int64, groupId int64, page int, pageSize int) (mids []int64, err error) {
	if uid == 0 || groupId < 0 {
		return mids, errors.New("not valid param")
	}
	conn := f.RedisPoolS.Get()
	defer conn.Close()

	feedKey := getFeedKey(uid, groupId)
	startInd := (page - 1) * pageSize
	endInd := page*pageSize - 1
	res, err := redis.Strings(conn.Do("zrevrange", feedKey, startInd, endInd))
	if err != nil {
		logger.Err(logType, err.Error())
		return mids, err
	}

	mids = make([]int64, len(res))
	for ind, val := range res {
		mids[ind], _ = strconv.ParseInt(val, 10, 64)
	}
	return mids, nil
}

func (f *FeedRbRepository) GetFirstMid(uid int64, groupId int64) (mid int64, err error) {
	if uid == 0 || groupId < 0 {
		return mid, errors.New("invalid uid or groupId")
	}
	conn := f.RedisPoolS.Get()
	defer conn.Close()

	feedKey := getFeedKey(uid, groupId)
	res, err := redis.Strings(conn.Do("zrevrange", feedKey, 0, 0))
	if err != nil {
		logger.Err(logType, err.Error())
		return mid, err
	}
	if len(res) == 0 {
		return mid, errors.New("not found key")
	}
	mid, _ = strconv.ParseInt(res[0], 10, 64)
	return mid, err
}

func (f *FeedRbRepository) GetLastMid(uid int64, groupId int64) (mid int64, err error) {
	if uid == 0 || groupId < 0 {
		return mid, errors.New("invalid uid or groupId")
	}
	conn := f.RedisPoolS.Get()
	defer conn.Close()

	feedKey := getFeedKey(uid, groupId)
	res, err := redis.Strings(conn.Do("zrange", feedKey, -1, -1))
	if err != nil {
		logger.Err(logType, err.Error())
		return mid, err
	}
	if len(res) == 0 {
		return mid, errors.New("not found key")
	}
	mid, _ = strconv.ParseInt(res[0], 10, 64)
	return mid, err
}

//该函数注释保留作为如何获取kv对的示例
//timeBefore为ms级时间戳
////获取指定时间之前(<timeBefore)的feed数据（时间戳按从大到小的顺序）
//func (f *FeedRbRepository) GetByTimeBefore(uid int64, groupId int64, timeBefore int64, size int) (res map[string]int64, err error) {
//	if uid == 0 || groupId < 0 {
//		return res, errors.New("invalid uid or groupId")
//	}
//	conn := f.RedisPoolS.Get()
//	defer conn.Close()
//
//	feedKey := getFeedKey(uid, groupId)
//	//ZREVRANGEBYSCORE key (max (min LIMIT offset count
//	values, err := redis.Values(conn.Do("zrevrangebyscore", feedKey, "("+strconv.FormatInt(timeBefore, 10), "(0", "withscores", "limit", 0, size))
//	fmt.Println("res", res)
//	if err != nil {
//		logger.Err(logType, err.Error())
//		return res, err
//	}
//	type tempMid struct {
//		mid  int64
//		time int64
//	}
//	pairs := make([]tempMid, len(values)/2)
//	for i := range pairs {
//		values, err = redis.Scan(values, &pairs[i].mid, &pairs[i].time)
//		if err != nil {
//			logger.Err(logType, err.Error())
//		}
//	}
//	return res, nil
//}

//timeBefore为ms级时间戳
//获取指定时间之前的feed数据（时间戳按从大到小的顺序）
func (f *FeedRbRepository) GetByTimeBefore(uid int64, groupId int64, timeBefore int64, size int) (mids []int64, err error) {
	if uid == 0 || groupId < 0 {
		return mids, errors.New("invalid uid or groupId")
	}
	conn := f.RedisPoolS.Get()
	defer conn.Close()

	feedKey := getFeedKey(uid, groupId)
	//ZREVRANGEBYSCORE key (max (min LIMIT offset count
	res, err := redis.Strings(conn.Do("zrevrangebyscore", feedKey, "("+strconv.FormatInt(timeBefore, 10), "(0", "limit", 0, size))
	if err != nil {
		logger.Err(logType, err.Error())
		return mids, err
	}
	mids = make([]int64, len(res))
	for ind, val := range res {
		mids[ind], _ = strconv.ParseInt(val, 10, 64)
	}
	return mids, nil
}

//timeAfter为ms级时间戳
//获取指定时间之后的feed数据（时间戳按从小到大的顺序）
func (f *FeedRbRepository) GetByTimeAfter(uid int64, groupId int64, timeAfter int64, size int) (mids []int64, err error) {
	if uid == 0 || groupId < 0 {
		return mids, errors.New("invalid uid or groupId")
	}
	conn := f.RedisPoolS.Get()
	defer conn.Close()

	feedKey := getFeedKey(uid, groupId)
	curTime := time.Now().UnixNano()
	curTime = curTime / 1e6
	//ZRANGEBYSCORE key (min (max LIMIT offset count
	res, err := redis.Strings(conn.Do("zrangebyscore", feedKey, "("+strconv.FormatInt(timeAfter, 10), "("+strconv.FormatInt(curTime, 10), "limit", 0, size))
	if err != nil {
		logger.Err(logType, err.Error())
		return mids, err
	}
	mids = make([]int64, len(res))
	for ind, val := range res {
		mids[ind], _ = strconv.ParseInt(val, 10, 64)
	}
	return mids, nil
}

//在feed池中新增数据
func (f *FeedRbRepository) AppendNewMid(uid int64, groupId int64, mid int64) (err error) {
	if uid == 0 || groupId < 0 {
		return errors.New("invalid uid or groupId")
	}
	conn := f.RedisPoolM.Get()
	defer conn.Close()

	feedKey := getFeedKey(uid, groupId)
	mblogTime := idGenerate.GetDetailTimeById(mid)
	_, err = conn.Do("zadd", feedKey, mblogTime, mid)
	if err != nil {
		logger.Err(logType, err.Error())
	}
	return err
}

//feed被删除之后即无法再获取，故删除feed的操作可以异步执行
func (f *FeedRbRepository) DelFeed(uid int64, groupId int64) {
	go func() {
		first, err := f.GetFirstMid(uid, groupId)
		if err != nil {
			logger.Err(logType, err.Error())
			return
		}
		last, err := f.GetLastMid(uid, groupId)
		if err != nil {
			logger.Err(logType, err.Error())
			return
		}

		conn := f.RedisPoolM.Get()
		defer conn.Close()

		feedKey := getFeedKey(uid, groupId)
		firstStr := strconv.FormatInt(first, 10)
		lastStr := strconv.FormatInt(last, 10)
		_, err = conn.Do("zremrangebylex", feedKey, "["+firstStr, "["+lastStr)
		if err != nil {
			logger.Err(logType, err.Error())
			return
		}
	}()
}

//从feed池中批量删除数据
func (f *FeedRbRepository) RemoveMids(uid int64, groupId int64, mids []int64) (err error) {
	if len(mids) == 0 {
		return errors.New("empty mids del")
	}
	conn := f.RedisPoolM.Get()
	defer conn.Close()

	feedKey := getFeedKey(uid, groupId)
	_, err = conn.Do("zrem", redis.Args{}.Add(feedKey).AddFlat(mids)...)
	if err != nil {
		logger.Err(logType, err.Error())
		return
	}
	return err
}
