package mysql

import (
	"fmt"
	"go-project/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func Init(cfg *setting.MySQLConfig) {
	connect := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DatabaseName,
	)

	var err error
	DB, err = gorm.Open("mysql", connect)
	if err != nil {
		panic("MySQL初始化失败! " + err.Error())
	}

	DB.DB().SetMaxIdleConns(cfg.MaxIdleConns)       // 连接池最大允许的空闲连接数
	DB.DB().SetMaxOpenConns(cfg.MaxOpenConns)       // 设置最大连接数
	DB.DB().SetConnMaxLifetime(cfg.ConnMaxLifetime) // 设置了连接可复用的最大时间
	DB.SingularTable(true)                          // 表名禁用复数
}
