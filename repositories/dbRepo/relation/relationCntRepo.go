package relationDbRepo

import (
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
	"sync"
)

func (r *RelationDbRepository) UpdateFanCntByUid(uid int64, incr int) bool {
	relationCnt := dm.RelationCnt{}
	infectRow, err := r.sourceM.Table(getRelationCntTableName(uid)).Where("uid = ?", uid).Incr("fan_cnt", incr).Update(&relationCnt)
	if err != nil {
		logger.Err(logType, err.Error())
		return false
	}
	//关系计数表未插入该uid,进行插入操作
	if infectRow == 0 {
		sql := "INSERT INTO " + getRelationCntTableName(uid) + " (uid,fan_cnt,follow_cnt) VALUES (?,?,0) ON DUPLICATE KEY UPDATE fan_cnt=fan_cnt+?"
		_, err := r.sourceM.Exec(sql, uid, incr, incr)
		if err != nil {
			logger.Err(logType, err.Error())
			return false
		}
	}
	return true
}

func (r *RelationDbRepository) UpdateFollowCntByUid(uid int64, incr int) bool {
	relationCnt := dm.RelationCnt{}
	infectRow, err := r.sourceM.Table(getRelationCntTableName(uid)).Where("uid = ?", uid).Incr("follow_cnt", incr).Update(&relationCnt)
	if err != nil {
		logger.Err(logType, err.Error())
		return false
	}
	//关系计数表未插入该uid,进行插入操作
	if infectRow == 0 {
		sql := "INSERT INTO " + getRelationCntTableName(uid) + " (uid,fan_cnt,follow_cnt) VALUES (?,0,?)  ON DUPLICATE KEY UPDATE follow_cnt=follow_cnt+?"
		_, err := r.sourceM.Exec(sql, uid, incr, incr)
		if err != nil {
			logger.Err(logType, err.Error())
			return false
		}
	}
	return true
}

func (r *RelationDbRepository) getRelationCntByUidsFromOneTable(uids []int64, tableName string) map[int64]dm.RelationCnt {
	res := make(map[int64]dm.RelationCnt)
	var relationCnts []dm.RelationCnt
	err := r.sourceS.Table(tableName).In("uid", uids).In("uid", uids).Find(&relationCnts)
	if err != nil {
		logger.Err(logType, err.Error())
		return res
	}
	for _, relationCnt := range relationCnts {
		res[relationCnt.Uid] = relationCnt
	}
	return res
}

func (r *RelationDbRepository) GetRelationCntByUids(uids []int64) map[int64]dm.RelationCnt {
	res := make(map[int64]dm.RelationCnt)

	//根据uid获取对应relationCnt表名并分组查找
	uidMap := make(map[string]([]int64))
	for _, uid := range uids {
		uidMap[getRelationCntTableName(uid)] = append(uidMap[getRelationCntTableName(uid)], uid)
	}

	//根据map开多个线程并行进行表查询
	wg := sync.WaitGroup{}
	for tableName, tempUids := range uidMap {
		wg.Add(1)
		go func(table string, mids []int64) {
			defer wg.Done()

			tempCntRes := r.getRelationCntByUidsFromOneTable(tempUids, tableName)

			for uid, cntInfo := range tempCntRes {
				res[uid] = cntInfo
			}
		}(tableName, tempUids)
	}
	wg.Wait()

	return res
}
