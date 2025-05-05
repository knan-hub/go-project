package model

import "gorm.io/gorm"

type ServerNode struct {
	gorm.Model
	InternalIp string `gorm:"size:255;not null;default:'';comment:内网ip" json:"internalIp"`
	ExternalIp string `gorm:"size:255;not null;default:'';" json:"externalIp"`
	Num        int32  `gorm:"not null;default:0;comment:运行的docker数" json:"num"`
	Status     int8   `gorm:"not null;default:1;comment: 连接状态,0掉线,1在线,2禁用" json:"status"`
}
