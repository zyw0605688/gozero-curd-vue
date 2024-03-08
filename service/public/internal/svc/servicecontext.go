package svc

import (
	"github.com/zyw0605688/gozero-curd-vue/service/public/internal/config"
	"github.com/zyw0605688/gozero-curd-vue/service/public/model"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config   config.Config
	PublicDb *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		PublicDb: model.NewGorm(c.PublicDb),
	}
}
