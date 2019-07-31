package chatRdRepo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

type BarrageInfoInSet struct {
	Uid       int64
	Message   string
	VideoTime int64
}

//往指定集合中塞入新的弹幕信息
func (r *ChatRbRepository) AppendNewBarrage(uid int64, setInd int, barrageInfo datamodels.ChatBarrageInfo) error {
	key := getBarrageSetInfo(uid, setInd)
	info := BarrageInfoInSet{
		Uid:       barrageInfo.Uid,
		Message:   barrageInfo.Message,
		VideoTime: barrageInfo.VideoTime,
	}
	str, _ := json.Marshal(&info)
	score := barrageInfo.CreateTime.UnixNano() / 1e6
	fmt.Println(key, str, score)
	res, err := r.redisPool.ZAdd(key, redis.Z{
		Score:  float64(score),
		Member: str,
	}).Result()
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	if res == 0 {
		err = errors.New("insert new barrage failed")
		logger.Err(logType, err.Error())
		return err
	}
	return nil
}

//从指定集合中获取timeAfter时间之后的size条弹幕信息
func (r *ChatRbRepository) GetBarragesByUid(uid int64, setInd int, size int, timeAfter time.Time) ([]datamodels.ChatBarrageInfo, error) {
	key := getBarrageSetInfo(uid, setInd)

	startTime := timeAfter.UnixNano() / 1e6
	curTime := time.Now().UnixNano() / 1e6
	//ZRANGEBYSCORE key (min (max LIMIT offset count
	redisRes, err := r.redisPool.ZRangeByScoreWithScores(key, redis.ZRangeBy{
		Min:    "(" + strconv.FormatInt(startTime, 10),
		Max:    "(" + strconv.FormatInt(curTime, 10),
		Offset: 0,
		Count:  int64(size),
	}).Result()
	if err != nil {
		logger.Err(logType, err.Error())
		return nil, err
	}
	res := make([]datamodels.ChatBarrageInfo, len(redisRes))
	for ind, val := range redisRes {
		err = json.Unmarshal([]byte(val.Member.(string)), &res[ind])
		res[ind].CreateTime = time.Unix(int64(val.Score)/1e3, int64(val.Score)%1e3*1e6)
	}

	return res, nil
}

//删除某个房间指定时间之前的弹幕
func (r *ChatRbRepository) RemoveBarrageByUid(uid int64, setInd int, timeBefore time.Time) error {
	key := getBarrageSetInfo(uid, setInd)

	//ZREMRANGEBYSCORE key min max
	_, err := r.redisPool.ZRevRangeByScore(key, redis.ZRangeBy{
		Min: "0",
		Max: strconv.FormatInt(timeBefore.UnixNano()/1e6, 10),
	}).Result()
	if err != nil {
		logger.Err(logType, err.Error())
		return err
	}
	return nil
}
