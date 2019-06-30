package relationService

import (
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
)

func (s *relationService) GetFanCntByUids(uids []int64) map[int64]int64 {
	res, err := s.rdRepo.GetFanCnt(uids)
	if err != nil {
		logger.Err(logType, err.Error())
	}
	return res
}

func (s *relationService) GetFansByUid(uid int64, page int, pageSize int) (fans []dm.FanInfo, totalCnt int64) {
	fans, _ = s.repo.SelectMultiFansByUid(uid, page, pageSize)
	cntMap, err := s.rdRepo.GetFanCnt([]int64{uid})
	if len(cntMap) == 0 || err != nil {
		return fans, 0
	}
	return fans, cntMap[uid]
}
