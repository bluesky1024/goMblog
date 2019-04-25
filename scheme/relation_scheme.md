##  relation_scheme
### 功能
* 用户互相的关注/取关
* 获取全部关注/粉丝列表
* 对关注用户进行分组管理(至多20个分组，一个关注用户只能纳入一个分组)

---

### 方案设计
#### 量级预估
* 关注数：2000
* 粉丝数：1e8
#### 存储方式
* mysql备份存储落地数据
* redis的list结构存储热数据 （The max length of a list is 2^32 - 1 elements( 4294967295)）

---

#### 表设计
##### 基本设计
* 以follow_info表中数据为准，在执行关注取关操作时，优先存储到该表中
* 分组信息表，若关注人可以分配个多个分组，需要冗余关注信息，目前暂仅支持关注人只属于一个分组
* 执行关注取关操作时，同时将动作存入kafka，在kafka中完善fan_info表和维护互关表
###### 用户关注分组信息
```sql
CREATE TABLE `follow_group_201903` (
  `id` bigint(18) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `uid` bigint(18) unsigned NOT NULL COMMENT '用户uid',
  `group_name` varchar(64) NOT NULL COMMENT '分组名称',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态：1 有效；0 无效',
  `create_time` timestamp NOT NULL DEFAULT '1970-01-01 08:00:01',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uid_groupname` (`uid`,`group_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户分组信息表' 
```

---

###### 关注信息表
```sql
CREATE TABLE `follow_info_201903` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `uid` bigint(18) unsigned NOT NULL COMMENT '用户uid',
  `follow_uid` varchar(64) NOT NULL COMMENT '关注人uid',
  `create_time` timestamp NOT NULL DEFAULT '1970-01-01 08:00:01',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '关注状态：1 有效；0 无效',
  `group_id` bigint(18) unsigned NOT NULL COMMENT '分组id：0 未分组；其他 分组id'
  `is_friend` tinyint(4) NOT NULL DEFAULT 1 COMMENT '（还是不要了吧。。。）是否互关：1 是；0 不是',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uid_follow_ind` (`uid`,`follow_uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户关注信息表' 
```

---

###### 粉丝信息表
```sql
CREATE TABLE `fan_info_201903` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `uid` bigint(18) unsigned NOT NULL COMMENT '用户uid',
  `fan_uid` varchar(64) NOT NULL COMMENT '粉丝uid',
  `create_time` timestamp NOT NULL DEFAULT '1970-01-01 08:00:01',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '关注状态：1 有效；0 无效',
  `is_friend` tinyint(4) NOT NULL DEFAULT 1 COMMENT '是否互关：1 是；0 不是',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uid_fan_ind` (`uid`,`fan_uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户粉丝信息表' 
```

---

#### 难点
* 粉丝列表人数过多怎么处理？
  - 使用一致性hash进行分表
* kafka生产消费队列崩溃，导致数据丢失怎么处理？
* kafka只能保证在一个partition内的消息是有序的，如果消息被传递到多个不同的partition内可能存在关注和取关等操作的逆序，此处若需要水平扩展：
    * 生产者：可以按照uid进行hash，不同的操作固定放置到特定的partition中
    * 消费者：
