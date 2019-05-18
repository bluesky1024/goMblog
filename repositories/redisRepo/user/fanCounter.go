package userRdRepo

//用户粉丝数的计数处理逻辑

import (
	"github.com/bluesky1024/goMblog/tools/logger"
	"strconv"
	"time"
)

//增加粉丝和删除粉丝的错误需要记录下来，加入队列进行及时补偿
func (u *UserRbRepository) AddFan(uid int64) (int64, error) {
	fanCnt, err := u.redisPool.Incr(getFanCntKey(uid)).Result()
	if err != nil {
		logger.Err(logType, err.Error())
		return 0, err
	}
	return fanCnt, nil
}

func (u *UserRbRepository) LoseFan(uid int64) (int64, error) {
	fanCnt, err := u.redisPool.Decr(getFanCntKey(uid)).Result()
	if err != nil {
		logger.Err(logType, err.Error())
		return 0, err
	}
	return fanCnt, nil
}

func (u *UserRbRepository) GetFanCnt(uids []int64) (map[int64]int64, error) {
	resMap := make(map[int64]int64)
	fanCntKeys := make([]string, len(uids))
	for i := 0; i < len(uids); i++ {
		fanCntKeys[i] = getFanCntKey(uids[i])
	}
	res, err := u.redisPool.MGet(fanCntKeys...).Result()
	if err != nil {
		logger.Err(logType, err.Error())
		return resMap, err
	}
	for i, v := range res {
		if v == nil {
			resMap[uids[i]] = -1
			continue
		}
		resMap[uids[i]], _ = strconv.ParseInt(v.(string), 10, 64)
	}
	return resMap, nil
}

func (u *UserRbRepository) ExpFanCnt(uid int64, timeMore time.Duration) (bool, error) {
	success, err := u.redisPool.Expire(getFanCntKey(uid), timeMore).Result()
	if err != nil {
		logger.Err(logType, err.Error())
		return false, err
	}
	return success, nil
}
