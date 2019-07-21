## chat_scheme (FLAG坑王...)
### 需求
* 即时弹幕系统
    - 数百人同时处于一个"直播间"
    - 其中一人发送的弹幕信息，“直播间”内所有人能在客户端看到
    - 数据不会落地，实时推送给前端建立连接的浏览器
* 私密聊天系统
    - 理论上点对点私密聊天应该对数据进行加密？
    - 当其中一人不在线，数据是否需要进行缓存？
    - 讲道理，kafka本身就能缓存指定时间内的数据。。。
    
#### 方案设计-即时弹幕系统
##### 基本技术
* 前端js建立websocket
* 后端iris支持websocket
* kafka提供直播间消息

##### 存储方案
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

##### 基本流程
* 
* 进入room--domain/chat/videoview/[roomNo]
    - 返回该roomNo对应的视频信息，同时在前端js渲染中加入roomNo
    - 页面渲染完成后开始建立websocket--GET domain/chat/barrage/websocket/[roomNo]
    - 建立成功后，前端向后端发送消息,event--"connected"
* 后端针对该websocket连接接受到event--“connected”进行一系列操作
    - 每个websocket链接都有一个独特id（此处目前是框架默认的randomString(64)，未保证完全唯一，在分布式环境中更难以应用，待修改为分布式唯一id生成器）
    - 后端响应前端发起的"connected"event,新启一个协程,进行轮训，从该房间的弹幕池中拉取数据
        + 数据首先肯定是在kafka队列中，每个连接用户需要独立的获取弹幕池中的数据，而且。。。需要时间上的配合，这个可以在前端js编写逻辑。但后端还是需要提供某个时间段内的所有弹幕数据。
        + 此处需要注意应该使用拉模式还是推模式。但想想这个跟微博的feed不同，feed如果用户不去拉，那么用拉模式能显著减少不必要的推送。但弹幕系统中之只要在直播间中，就一定会实时推送弹幕信息，所以似乎推模式和拉模式差别不大。
        + 目前初定这样的处理方式。一个直播间的kafka消息队列,指定n个不同的consumer group 去消费，按照一定的逻辑（可能就是随机）略过百分之xx的弹幕后，将数据存入各个consumer group 的redis 有序集合中，此处集合可以用BarrageInfo的序列化字符串为val,发送弹幕时间为score。
        + 当一个用户进入该直播间，长连接被建立，则需要将该长连接，以connetedId作为唯一标识分配给一个上述的consumer group 的redis。响应了onConnected后，每个长连接对应建立一个协程去对应的redis池中拉取弹幕信息