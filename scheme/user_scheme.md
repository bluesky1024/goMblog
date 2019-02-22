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
) ENGINE=InnoDB AUTO_INCREMENT=48 DEFAULT CHARSET=utf8 COMMENT='用户基本信息表' 
```
#### 用户详细信息
CREATE TABLE `user_detail` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `uid` int(10) unsigned NOT NULL COMMENT '用户uid',
)

##### 重点字段
字段|意义|长度
---|:--:|---:
uid|用户id|
nick_name|用户昵称|内容
password|用户密码hash值|内容