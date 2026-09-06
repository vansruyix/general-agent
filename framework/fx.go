package framework

import (
	"general-agent/config"
	"general-agent/framework/http"
	logger "general-agent/framework/log"
	"general-agent/framework/mysql"
	"general-agent/framework/postgres"
	"general-agent/framework/scheduler"
	"general-agent/framework/valid"
	"general-agent/framework/xredis"
	"general-agent/framework/xresty"

	"go.uber.org/fx"
)

var Module = fx.Module("framework",
	fx.Provide(config.NewGlobalConfig, postgres.NewPostgresDB, xredis.NewRedisClient, xresty.NewHttpClient),
	mysql.Module,
	logger.Module,
	http.Module,
	valid.Module,
	scheduler.Module,
)
