package datasource

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func loadDbSource(host string, port string, user string, password string, db_name string) (*xorm.Engine, error) {
	connUrl := user + ":" + password + "@" + "tcp(" + host + ":" + port + ")/" + db_name + "?charset=utf8"
	orm, err := xorm.NewEngine("mysql", connUrl)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	orm.SetMaxOpenConns(200)

	//调试用，可以在控制台打印生成的sql语句
	//orm.ShowSQL(true)

	return orm, nil
}
