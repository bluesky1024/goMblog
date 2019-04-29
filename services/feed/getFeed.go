package feedService

import (
	"github.com/bluesky1024/goMblog/tools/idGenerate"
	"time"
)

func (f *feedService) GetFeedByUidAndGroupId(uid int64, groupId int64, page int, pageSize int) (mids []int64, err error) {
	return f.feedRdRepo.GetFeeds(uid, groupId, page, pageSize)
}

func (f *feedService) GetFeedMoreByMid(uid int64, groupId int64, Mid int64, size int) (mids []int64, err error) {
	var timeBefore int64
	if Mid == 0 {
		timeBefore = time.Now().UnixNano()
		timeBefore = timeBefore / 1e6
	} else {
		timeBefore = idGenerate.GetDetailTimeById(Mid)
	}
	return f.feedRdRepo.GetByTimeBefore(uid, groupId, timeBefore, size)
}

func (f *feedService) GetFeedNewerByMid(uid int64, groupId int64, Mid int64, size int) (mids []int64, err error) {
	timeAfter := idGenerate.GetDetailTimeById(Mid)
	return f.feedRdRepo.GetByTimeAfter(uid, groupId, timeAfter, size)
}
