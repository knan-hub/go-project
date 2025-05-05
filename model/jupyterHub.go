package model

import (
	"time"

	"gorm.io/gorm"
)

type JupyterHub struct {
	gorm.Model
	ServerNodeId     int32      `gorm:"not null;index:idx_server_node" json:"serverNodeId"`
	Status           int8       `gorm:"not null;default:1;index:idx_status_port,priority:1;comment:状态,1服务中,2半关闭,3已关闭" json:"status"`
	Type             int8       `gorm:"not null;default:0;uniqueIndex:idx_type_username;comment:类型,0共享,1独享" json:"type"`
	Username         string     `gorm:"size:32;not null;default:share;uniqueIndex:idx_type_username;comment:jupyter用户名,默认是share" json:"username"`
	DockerId         string     `gorm:"size:255;not null;default:'';uniqueIndex;comment:docker container id" json:"dockerId"`
	Port             int32      `gorm:"not null;index:idx_status_port,priority:2;comment:jupyterHub的启动端口" json:"port"`
	Num              int32      `gorm:"not null;default:0;comment:正在运行的用户数" json:"num"`
	StartAt          time.Time  `gorm:"index;comment:启动时间点" json:"startAt"`
	FirstConnectedAt time.Time  `gorm:"default:null;index:idx_connect_time;comment:首次连接时间，超过一小时则变为半关闭，超过三小时自动关闭" json:"firstConnectedAt"`
	ServerNode       ServerNode `gorm:"foreignKey:ServerNodeId"`
}
