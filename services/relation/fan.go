package relationService

import (
	dm "github.com/bluesky1024/goMblog/datamodels"
)

func (s *relationService) GetFansByUid(uid int64, page int, pageSize int) (fans []dm.FanInfo, cnt int64) {
	fans, cnt = s.repo.SelectMultiFansByUid(uid, page, pageSize)
	totalCnt, err := s.repo.CountFansByUid(uid)
	if err != nil {
		return nil, 0
	}
	return fans, totalCnt
}
