package feedService

import (
	"errors"
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
)

//判定是否需要以推模式更新该用户的微博到其粉丝的feed中
func (f *feedService) checkNeedPullUserMblog(uid int64) (need bool, err error) {
	//判断该用户的属性(若粉丝数太多，feed采用推模式，不用给feed池里逐一添加,目前仅参考粉丝数，应该改为一个独立字段)
	userInfo, err := f.userSrv.GetByUid(uid)
	if err != nil {
		err = errors.New("uid is invalid")
		logger.Err(logType, err.Error())
		return false, err
	}
	if userInfo.FriendsCount > int64(dm.FansCntGate) {
		return false, nil
	}
	return true, nil
}
