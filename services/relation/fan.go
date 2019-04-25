package relationService

import (
	dm "github.com/bluesky1024/goMblog/datamodels"
)

func (s *relationService) GetFansByUid(uid int64, page int, pageSize int) (fans []dm.FanInfo, cnt int64) {
	return s.repo.SelectMultiFansByUid(uid, page, pageSize)
}
