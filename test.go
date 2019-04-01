package main

import (
	"fmt"
	"github.com/bluesky1024/goMblog/datasource/redisSource"
	"github.com/bluesky1024/goMblog/repositories/redisRepo/feed"

	"sync"
	"time"

	idGen "github.com/bluesky1024/goMblog/tools/idGenerate"

	"github.com/bluesky1024/goMblog/config"
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/datasource/dbSource"
	"github.com/bluesky1024/goMblog/repositories/dbRepo/user"
	mblogSrv "github.com/bluesky1024/goMblog/services/mblog"
	//"github.com/go-xorm/xorm"

)

func testConfig() {
	config_map := conf.InitConfig("dbConfig.user")

	fmt.Println(config_map)
}

func testDatasource() {
	fmt.Println(time.Now())
	wg := sync.WaitGroup{}
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			fmt.Println(i)
			test, err1 := datasource.LoadUsers(true)
			fmt.Println("err1", err1)
			user := new(dm.User)
			result, err := test.Where("Id=?", 1).Get(user)
			fmt.Println("result:", result, "err:", err)
			//fmt.Println(user)
		}()
	}
	wg.Wait()
	fmt.Println(time.Now())
}

func testUserRepoUpdate() {
	userSourceM, _ := datasource.LoadUsers(true)
	userSourceS, _ := datasource.LoadUsers(false)
	repo := userDbRepo.NewUserRepository(userSourceM, userSourceS)

	user := dm.User{
		Uid:      123,
		Password: "newpass11",
	}
	a, err := repo.UpdateByUid(user)
	fmt.Println("affected:", a, "err:", err)
}

func testUserRepoInsert() {
	userSourceM, _ := datasource.LoadUsers(true)
	userSourceS, _ := datasource.LoadUsers(false)
	repo := userDbRepo.NewUserRepository(userSourceM, userSourceS)

	user := dm.User{
		Uid:      444,
		Password: "abc",
	}
	a, err := repo.Insert(user)
	fmt.Println("affected:", a, "err:", err)
}

func testGenMid() {
	fmt.Println(time.Now())
	var wg sync.WaitGroup
	wg.Add(2000)
	idGen.InitMidPool(10)
	for i := 0; i < 1000; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				idGen.GenMidId()
			}
			wg.Done()
		}()
	}

	idGen.InitUidPool(10)
	for i := 0; i < 1000; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				idGen.GenUidId()
			}
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println(time.Now(), "end")
}

func testPasswordHashed() {
	userSourceM, _ := datasource.LoadUsers(true)
	userSourceS, _ := datasource.LoadUsers(false)
	repo := userDbRepo.NewUserRepository(userSourceM, userSourceS)
	user := &dm.User{}
	user.NickName = "blue"
	found := repo.Select(user)
	//user, found := repo.SelectByNickname("bluesky1024")
	fmt.Println("user", user, "found:", found)
}

func testMblogOrm(){
	//mblogSourceM, _ := datasource.LoadMblogSour(true)
	//mblogSourceS, _ := datasource.LoadMblogSour(false)

}

func testGetTimeInMid(){
	idGen.InitMidPool(1)
	mid,_ :=idGen.GenMidId()
	fmt.Printf("mid: %b\n",mid)

	var baseMblogTableName string = "mblog_info"

	//时间戳在mid/uid中位置
	var timestampPos int64 = 0x3fffffffffe00000

	mblogTime := (mid&timestampPos)>>21

	fmt.Printf("mblogTime: %b \t %d\n",mblogTime,mblogTime)

	realTime := mblogTime/1e3 + idGen.Epoch/1e3

	fmt.Printf("mblogTime: %b \t %d\n",realTime,realTime)

	//Format("2006-01-02 15:04:05")
	timeFormat := time.Unix(realTime,0).Format("200601")
	fmt.Println(baseMblogTableName+"_"+timeFormat)
}

func testSendMblog(){
	//idGen.InitUidPool(1)
	idGen.InitMidPool(1)
	//uid,_ :=idGen.GenUidId()
	uid := 160846466519040

	srv := mblogSrv.NewMblogServicer()
	mblog,err := srv.Create(int64(uid),"123",1,0,0)
	fmt.Println("mblog:",mblog,"err:",err)
}

func testRedisRepo(){
	redisPoolM,_ := redisSource.LoadRedisSource("127.0.0.1","6379")
	redisPoolS,_ := redisSource.LoadRedisSource("127.0.0.1","6379")
	feedRepo := feedRdRepo.NewFeedRdRepo(redisPoolM,redisPoolS)

	res := feedRepo.RemoveMids(1,1,[]int64{222,444})
	fmt.Println(res)

	//feedRepo.DelFeed(1,1)

	//err := feedRepo.AppendNewMid(1,1,4612233247743545344)
	//fmt.Println(err)
	//
	//mids := feedRepo.GetFeeds(1,1,1,3)
	//fmt.Println(mids)
	//
	//mid,err := feedRepo.GetFirstFeed(1,1)
	//fmt.Println(mid,err)
	//
	//mid,err = feedRepo.GetLastFeed(1,1)
	//fmt.Println(mid,err)
}

func main() {
	testRedisRepo()
}