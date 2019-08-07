## chat_scheme (FLAG坑王...)
### 1. 需求
* 即时弹幕系统
    - 数百人同时处于一个"直播间"
    - 其中一人发送的弹幕信息，“直播间”内所有人能在客户端看到
    - 数据不会落地，实时推送给前端建立连接的浏览器
* 私密聊天系统
    - 理论上点对点私密聊天应该对数据进行加密？
    - 当其中一人不在线，数据是否需要进行缓存？
    - 讲道理，kafka本身就能缓存指定时间内的数据。。。
    
### 2. 方案设计-即时弹幕系统
##### 2.1 基本技术
* 前端js建立websocket
* 后端iris支持websocket
* kafka提供直播间消息

##### 2.2 存储方案
* roomNo目测可以直接替换为微博mid,类似于直播间或者特定的视频微博都可以用mid作为唯一标签
* 需要开一个针对roomNo为topic的kafka消息队列进行弹幕信息收集和分发
    - 那么问题来了，什么时候开这个消息队列？
    - 什么时候删除这个消息队列？
* 每增加一个观众，也就是一个用户进入了对应的room进行观看，就需要针对上述kafka队列新起一个consumer group消费队列中的消息，模式 newest
    - 首要问题，一个kafka队列能起多少个group去消费？
    - 从kafka读取出来之后，针对观众是否需要展示全量弹幕？
    - 如果是提取部分弹幕，那是不是按概率随机略过一部分
    - 如果量太大，读取之后怎么分批发送到前端？
    - 有一些网站的弹幕（其实就是b站）可以做到跟视频时间强相关，怎么把控这部分时间？
    - 如果发送阻塞了怎么处理？就会导致弹幕和视频内容不能相关

##### 2.3 基本流程
* 进入room--domain/chat/videoview/[roomNo]
    - 返回该roomNo对应的视频信息，同时在前端js渲染中加入roomNo
    - 页面渲染完成后开始建立websocket--GET domain/chat/barrage/websocket/[roomNo]
    - 建立成功后，前端向后端发送消息,event--"connected"
* 后端针对该websocket连接接受到event--“connected”进行一系列操作
    - 每个websocket链接都有一个独特id（此处目前是框架默认的randomString(64)，未保证完全唯一，在分布式环境中更难以应用，待修改为分布式唯一id生成器）
    - 后端响应前端发起的"connected"event,新启一个协程,进行轮训，从该房间的弹幕池中拉取数据
        + 数据首先肯定是在kafka队列中，每个连接用户需要独立的获取弹幕池中的数据，而且。。。需要时间上的配合，这个可以在前端js编写逻辑。但后端还是需要提供某个时间段内的所有弹幕数据。
        + 此处需要注意应该使用拉模式还是推模式。但想想这个跟微博的feed不同，feed如果用户不去拉，那么用拉模式能显著减少不必要的推送。但弹幕系统中之只要在直播间中，就一定会实时推送弹幕信息，所以似乎推模式和拉模式差别不大。
        + 不过考虑前端用户可以自由控制是否关闭弹幕，与其还需要增加交互请求确认是否需要继续推送弹幕信息，还不如改为拉模式。也可以少开一个死循环协程。缺点：多了更多的前端到后端的请求，每次拉取弹幕都得先请求。似乎多了更多的流量。光这点似乎就可以不考虑了。
        + 目前初定这样的处理方式。一个直播间的kafka消息队列,指定n个不同的consumer group 去消费，按照一定的逻辑（可能就是随机）略过百分之xx的弹幕后，将数据存入各个consumer group 的redis 有序集合中，此处集合可以用BarrageInfo的序列化字符串为val,发送弹幕时间为score。
        + 其实也不需要启动n个不同的consumer group去消费一个房间的弹幕，甚至不需要kafka队列，kafka队列还是需要，给弹幕发布一个中间缓冲。但消费group可以只用一个，在一个消费任务里，将弹幕信息写入多个redis set,来扛住高并发的读
        + 当一个用户进入该直播间，长连接被建立，则需要将该长连接，以connetedId作为唯一标识分配给一个上述的consumer group 的redis。响应了onConnected后，每个长连接对应建立一个协程去对应的redis池中拉取弹幕信息
* 前端接收到后端的消息信息后，需要根据消息中的时间进行展示。这块先不想了。

##### 2.4 任务调度
###### 2.4.1 基本需求
* 要求在接收到开启房间的消息后，根据room基本配置开启指定个数的consumer group协程进行指定房间号的弹幕的消费任务
* 当接收到关闭房间的消息后，需要能够即时终止消费

###### 2.4.2 实现方式
* 要求在接收到开启房间的消息后，根据room基本配置开启指定个数的consumer group协程进行指定房间号的弹幕的消费任务
* 当接收到关闭房间的消息后，需要能够即时终止消费

### 3. 数据存储
* 房间配置信息
```sql
CREATE TABLE `chat_room_configure` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `room_id` bigint(18) unsigned NOT NULL COMMENT '房间id',
  `room_name` varchar(64) NOT NULL COMMENT '房间名',
  `room_owner_uid` bigint(18) unsigned NOT NULL COMMENT '房间所有人uid',
  `redis_set_cnt` int(10) unsigned DEFAULT 1 NOT NULL COMMENT '用于前端读取弹幕的redis set个数，同时也是kafka的消费群组个数',
  `status` tinyint(4) unsigned DEFAULT 0 NOT NULL COMMENT '房间状态：0:无效；1:有效',
  `work_status` tinyint(4) unsigned DEFAULT 0 NOT NULL COMMENT '房间开播状态：0:未开播；1:启动中；2:开播中；3:关闭中',
  `create_time` timestamp NOT NULL DEFAULT '1970-01-01 08:00:01',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY (`room_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='直播房间基本配置信息表'
```

### 4. 异常处理
#### 4.1 队列处理服务重新发布
这块其实一直不太了解，发布上线影响线上数据处理的过程
#### 4.2 redis扩容
redis是为了应对弹幕读的部分，当观众比较多的时候，通过增加同一批弹幕的redis set的数量，可以提升弹幕的读速度。除了redis服务本身扩充外，再增加对应的room配置中的rediscnt即可，
#### 4.3 消费队列机器扩充
房间消息队列初始状态多设置几个分区，一个消费群里的消费者数量可以少于分区树。后期扩容的时候，直接增加消费者机器，在消费者机器少于等于分区数前提下，消费者扩充即可增强弹幕分发处理速度
