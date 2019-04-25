package datasource

import (
	//"fmt"
	//	"errors"
	//	"log"
	"sync"

	"github.com/bluesky1024/goMblog/config"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

var mblogMInstance *xorm.Engine
var mblogSInstance *xorm.Engine
var mblogLock *sync.Mutex = &sync.Mutex{}

//单例模式获取
func LoadMblogSour(master bool) (*xorm.Engine, error) {
	var err error = nil
	if master {
		if mblogMInstance == nil {
			userLock.Lock()
			defer userLock.Unlock()
			if mblogMInstance == nil {
				mblogConfigMap := conf.InitConfig("dbConfig.mblog")
				mblogMInstance, err = loadDbSource(mblogConfigMap["m_host"],
					mblogConfigMap["m_port"],
					mblogConfigMap["m_user"],
					mblogConfigMap["m_password"],
					mblogConfigMap["m_dbname"])
				if err != nil {
					return nil, err
				}
			}
			tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, "201903_")
			mblogMInstance.SetTableMapper(tbMapper)
		}
		return mblogMInstance, err
	} else {
		if mblogSInstance == nil {
			userLock.Lock()
			defer userLock.Unlock()
			if mblogSInstance == nil {
				mblogConfigMap := conf.InitConfig("dbConfig.user")
				mblogSInstance, err = loadDbSource(mblogConfigMap["s_host"],
					mblogConfigMap["s_port"],
					mblogConfigMap["s_user"],
					mblogConfigMap["s_password"],
					mblogConfigMap["s_dbname"])
				if err != nil {
					return nil, err
				}
			}
		}
		return mblogSInstance, err
	}
}
