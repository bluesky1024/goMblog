package mblogDbRepo

import (
	"sync"
	"time"

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
	for tableName, tempMids := range midMap {
		wg.Add(1)
		go func(table string, mids []int64) {
			defer wg.Done()

			var tempMblogs []dm.MblogInfo
			err := r.sourceS.Table(table).In("mid", mids).Desc("id").Find(&tempMblogs)
			if err != nil {
				logger.Err(logType, err.Error())
				return
			}
			for _, mblog := range tempMblogs {
				mblogs[int64(mblog.Mid)] = mblog
			}
		}(tableName, tempMids)
	}
	wg.Wait()

	////顺序遍历各张表
	//for tableName, tempMids := range midMap {
	//	var tempMblogs []dm.MblogInfo
	//	err := r.sourceS.Table(tableName).In("mid", tempMids).Desc("id").Find(&tempMblogs)
	//	if err != nil {
	//		logger.Err(logType, err.Error())
	//		return
	//	}
	//	for ind, mblog := range tempMblogs {
	//		mblogs[int64(ind)] = mblog
	//	}
	//}
	return mblogs
}

//根据uid获取微博,默认顺序是根据更新时间从大到小排列
func (r *MblogDbRepository) SelectNormalByUid(uid int64, page int, pageSize int, readAbles []int8, startTime int64, endTime int64) (mblogs []dm.MblogInfo, cnt int64) {
	//映射表
	var midsMatch []dm.UidToMblog
	start := (page - 1) * pageSize
	tempSession := r.sourceS.Table(getUidToMblogTableName(uid)).Where("uid = ?", uid).And("status = ?", dm.MblogStatusNormal).In("read_able", readAbles)
	if startTime > 0 {
		tm := time.Unix(startTime, 0)
		tempSession = tempSession.Where("update_time >= ?", tm.Format("02/01/2006 15:04:05"))
	}
	if endTime > 0 {
		tm := time.Unix(endTime, 0)
		tempSession = tempSession.Where("update_time <= ?", tm.Format("02/01/2006 15:04:05"))
	}
	err := tempSession.Limit(pageSize, start).Desc("update_time").Find(&midsMatch)
	if err != nil {
		logger.Err(logType, err.Error())
		return mblogs, 0
	}
	//总数
	var uidToMblog = new(dm.UidToMblog)
	tempSession = r.sourceS.Table(getUidToMblogTableName(uid)).Where("uid = ?", uid).And("status = ?", dm.MblogStatusNormal)
	if startTime > 0 {
		tm := time.Unix(startTime, 0)
		tempSession = tempSession.Where("update_time >= ?", tm.Format("02/01/2006 15:04:05"))
	}
	if endTime > 0 {
		tm := time.Unix(endTime, 0)
		tempSession = tempSession.Where("update_time <= ?", tm.Format("02/01/2006 15:04:05"))
	}
	cnt, err = tempSession.In("read_able", readAbles).Count(uidToMblog)
	if err != nil {
		logger.Err(logType, err.Error())
		return mblogs, 0
	}
	//拼凑mids
	mids := make([]int64, len(midsMatch))
	for ind, midMatch := range midsMatch {
		mids[ind] = midMatch.Mid
	}

	//获取实际微博数据
	mblogMap := r.SelectMultiByMids(mids)

	//按midsMatch顺序拼凑返回的数组
	mblogs = make([]dm.MblogInfo, len(mblogMap))
	ind := 0
	for _, mid := range mids {
		if mblogInfo, ok := mblogMap[mid]; ok {
			mblogs[ind] = mblogInfo
			ind++
		}
	}
	return mblogs, cnt
}
