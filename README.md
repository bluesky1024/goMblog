## goMblog
### 微博用户操作
FUNCTION|PATH|METHOD|NEED_LOGIN
---|:--:|:--:|:--:
注册页面|/user/register|GET|NO
注册请求|/user/register|POST|NO
登录页面|/user/login|GET|NO
登录请求|/user/login|POST|NO
下线请求|/user/logout|POST|NO

### 个人主页
FUNCTION|PATH|METHOD|NEED_LOGIN
---|:--:|:--:|:--:
指定用户主页|/personal/profile/[uid]|GET|NO
用户个人信息|/personal/baseinfo/[uid]|GET|YES 
用户个人信息提交保存|/personal/baseinfo/edit/[uid]|POST|YES

### 用户关系
FUNCTION|PATH|METHOD|NEED_LOGIN
---|:--:|:--:|:--:
用户关注列表|/relation/follows/[uid]|GET|NO
用户粉丝列表|/relation/fans/[uid]|GET|NO
发起关注|/relation/follow|POST|YES
取消关注|/relation/unfollow|POST|YES

### 个人微博feed
FUNCTION|PATH|METHOD|NEED_LOGIN
---|:--:|:--:|:--:
个人关注feed页|/feed|GET|YES
个人关注feed页(定领域)|/feed/[domain]|GET|YES

### 直播弹幕
FUNCTION|PATH|METHOD|NEED_LOGIN
---|:--:|:--:|:--:
个人直播页|/chat/videoview/[roomId]|GET|NO
无效直播页|/chat/no/room|GET|NO
注册直播间|/chat/room/register|POST|YES
开始直播|/chat/room/start|POST|YES
结束直播|/chat/room/stop|POST|YES
弹幕接收(长连接)|/chat/barrage/websocket/[roomId]|GET|NO
弹幕发送(长连接)|/chat/barrage/websocket/[roomId]|GET|YES

### 微博广场
FUNCTION|PATH|METHOD|NEED_LOGIN
---|:--:|:--:|:--:
公共热门页|/public/[domain]|GET|NO
公共热门页(定领域)|/public/[domain]|GET|NO

### 微博操作
FUNCTION|PATH|METHOD|NEED_LOGIN
---|:--:|:--:|:--:
用户发布微博|/mblog/send|POST|YES
单条微博详情页|/mblog/[mid—encode]|GET|NO
删除微博|/mlog/delete/[mid—encode]|POST|YES
点赞微博|/mblog/like|POST|YES
转发微博|/mblog/transmit|POST|YES
发布微博评论|/mblog/comment|POST|YES
获取微博评论|/mblog/comment/list|GET|NO

### 搜索
FUNCTION|PATH|METHOD|NEED_LOGIN
---|:--:|:--:|:--:
搜索页|/search|GET|NO
搜索请求|/search|POST|YES
