package datasource

import (
	//"fmt"
	//	"errors"
	//	"log"
	"sync"

	"github.com/bluesky1024/goMblog/config"
	"github.com/go-xorm/xorm"
)

var userMInstance *xorm.Engine
var userSInstance *xorm.Engine
var userLock *sync.Mutex = &sync.Mutex{}

//单例模式获取
func LoadUsers(master bool) (*xorm.Engine, error) {
	var err error = nil
	if master {
		if userMInstance == nil {
			userLock.Lock()
			defer userLock.Unlock()
			if userMInstance == nil {
				userConfigMap := conf.InitConfig("dbConfig.user")
				userMInstance, err = loadDbSource(userConfigMap["m_host"],
					userConfigMap["m_port"],
					userConfigMap["m_user"],
					userConfigMap["m_password"],
					userConfigMap["m_dbname"])
				if err != nil {
					return nil, err
				}
			}
		}
		return userMInstance, err
	} else {
		if userSInstance == nil {
			userLock.Lock()
			defer userLock.Unlock()
			if userSInstance == nil {
				userConfigMap := conf.InitConfig("dbConfig.user")
				userSInstance, err = loadDbSource(userConfigMap["s_host"],
					userConfigMap["s_port"],
					userConfigMap["s_user"],
					userConfigMap["s_password"],
					userConfigMap["s_dbname"])
				if err != nil {
					return nil, err
				}
			}
		}
		return userSInstance, err
	}
}
