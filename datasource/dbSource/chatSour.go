package datasource

import (
	//"fmt"
	//	"errors"
	//	"log"
	"sync"

	"github.com/bluesky1024/goMblog/config"
	"github.com/go-xorm/xorm"
)

var chatMInstance *xorm.Engine
var chatSInstance *xorm.Engine
var chatLock *sync.Mutex = &sync.Mutex{}

//单例模式获取
func LoadChats(master bool) (*xorm.Engine, error) {
	var err error = nil
	if master {
		if chatMInstance == nil {
			chatLock.Lock()
			defer chatLock.Unlock()
			if chatMInstance == nil {
				chatConfigMap := conf.InitConfig("dbConfig.chat")
				chatMInstance, err = loadDbSource(chatConfigMap["m_host"],
					chatConfigMap["m_port"],
					chatConfigMap["m_user"],
					chatConfigMap["m_password"],
					chatConfigMap["m_dbname"])
				if err != nil {
					return nil, err
				}
			}
		}
		return chatMInstance, err
	} else {
		if chatSInstance == nil {
			chatLock.Lock()
			defer chatLock.Unlock()
			if chatSInstance == nil {
				chatConfigMap := conf.InitConfig("dbConfig.chat")
				chatSInstance, err = loadDbSource(chatConfigMap["s_host"],
					chatConfigMap["s_port"],
					chatConfigMap["s_user"],
					chatConfigMap["s_password"],
					chatConfigMap["s_dbname"])
				if err != nil {
					return nil, err
				}
			}
		}
		return chatSInstance, err
	}
}
