## user_scheme
### 功能
#### 注册
* 手机邮箱注册（不做验证）
* 填写内容：昵称/手机/邮箱/登录密码
#### 登录
* 访问域名下所有域名均能获取到用户身份
#### 附属要求
* uid唯一

### 方案设计
/usr/local/etc/my.cnf
/usr/local/Cellar/mysql/8.0.15/.bottle/etc/my.cnf
#### 用户基本信息
##### 表结构设计
```sql
CREATE TABLE `user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `uid` bigint(18) unsigned NOT NULL COMMENT '用户uid',
  `nick_name` varchar(64) NOT NULL COMMENT '昵称',
  `telephone` varchar(32) NOT NULL COMMENT '手机号',
  `email` varchar(64) NOT NULL COMMENT '邮箱',
  `create_time` timestamp NOT NULL DEFAULT '1970-01-01 08:00:01',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `password` varchar(256) NOT NULL COMMENT '用户密码hash值',
  `profile_image` varchar(256) DEFAULT NULL COMMENT '用户头像',
  `follows_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '粉丝数',
  `friends_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '关注数',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uid` (`uid`),
  UNIQUE KEY `nick_name` (`nick_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户基本信息表' 
```
##### 分表
* 一致性hash分表,与uid_to_mblog表一致(暂无法这么写，还没想好怎么根据nick_name进行查询)

#### 用户详细信息
CREATE TABLE `user_detail` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `uid` int(10) unsigned NOT NULL COMMENT '用户uid',
)

#### 用户粉丝数
##### 分析
* 用户粉丝数相较用户基本信息，属于变动更高频的对象
* 变动过程由redis进行维护，数据最终落地到mysql，什么时候落地？定时批量落地？？？

##### 关键问题
* 降低存储成本？
* 提升存储速度？
* 保证可靠性？

##### 一些瞎想
* 是否采用redis集群进行存储？
  - redis集群优势：自动主从同步，自动根据key来hash分节点，当主节点挂了时，从节点可以转正，变成主节点来提供服务，方便容灾
  - 问题：redis节点变更的时候，如增加节点，移除节点，怎么保证数据的正常获取？参考：https://blog.51cto.com/xiaotaoge/1899800 。似乎扩容过程，数据的读写不受影响。
  - 扩展一下，顺便考虑一下mysql扩充分表(一致性hash)：对于table_1,tale_2,table_3中的table_1分表为table_1_1,table_1_2,怎么平滑扩容？？？
    + 提前进行表的双写，在指定时间后对数据的变更同时写入两张表
    + 将该时间点之前的数据复制到table_1_2
    + table_1_2正式提供读服务，分表结束
  

 



##### 重点字段
字段|意义|长度
---|:--:|---:
uid|用户id|
nick_name|用户昵称|内容
password|用户密码hash值|内容