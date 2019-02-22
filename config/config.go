/*配置文件加载器，后续优化*/
package conf

import (
	"bufio"
	"io"
	"os"
	"strings"
)

var configs_map map[string](map[string]string)

func init() {
	configs_map = make(map[string](map[string]string))
}

func InitConfig(path string) map[string]string {
	if _, ok := configs_map[path]; ok {
		return configs_map[path]
	}

	//拆分path(例：db_config.user_db)
	path_arr := strings.Split(path, ".")
	path = os.Getenv("GOPATH")
	path = path + "/src/github.com/bluesky1024/goMblog/config"
	for _, temp := range path_arr {
		path = path + "/" + temp
	}
	path = path + ".ini"

	//初始化
	myMap := make(map[string]string)

	//打开文件指定目录，返回一个文件f和错误信息
	f, err := os.Open(path)
	defer f.Close()

	//异常处理 以及确保函数结尾关闭文件流
	if err != nil {
		panic(err)
	}

	//创建一个输出流向该文件的缓冲流*Reader
	r := bufio.NewReader(f)
	for {
		//读取，返回[]byte 单行切片给b
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		//去除单行属性两端的空格
		s := strings.TrimSpace(string(b))
		//fmt.Println(s)

		//判断等号=在该行的位置
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}
		//取得等号左边的key值，判断是否为空
		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}

		//取得等号右边的value值，判断是否为空
		value := strings.TrimSpace(s[index+1:])
		if len(value) == 0 {
			continue
		}
		//这样就成功吧配置文件里的属性key=value对，成功载入到内存中c对象里
		myMap[key] = value
	}
	configs_map[path] = myMap
	return myMap
}
