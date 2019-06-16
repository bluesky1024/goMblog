package relationDbRepo

import (
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/tools/logger"
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
