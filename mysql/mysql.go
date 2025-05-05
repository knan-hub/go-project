package mysql

import (
	"database/sql"
	"fmt"
	"go-project/model"
	"go-project/setting"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var sqlDB *sql.DB

func Init(cfg *setting.MySQLConfig) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DatabaseName,
	)

	var err error
	DB, err = gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{LogLevel: logger.Info}),
		},
	)
	if err != nil {
		panic("MySQL初始化失败! " + err.Error())
	}

	// 获取通用数据库对象，然后设置连接池参数
	sqlDB, err = DB.DB()
	if err != nil {
		panic("获取数据库连接池失败! " + err.Error())
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)       // 连接池最大允许的空闲连接数
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)       // 设置最大连接数
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime) // 设置了连接可复用的最大时间

	syncToDB()

	return
}

func syncToDB() {
	DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.ServerNode{})
	DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&model.JupyterHub{})
}

func Close() {
	_ = sqlDB.Close()
}
