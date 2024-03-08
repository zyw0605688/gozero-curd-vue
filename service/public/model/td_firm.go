package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TdFirm struct {
	FirmId    string `json:"firm_id" gorm:"primaryKey;column:firm_id;type:varchar(36);comment:唯一id;"`
	FirmName  string `json:"firm_name" gorm:"column:firm_name;not null;uniqueIndex;type:varchar(255);comment:厂商名称;"`
	FirmAlias string `json:"firm_alias" gorm:"column:firm_alias;type:varchar(255);comment:厂商别名;"`
	FirmCode  string `json:"firm_code" gorm:"column:firm_code;type:varchar(255);comment:厂商编码;"`
	FirmDesc  string `json:"firm_desc" gorm:"column:firm_desc;type:varchar(255);comment:厂商描述;"`
}

func (table TdFirm) TableName() string {
	return "td_firm"
}

func (u *TdFirm) BeforeCreate(tx *gorm.DB) (err error) {
	u.FirmId = uuid.NewString()
	return
}
