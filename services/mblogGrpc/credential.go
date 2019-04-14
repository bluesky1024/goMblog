package mblogGrpc

import (
	"context"
	//"fmt"
	"github.com/bluesky1024/goMblog/tools/auth"
	"google.golang.org/grpc/metadata"
)

var(
	appId = "123"
	token = "abc"
)

// customCredential 自定义认证
type customCredential struct{}

func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appid":  appId,
	}, nil
}

func (c customCredential) RequireTransportSecurity() bool {
	return false
}


func getSign(ctxOri context.Context,req interface{}) (ctx context.Context) {
	ctx = metadata.AppendToOutgoingContext(ctxOri,"app_id",appId)
	ctx = metadata.AppendToOutgoingContext(ctx,"app_sign",authen.GetSign(req,token))
	return ctx
}
