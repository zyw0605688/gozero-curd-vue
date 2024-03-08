package main

import (
	"flag"
	"fmt"

	"github.com/zyw0605688/gozero-curd-vue/service/public/internal/config"
	"github.com/zyw0605688/gozero-curd-vue/service/public/internal/handler"
	"github.com/zyw0605688/gozero-curd-vue/service/public/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/public-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
