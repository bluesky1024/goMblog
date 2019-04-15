package mblogDbRepo

import (
	"sync"

	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
	"github.com/go-xorm/xorm"
)

var logType string = "mblogDbRepo"

type MblogDbRepository struct {
	sourceM *xorm.Engine
	sourceS *xorm.Engine
}

func NewMblogRepository(sourceM *xorm.Engine, sourceS *xorm.Engine) *MblogDbRepository {
	return &MblogDbRepository{
		sourceM: sourceM,
		sourceS: sourceS,
	}
}

// 插入新的微博记录，需要同时在映射表和微博信息表插入数据，采用事务方式
func (r *MblogDbRepository) Insert(mblog dm.MblogInfo) (affected int64, err error) {
	session := r.sourceM.NewSession()
	defer session.Clone()

	err = session.Begin()
	if err != nil {
		logger.Err(logType, err.Error())
		return 0, err
	}
	//插入微博信息表
	affected, err = r.sourceM.Table(getMblogTableName(mblog.Mid)).Insert(mblog)
	if err != nil {
		session.Rollback()
		logger.Err(logType, err.Error())
		return 0, err
	}
	//插入用户微博映射表
	uidToMblogData := dm.UidToMblog{
		Uid:      mblog.Uid,
		Mid:      mblog.Mid,
		Status:   mblog.Status,
		ReadAble: mblog.ReadAble,
	}
	_, err = r.sourceM.Table(getUidToMblogTableName(uidToMblogData.Uid)).Insert(uidToMblogData)
	if err != nil {
		session.Rollback()
		logger.Err(logType, err.Error())
		return 0, err
	}

	err = session.Commit()
	if err != nil {
		session.Rollback()
		logger.Err(logType, err.Error())
		return 0, err
	}

	return affected, nil
}

func (r *MblogDbRepository) SelectByMid(mid int64) (mblog dm.MblogInfo, found bool) {
	if mid <= 0 {
		logger.Err(logType, "invalid mid")
		return dm.MblogInfo{}, false
	}
	found, err := r.sourceS.Table(getMblogTableName(mid)).Where("mid = ?", mid).Get(&mblog)
	if err != nil {
		logger.Err(logType, err.Error())
		return dm.MblogInfo{}, false
	}
	return mblog, found
}

func (r *MblogDbRepository) SelectMultiByMids(mids []int64) (mblogs map[int64]dm.MblogInfo) {
	mblogs = make(map[int64]dm.MblogInfo)

	//根据mid获取对应mbloginfo表名并分组查找
	midMap := make(map[string]([]int64))
	for _, mid := range mids {
		midMap[getMblogTableName(mid)] = append(midMap[getMblogTableName(mid)], mid)
	}

	//根据map开多个线程并行进行表查询
	wg := sync.WaitGroup{}
	wg.Add(len(midMap))
	for tableName, tempMids := range midMap {
		go func(table string) {
			defer wg.Done()

			var tempMblogs []dm.MblogInfo
			err := r.sourceS.Table(table).In("mid", tempMids).Desc("id").Find(&tempMblogs)
			if err != nil {
				logger.Err(logType, err.Error())
				return
			}
			for ind, mblog := range tempMblogs {
				mblogs[int64(ind)] = mblog
			}
		}(tableName)
	}
	return mblogs
}

func (r *MblogDbRepository) SelectByUid(uid int64, page int, pageSize int) (mblogs map[int64]dm.MblogInfo) {
	var midsMatch []dm.UidToMblog
	start := (page - 1) * pageSize
	err := r.sourceS.Table(getUidToMblogTableName(uid)).Where("uid = ?", uid).Limit(pageSize, start).Asc("mid").Find(&midsMatch)
	if err != nil {
		logger.Err(logType, err.Error())
		return mblogs
	}

	//拼凑mids
	mids := make([]int64, len(midsMatch))
	for ind, midMatch := range midsMatch {
		mids[ind] = midMatch.Mid
	}
	return r.SelectMultiByMids(mids)
}
