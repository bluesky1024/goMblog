package feedRdRepo

import (
	"errors"
	"fmt"
	"github.com/bluesky1024/goMblog/tools/idGenerate"
	"github.com/bluesky1024/goMblog/tools/logger"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

var logType string = "feedRdRepo"

type FeedRbRepository struct {
	RedisPoolM *redis.Pool
	RedisPoolS *redis.Pool
}

func NewFeedRdRepo(redisPoolM *redis.Pool,redisPoolS *redis.Pool) *FeedRbRepository {
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

func (f *FeedRbRepository) GetFeeds(uid int64,groupId int64,page int,pageSize int) (mids []int64) {
	if uid == 0 || groupId < 0 {
		return mids
	}
	conn := f.RedisPoolS.Get()
	defer conn.Close()

	feedKey := getFeedKey(uid,groupId)
	//feedKey = "ss"
	startInd := (page-1)*pageSize
	endInd := page*pageSize-1
	res,err := redis.Strings(conn.Do("zrange",feedKey,startInd,endInd))
	if err != nil {
		logger.Err(logType,err.Error())
		return mids
	}

	mids = make([]int64,len(res))
	for ind,val := range res{
		mids[ind],_ = strconv.ParseInt(val,10,64)
	}
	return mids
}

func (f *FeedRbRepository) GetFirstFeed(uid int64,groupId int64) (mid int64,err error) {
	if uid == 0 || groupId < 0 {
		return mid,errors.New("invalid uid or groupId")
	}
	conn := f.RedisPoolS.Get()
	defer conn.Close()

	feedKey := getFeedKey(uid,groupId)
	res,err := redis.Strings(conn.Do("zrange",feedKey,0,0))
	if err != nil {
		logger.Err(logType,err.Error())
		return mid,err
	}
	if len(res) == 0{
		return mid,errors.New("not found key")
	}
	mid,_ = strconv.ParseInt(res[0],10,64)
	return mid,err
}

func (f *FeedRbRepository) GetLastFeed(uid int64,groupId int64) (mid int64,err error) {
	if uid == 0 || groupId < 0 {
		return mid,errors.New("invalid uid or groupId")
	}
	conn := f.RedisPoolS.Get()
	defer conn.Close()

	feedKey := getFeedKey(uid,groupId)
	res,err := redis.Strings(conn.Do("zrange",feedKey,-1,-1))
	if err != nil {
		logger.Err(logType,err.Error())
		return mid,err
	}
	if len(res) == 0{
		return mid,errors.New("not found key")
	}
	mid,_ = strconv.ParseInt(res[0],10,64)
	return mid,err
}

func (f *FeedRbRepository) AppendNewMid(uid int64,groupId int64,mid int64) (err error) {
	if uid == 0 || groupId < 0 {
		return errors.New("invalid uid or groupId")
	}
	conn := f.RedisPoolM.Get()
	defer conn.Close()

	feedKey := getFeedKey(uid,groupId)
	mblogTime := idGenerate.GetTimeInfoById(mid)
	_,err = conn.Do("zadd",feedKey,mblogTime.Unix(),mid)
	if err != nil {
		logger.Err(logType,err.Error())
	}
	return err
}

//feed被删除之后即无法再获取，故删除feed的操作可以异步执行
func (f *FeedRbRepository) DelFeed(uid int64,groupId int64) {
	go func() {
		first,err := f.GetFirstFeed(uid,groupId)
		if err != nil{
			logger.Err(logType,err.Error())
			return
		}
		last,err := f.GetLastFeed(uid,groupId)
		if err != nil{
			logger.Err(logType,err.Error())
			return
		}

		conn := f.RedisPoolM.Get()
		defer conn.Close()

		feedKey := getFeedKey(uid,groupId)
		firstStr := strconv.FormatInt(first,10)
		lastStr := strconv.FormatInt(last,10)
		res,err := conn.Do("zremrangebylex",feedKey,"["+firstStr,"["+lastStr)
		if err != nil{
			logger.Err(logType,err.Error())
			return
		}
		fmt.Println(res)
	}()
}

func (f *FeedRbRepository) RemoveMids(uid int64,groupId int64,mids []int64) (err error) {
	if len(mids) == 0 {
		return errors.New("empty mids del")
	}
	conn := f.RedisPoolM.Get()
	defer conn.Close()

	feedKey := getFeedKey(uid,groupId)
	_,err = conn.Do("zrem",redis.Args{}.Add(feedKey).AddFlat(mids)...)
	if err != nil{
		logger.Err(logType,err.Error())
		return
	}
	return err
}