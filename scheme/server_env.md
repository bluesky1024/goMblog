## go-nginx 环境配置
### 基础组件
```bash
#基础镜像下载
docker pull centos:latest
#运行centos
docker run -it --name="go_mblog" -p80:80 centos:latest
yum -y install vim wget git
yum -y install gcc pcre pcre-devel zlib zlib-devel openssl openssl-devel
```

### go环境搭建
* go压缩包下载
```bash
wget https://dl.google.com/go/go1.12.linux-amd64.tar.gz
```
* 解压至指定文件夹
```bash
tar -C /usr/local/ -xzf go1.12.linux-amd64.tar.gz
```
* 配置环境变量
```bash
vim ~/.bash_profile
#export GOPATH=$HOME/go
#export GOROOT=/usr/local/go
#export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
source ~/.bash_profile
```
* go开发文件夹
```bash
mkdir -p $HOME/go/src
mkdir -p $HOME/go/bin
mkdir -p $HOME/go/pkg
```
* 测试go
```bash
go env
```

### nginx环境搭建
* nginx压缩包下载
```bash
wget http://nginx.org/download/nginx-1.15.9.tar.gz
```
* 解压至指定文件夹
```bash
tar -C /usr/local -xzf nginx-1.15.9.tar.gz
mv nginx-1.15.9 nginx
```
* nginx编译
```bash
./configure --prefix=/usr/local/nginx --conf-path=/usr/local/nginx/nginx.conf
make
make install
```
* nginx启动关闭
```bash
cd /usr/local/nginx/sbin/
./nginx 
./nginx -s stop
./nginx -s quit
./nginx -s reload
```
* 测试nginx
```bash
/usr/local/nginx/sbin/nginx
#浏览器中查看localhost:welcome to nginx! 
```

### mysql环境搭建
* mysql_rpm包下载
```bash
wget https://dev.mysql.com/get/mysql80-community-release-el7-2.noarch.rpm
```
...待续

### hbase环境搭建
* java安装
```bash
yum -y install java-1.8.0-openjdk*
# 找到java安装位置
whereis java
# java: /usr/bin/java /usr/lib/java /etc/java /usr/share/java /usr/share/man/man1/java.1.gz
ll /usr/bin/javac
# lrwxrwxrwx 1 root root 23 Mar 14 10:01 /usr/bin/javac -> /etc/alternatives/javac
ll /etc/alternatives/javac
# lrwxrwxrwx 1 root root 70 Mar 14 10:01 /etc/alternatives/javac -> /usr/lib/jvm/java-1.8.0-openjdk-1.8.0.201.b09-2.el7_6.x86_64/bin/javac
export $JAVA_HOME=/etc/alternatives/javac -> /usr/lib/jvm/java-1.8.0-openjdk-1.8.0.201.b09-2.el7_6.x86_64/  > /etc/bash_profile
```
* hbase安装
```bash
wget https://mirrors.cnnic.cn/apache/hbase/stable/hbase-1.4.9-bin.tar.gz
tar -xvf hbase-1.4.9-bin.tar.gz
mv hbase-1.4.9 ~/hbase
# 修改系统配置文件
vim /etc/porifle
export HBASE_HOME=/root/hbase
export PATH=.:${HBASE_HOME}/bin:$PATH
# hbase配置文件
vim ~/hbase/conf/hbase-env.sh
export JAVA_HOME=/usr/lib/jvm/java-1.8.0-openjdk-1.8.0.201.b09-2.el7_6.x86_64/
```
* hbase基本操作
```bash
# 启动hbase
~/hbase/bin/start-hbase.sh
# 进入hbase终端
~/hbase/bin/hbase shell
# 终端中基础指令
# 
```

### redis—cluster环境搭建
* 参考链接
https://www.cnblogs.com/wuxl360/p/5920330.html

* 搭建步骤
```bash
# redis服务安装
# mac
brew install redis
# centos
wget http://download.redis.io/releases/redis-3.2.4.tar.gz
tar -zxvf redis-3.2.4.tar.gz
yum install gcc gcc-c++ automake autoconf libtool 
make && make install

# 以下为集群搭建
mkdir redis_cluster

# redis-cluster 要求至少6个节点
mkdir 10011 10012 10013 10014 10015 10016
cp redis.conf redis_cluster/10011 ...

# 编辑redis.conf
port  10011                             //端口10011      
bind 本机ip                              //默认ip为127.0.0.1 需要改为其他节点机器可访问的ip 否则创建集群时无法访问对应的端口，无法创建集群
daemonize    yes                        //redis后台运行
pidfile  /var/run/redis_10011.pid       //pidfile文件对应10011
cluster-enabled  yes                    //开启集群  把注释#去掉
cluster-config-file  nodes_10011.conf   //集群的配置  配置文件首次启动自动生成 10011
cluster-node-timeout  15000             //请求超时  默认15秒，可自行设置
appendonly  yes                         //aof日志开启  有需要就开启，它会每次写操作都记录一条日志　
# 其余5个节点类似

# 启动6个redis服务
redis-server 10011/redis.conf ...

# 使用redis提供的官方ruby集群管理工具(对应redis安装包/src/redis-trib.rb)
# centos
yum -y install ruby ruby-devel rubygems rpm-build
# mac
brew install ruby

# 安装ruby对于redis的支持库函数
gem install redis

# 运行脚本进行redis集群管理
ruby redis-trib.rb create --replicas 1 127.0.0.1:10011 127.0.0.1:10012 127.0.0.1:10013 127.0.0.1:10014 127.0.0.1:10015 127.0.0.1:10016

# 新增主节点
redis-trib.rb add-node 127.0.0.1:10017 127.0.0.1:10011

# 查看新增的主节点
redis-cli -c -p 10011 CLUSTER nodes | grep 10017
301b60cdb455b9ae27b7b562524c0d039e640815 127.0.0.1:10017 master - 0 1487342302506 0 connected

# 新增从节点
redis-trib.rb add-node  
--slave --master-id 301b60cdb455b9ae27b7b562524c0d039e640815 127.0.0.1:10018
 192.168.11.3:6380

 # 查看整个集群的状态
redis-cli -c -p 10011 CLUSTER nodes

# 换了ip但端口没变动的情况下，重启ruby redis-trib.rb会报node不为空
# 节点信息是保存在redis中，重新配置了集群，可以考虑把rdb，aof文件清空，对redis执行flushdb，然后重启redis-server
# 配置文件和日志文件保存在 /usr/local/var/db/redis/
```


