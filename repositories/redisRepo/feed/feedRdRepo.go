package feedRdRepo

import (
	"errors"
	"github.com/bluesky1024/goMblog/tools/idGenerate"
	"github.com/bluesky1024/goMblog/tools/logger"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

var logType = "feedRdRepo"

type FeedRbRepository struct {
	RedisPool *redis.ClusterClient
}

func NewFeedRdRepo(redisPool *redis.ClusterClient) *FeedRbRepository {
	return &FeedRbRepository{
		RedisPool: redisPool,
	}
}

func (f *FeedRbRepository) GetFeeds(uid int64, groupId int64, page int, pageSize int) (mids []int64, err error) {
	if uid == 0 || groupId < 0 {
		return mids, errors.New("not valid param")
	}

	feedKey := getFeedKey(uid, groupId)
	startInd := (page - 1) * pageSize
	endInd := page*pageSize - 1
	res, err := f.RedisPool.ZRevRange(feedKey, int64(startInd), int64(endInd)).Result()
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

	feedKey := getFeedKey(uid, groupId)
	res, err := f.RedisPool.ZRevRange(feedKey, 0, 0).Result()
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

	feedKey := getFeedKey(uid, groupId)
	res, err := f.RedisPool.ZRevRange(feedKey, -1, -1).Result()
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

//timeBefore为ms级时间戳
//获取指定时间之前的feed数据（时间戳按从大到小的顺序）
func (f *FeedRbRepository) GetByTimeBefore(uid int64, groupId int64, timeBefore int64, size int) (mids []int64, err error) {
	if uid == 0 || groupId < 0 {
		return mids, errors.New("invalid uid or groupId")
	}

	feedKey := getFeedKey(uid, groupId)
	//ZREVRANGEBYSCORE key (max (min LIMIT offset count
	res, err := f.RedisPool.ZRevRangeByScore(feedKey, redis.ZRangeBy{
		Max:    "(" + strconv.FormatInt(timeBefore, 10),
		Min:    "(0",
		Offset: 0,
		Count:  int64(size),
	}).Result()
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

	feedKey := getFeedKey(uid, groupId)
	curTime := time.Now().UnixNano()
	curTime = curTime / 1e6
	//ZRANGEBYSCORE key (min (max LIMIT offset count
	res, err := f.RedisPool.ZRangeByScore(feedKey, redis.ZRangeBy{
		Min:    "(" + strconv.FormatInt(timeAfter, 10),
		Max:    "(" + strconv.FormatInt(curTime, 10),
		Offset: 0,
		Count:  int64(size),
	}).Result()
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

	feedKey := getFeedKey(uid, groupId)
	mblogTime := idGenerate.GetDetailTimeById(mid)
	//_, err = conn.Do("zadd", feedKey, mblogTime, mid)
	res, err := f.RedisPool.ZAdd(feedKey, redis.Z{
		Score:  float64(mblogTime),
		Member: mid,
	}).Result()
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	if res == 0 {
		err = errors.New("append fail")
		logger.Err(logType, err.Error())
		return err
	}
	return nil
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

		feedKey := getFeedKey(uid, groupId)
		firstStr := strconv.FormatInt(first, 10)
		lastStr := strconv.FormatInt(last, 10)
		_, err = f.RedisPool.ZRemRangeByLex(feedKey, "["+firstStr, "["+lastStr).Result()
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

	feedKey := getFeedKey(uid, groupId)
	midsInterface := make([]interface{}, len(mids))
	for ind, v := range mids {
		midsInterface[ind] = v
	}
	_, err = f.RedisPool.ZRem(feedKey, midsInterface...).Result()
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	return err
}
