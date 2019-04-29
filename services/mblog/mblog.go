package mblogService

import (
	"errors"
	dm "github.com/bluesky1024/goMblog/datamodels"
	idGen "github.com/bluesky1024/goMblog/tools/idGenerate"
	"github.com/bluesky1024/goMblog/tools/logger"
)

func (m *mblogService) Create(uid int64, content string, readAble int8, originUid int64, originMid int64) (mblog dm.MblogInfo, err error) {
	if uid <= 0 || content == "" || len(content) > 140 {
		return dm.MblogInfo{}, errors.New("invalid mblog info")
	}

	insertData := dm.MblogInfo{
		Uid:       uid,
		Content:   content,
		OriginMid: originMid,
		OriginUid: originUid,
		ReadAble:  readAble,
		Status:    dm.MblogStatusNormal,
	}
	//生成专属mid
	insertData.Mid, err = idGen.GenMidId()
	if err != nil {
		return dm.MblogInfo{}, err
	}

	affected, err := m.repo.Insert(insertData)
	if err != nil || affected == 0 {
		logger.Err(logType, err.Error())
		return dm.MblogInfo{}, err
	}
	return insertData, nil
}

func (m *mblogService) GetByMid(mid int64) (mblog dm.MblogInfo, found bool) {
	if mid <= 0 {
		return dm.MblogInfo{}, false
	}
	return m.repo.SelectByMid(mid)
}

func (m *mblogService) GetMultiByMids(mids []int64) map[int64]dm.MblogInfo {
	res := m.repo.SelectMultiByMids(mids)
	return res
}

//根据微博时间新旧排序，结合微博可读权限获取指定uid的未删除、未封禁微博列表
func (m *mblogService) GetNormalByUid(uid int64, page int, pageSize int, readAble []int8, startTime int64, endTime int64) (mblogs []dm.MblogInfo, cnt int64) {
	if uid <= 0 {
		return nil, 0
	}
	mblogs, cnt = m.repo.SelectNormalByUid(uid, page, pageSize, readAble, 0, 0)
	return mblogs, cnt
}
