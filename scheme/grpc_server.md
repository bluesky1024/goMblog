## grpc_server
### 需求
* 将services下多个服务均纳入gprc服务中，彼此之间通过gprc进行服务调用，以便后续进行微服务化
* 多个服务之间怎么解耦合？？？
* 提供服务和调用服务的安全性？？？
* 服务发现和服务治理？？？

### 实现
#### 安全性
* TLS认证（通用配置）
    * 服务端配置证书及私钥
    * 客户端安装证书
* 自定义签名认证 （根据不同服务接口权限定制）


type FollowInfo struct {
    Id         int32
    Uid        int64
    FollowUid  int64
    Status     int8
    IsFriend   int8
    GroupId    int64
    CreateTime time.Time `xorm:"created"`
    UpdateTime time.Time `xorm:"updated"`
}

type FanInfo struct {
    Id         int32
    Uid        int64
    FanUid     int64
    Status     int8
    IsFriend   int8
    CreateTime time.Time `xorm:"created"`
    UpdateTime time.Time `xorm:"updated"`
}

type FollowGroup struct {
    Id int64
    Uid int64
    GroupName  string
    Status     int8
    CreateTime time.Time `xorm:"created"`
    UpdateTime time.Time `xorm:"updated"`
}

type FollowKafkaStruct struct {
    Uid       int64
    FollowUid int64
    Status    int8
}

