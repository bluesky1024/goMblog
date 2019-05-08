##  mblog_scheme
### 功能
#### 发布
* 用户登录态可发布微博,字数限制140字内
* 发布的微博可以指定那些人可见(自己可见/好友圈可见/公开)
#### 编辑
* 用户登录态可删除已发布微博
* 用户登录态可重新编辑已发布微博
#### 查阅位置
* 可在用户个人主页查看自身微博
* 可直接通过mid生成短链接查阅
* 可进入指定用户主页查看其所有微博
#### 转评赞
* 每条微博均可被转发，评论(100字内)，点赞
* 对转评赞进行计数并存储
#### 消息通知
* 对于在线用户进行消息推送（websocket???）

---

### 方案设计
#### 微博发布量级预估
* 日均发博数:1000000*5 + 10000000*0.1 + 100000000*0.01 = 700-0000/day
* 人均发博数:500
* 微博用户量:5ww

---

#### 关系型数据库存储设计
##### 微博基础信息
###### 需求
* 根据mid直接定位对应的基本信息表，然后获取微博内容
###### 基本表设计
```sql
CREATE TABLE `mblog_info_201903` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `mid` bigint(18) unsigned NOT NULL COMMENT '微博mid',
  `uid` bigint(18) unsigned NOT NULL COMMENT '发博人uid',
  `content` varchar(255) NOT NULL COMMENT '微博文本内容',
  `origin_uid` bigint(18) unsigned NOT NULL DEFAULT 0 COMMENT '原始发博人uid(用于转发微博)',
  `origin_mid` bigint(18) unsigned NOT NULL DEFAULT 0 COMMENT '原始微博mid(用于转发微博)',
  `create_time` timestamp NOT NULL DEFAULT '1970-01-01 08:00:01',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `trans_cnt` int(10) unsigned NOT NULL COMMENT '转发次数',
  `likes_cnt` int(10) unsigned DEFAULT NULL COMMENT '点赞数',
  `comment_cnt` int(10) unsigned DEFAULT NULL COMMENT '评论数',
  `status` tinyint(4) unsigned NOT NULL DEFAULT 1 COMMENT '微博状态：1.正常；2.删除；3.官方屏蔽',
  `read_able` tinyint(4) unsigned NOT NULL DEFAULT 1 COMMENT '可见性：1.公开； 2.朋友圈可见（互相关注）；3.自己可见',
  PRIMARY KEY (`id`),
  UNIQUE KEY `mid` (`mid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='微博信息表';
```
###### 分库方案
* 待补充
###### 分表方案
* 根据mid中时间戳信息即可获取发博时间，根据该信息进行分表
* 理论上微博信息按日分表存放--mblog_info_20190304
* 测试环境微博信息可按月分表存放--mblog_info_201903

---

##### 微博用户与微博映射关系表
---
###### 需求
* mblog_info表按月存放微博信息，无法根据指定uid获取其所有微博
* 所以需要冗余该表
###### 基本表设计
```sql
CREATE TABLE `uid_to_mblog_1` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `uid` bigint(18) NOT NULL COMMENT '发博人uid',
  `mid` bigint(18) unsigned NOT NULL COMMENT '微博mid',
  `create_time` timestamp NOT NULL DEFAULT '1970-01-01 08:00:01',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `status` tinyint(4) unsigned NOT NULL DEFAULT 1 COMMENT '微博状态：1.正常；2.删除；3.官方屏蔽',
  `read_able` tinyint(4) unsigned NOT NULL DEFAULT 1 COMMENT '可见性：1.公开； 2.朋友圈可见（互相关注）；3.自己可见',
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户微博映射表';
```
###### 分库方案
* 待补充
###### 分表方案
* 冗余该映射表，按发博人uid进行进行分表
* 一致性hash分表,建立10000个虚拟节点，目前对应10张真实表,uid%10000（0～999=>1,1999~2000=>2,...）
* 后期要对部分表进行扩充，可将变动范围控制于局部。

---

##### 微博评论信息
###### 基本表设计
```sql
CREATE TABLE `mblog_comment_201903` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `mid` bigint(18) unsigned NOT NULL COMMENT '微博mid',
  `uid` bigint(18) unsigned NOT NULL COMMENT '发博人id',
  `comment_uid` bigint(18) unsigned NOT NULL COMMENT '评论人uid',
  `content` varchar(255) NOT NULL COMMENT '评论文本内容',
  `create_time` timestamp NOT NULL DEFAULT '1970-01-01 08:00:01',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `likes_cnt` int(10) unsigned DEFAULT NULL COMMENT '点赞数',
  `status` tinyint(4) unsigned NOT NULL DEFAULT 1 COMMENT '微博状态：1.正常；2.删除；3.官方屏蔽',
  PRIMARY KEY (`id`),
  KEY `mid` (`mid`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='微博信息表'
```
###### 分库方案
* 待补充
###### 分表方案
* 根据mid


---

#### 非关系型缓存存储设计
##### 微博点赞数计数