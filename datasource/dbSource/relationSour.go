package datasource

import (
	//"fmt"
	//	"errors"
	//	"log"
	"sync"

	"github.com/bluesky1024/goMblog/config"
	"github.com/go-xorm/xorm"
)

var relationMInstance *xorm.Engine
var relationSInstance *xorm.Engine
var relationLock *sync.Mutex = &sync.Mutex{}

//单例模式获取
func LoadRelation(master bool) (*xorm.Engine, error) {
	var err error = nil
	if master {
		if relationMInstance == nil {
			relationLock.Lock()
			defer relationLock.Unlock()
			if relationMInstance == nil {
				relationConfigMap := conf.InitConfig("dbConfig.relation")
				relationMInstance, err = loadDbSource(relationConfigMap["m_host"],
					relationConfigMap["m_port"],
					relationConfigMap["m_user"],
					relationConfigMap["m_password"],
					relationConfigMap["m_dbname"])
				if err != nil {
					return nil, err
				}
			}
		}
		return relationMInstance, err
	} else {
		if relationSInstance == nil {
			relationLock.Lock()
			defer relationLock.Unlock()
			if relationSInstance == nil {
				relationConfigMap := conf.InitConfig("dbConfig.relation")
				relationSInstance, err = loadDbSource(relationConfigMap["s_host"],
					relationConfigMap["s_port"],
					relationConfigMap["s_user"],
					relationConfigMap["s_password"],
					relationConfigMap["s_dbname"])
				if err != nil {
					return nil, err
				}
			}
		}
		return relationSInstance, err
	}
}
