package curd

import (
	"context"

	"github.com/zyw0605688/gozero-curd-vue/service/public/internal/svc"
	"github.com/zyw0605688/gozero-curd-vue/service/public/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HelloWorldLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHelloWorldLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HelloWorldLogic {
	return &HelloWorldLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HelloWorldLogic) HelloWorld(req *types.Empty) (resp *types.HelloWorldResponse, err error) {
	resp = &types.HelloWorldResponse{
		Data: "hello, gozero-curd-vue!",
	}
	return
}
