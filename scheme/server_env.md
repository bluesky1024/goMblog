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

