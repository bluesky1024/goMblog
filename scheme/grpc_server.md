## grpc_server
### 需求
* 将services下多个服务均纳入gprc服务中，彼此之间通过gprc进行服务调用，以便后续进行微服务化
* 多个服务之间怎么解耦合？？？
* 提供服务和调用服务的安全性？？？
* 服务发现和服务治理？？？

### 实现
#### 解耦合
* 怎么减少主服务对其他服务对依赖？

#### 安全性
* TLS认证（通用配置）
    * 服务端配置证书及私钥
    * 客户端安装证书
* 自定义签名认证 （根据不同服务接口权限定制）

#### 服务发现/治理

ListenAndServe()

listen部分：
    本地地址解析
    判断ip协议（AF_INET|AF_INET6|...）
    建立socket
        系统调用建立系统文件描述符（系统调用需要加forklock读写锁）还对系统文件描述符设置了属性为非阻塞，不明觉厉
        对系统文件描述符进行包装（装入网络文件描述符里）
    一个tcp的监听者里其实就放了个网络文件描述符

server部分：
    上述监听者加一层包装，成keepalive监听者（怎么keep alive???）
    server函数接受一个参数（实现了 accept(),close(),addr() 函数的 net.listen）
    keepalive监听者又包了一层。。。成了 onceCloseListener， 应该是只关闭一次，通过sync包进行控制
    建立个无脑for死循环，除非error，一直进行 listener.Accept() （返回一个连接 conn 或者error） 监听
    conn 主要支持以下操作 （read() write() close()...）
    每次 accept() 一个请求，开一个 go 协程 go conn.server(ctx)
    conn.server(ctx) 懒得看了，中间主要是包括一个 serverHandler{c.server}.ServeHTTP(w, w.req) 这个在开始服务之前进行注册
    
listener.Accept() => netFd.Accept() => Fd.Accept() => poll.accept(sysFd int) => syscall.Accept(sysFd int)