## server_run(wait for shell...)
### mysql
```bash
	mysql.server start
```

### redis
```bash
	redis-server
```

### kafka
```bash
    docker-compose up
```

### goMblog(调试中，暂未编译)
```bash
    #grpc相关服务
    cd $GO_PATH/src/github.com/bluesky1024/goMblog/grpcServer/user
    go run auth.go server.go userGrpc.go
    cd $GO_PATH/src/github.com/bluesky1024/goMblog/grpcServer/mblog
    go run auth.go server.go mblogGrpc.go
    cd $GO_PATH/src/github.com/bluesky1024/goMblog/grpcServer/relation
    go run auth.go server.go relationGrpc.go

    #前端web服务
    cd $GO_PATH/src/github.com/bluesky1024/goMblog
    go run main.go route.go initServ.go
```



