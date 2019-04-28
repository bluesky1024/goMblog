package main

import (
	"encoding/json"
	"fmt"
	"github.com/bluesky1024/goMblog/datasource/redisSource"
	"github.com/bluesky1024/goMblog/repositories/redisRepo/feed"
	"reflect"

	"sync"
	"time"

	idGen "github.com/bluesky1024/goMblog/tools/idGenerate"

	"github.com/bluesky1024/goMblog/config"
	dm "github.com/bluesky1024/goMblog/datamodels"
	"github.com/bluesky1024/goMblog/datasource/dbSource"
	"github.com/bluesky1024/goMblog/repositories/dbRepo/user"
	mblogSrv "github.com/bluesky1024/goMblog/services/mblog"
	"github.com/bluesky1024/goMblog/services/userGrpc"
	//"github.com/go-xorm/xorm"
	"math"
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

func testMblogOrm() {
	//mblogSourceM, _ := datasource.LoadMblogSour(true)
	//mblogSourceS, _ := datasource.LoadMblogSour(false)

}

func testGetTimeInMid() {
	idGen.InitMidPool(1)
	mid, _ := idGen.GenMidId()
	fmt.Printf("mid: %b\n", mid)

	var baseMblogTableName string = "mblog_info"

	//时间戳在mid/uid中位置
	var timestampPos int64 = 0x3fffffffffe00000

	mblogTime := (mid & timestampPos) >> 21

	fmt.Printf("mblogTime: %b \t %d\n", mblogTime, mblogTime)

	realTime := mblogTime/1e3 + idGen.Epoch/1e3

	fmt.Printf("mblogTime: %b \t %d\n", realTime, realTime)

	//Format("2006-01-02 15:04:05")
	timeFormat := time.Unix(realTime, 0).Format("200601")
	fmt.Println(baseMblogTableName + "_" + timeFormat)
}

func testSendMblog() {
	//idGen.InitUidPool(1)
	idGen.InitMidPool(1)
	//uid,_ :=idGen.GenUidId()
	uid := 160846466519040

	srv, _ := mblogSrv.NewMblogServicer()
	mblog, err := srv.Create(int64(uid), "123", 1, 0, 0)
	fmt.Println("mblog:", mblog, "err:", err)
}

func testRedisRepo() {
	redisPoolM, _ := redisSource.LoadRedisSource("127.0.0.1", "6379")
	redisPoolS, _ := redisSource.LoadRedisSource("127.0.0.1", "6379")
	feedRepo := feedRdRepo.NewFeedRdRepo(redisPoolM, redisPoolS)

	res := feedRepo.RemoveMids(1, 1, []int64{222, 444})
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

func testUserGrpcServer() {
	userServ := userGrpc.NewUserGrpcServicer()
	if userServ == nil {
		fmt.Println("new server error")
		return
	}
	fmt.Println("gen server success")
	//userInfo,err := userServ.Create("testGrpc","123","188","a@a.com")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(userInfo)

	uids := make([]int64, 2)
	uids[0] = 160846466519040
	uids[1] = 1055425730449408
	users, err := userServ.GetMultiByUids(uids)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(users)
}

func testJsonEncodeDecode() {
	//realMsg := dm.FollowMsg{
	//	Uid:       1,
	//	FollowUid: 2,
	//	Status:    1,
	//}

	type InnerStruct struct {
		Id   int64
		Name string
		Test string
	}

	type tempStruct struct {
		Uid      int64
		InStruct InnerStruct
		A        int32
	}

	realMsg := tempStruct{
		Uid: 1,
		A:   2,
		InStruct: InnerStruct{
			Id:   123,
			Name: "ABC",
			Test: "EFD",
		},
	}

	msg := dm.KafkaMsg{
		MsgId:     111,
		Topic:     "abc",
		Partition: 1,
		Data:      realMsg,
	}

	jsonData, _ := json.Marshal(&msg)

	msgRecover := new(dm.KafkaMsg)
	err := json.Unmarshal(jsonData, msgRecover)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("%s \n", jsonData)
	fmt.Printf("%s \n", msgRecover.Data)

	//realMsgRecover := msgRecover.Data.(dm.FollowMsg)

	//遍历解析msgRecover
	//for ind, v := range msgRecover.Data.(map[string]interface{}) {
	//	fmt.Println(ind, reflect.TypeOf(v))
	//}
	//SearchDeepMap(msgRecover.Data.(map[string]interface{}))
	//fmt.Println(msgRecover.Data)
}

func SearchDeepMap(data map[string]interface{}) {
	for ind, v := range data {
		if reflect.TypeOf(v).Name() == "map[string]interface {}" {
			SearchDeepMap(v.(map[string]interface{}))
			continue
		}
		fmt.Println(ind, reflect.TypeOf(v).Name())
	}
}

func testMblog() {
	mblogserv, _ := mblogSrv.NewMblogServicer()
	readAble := []int8{1, 2}
	mblogs, _ := mblogserv.GetNormalByUid(160846466519040, readAble, 1, 10)
	fmt.Println(mblogs)
}

type PageParams struct {
	CurPage  int64
	Cnt      int64
	BaseUrl  string
	PageList []int64
	First    bool
	Last     bool
	Forward  bool
	Back     bool
}

func GenPageView(curPage int64, pageSize int64, cnt int64, baseUrl string) (pageParam PageParams) {
	var startListInd int64
	var endListInd int64
	pageCnt := int64(math.Ceil(float64(cnt) / float64(pageSize)))
	pageParam.Cnt = pageCnt
	pageParam.BaseUrl = baseUrl
	pageParam.CurPage = curPage
	switch {
	case curPage > 4:
		startListInd = curPage - 2
		pageParam.First = true
		pageParam.Back = true
		break
	case curPage == 4:
		startListInd = curPage - 2
		pageParam.First = true
		pageParam.Back = false
		break
	default:
		startListInd = 1
		pageParam.First = false
		pageParam.Back = false
		break
	}
	switch {
	case curPage+2 < pageCnt-1:
		endListInd = curPage + 2
		pageParam.Last = true
		pageParam.Forward = true
		break
	case curPage+2 == pageCnt-1:
		endListInd = curPage + 2
		pageParam.Last = true
		pageParam.Forward = false
		break
	default:
		endListInd = pageCnt
		pageParam.Last = false
		pageParam.Forward = false
		break
	}
	pageParam.PageList = make([]int64, endListInd-startListInd+1)
	for ind, _ := range pageParam.PageList {
		pageParam.PageList[ind] = startListInd + int64(ind)
	}
	return pageParam
}

func main() {
	pageparams := GenPageView(1, 3, 6, "a/b/c")
	fmt.Println(pageparams)
}
