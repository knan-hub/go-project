package cmd

import (
	"fmt"
	"go-project/consts"
	"go-project/dao"
	"go-project/logger"
	"go-project/redis"
	"go-project/router"
	"go-project/setting"
	"os"
	"strings"
)

func Execute() {
	configFile := "settings.yml" // 默认配置文件路径，可以通过命令行参数指定

	if len(os.Args) < 2 {
		fmt.Println("Please provide a command.eg: go-project config/settings.yml")
		return
	}

	if os.Args[1] == "dev" || os.Args[1] == "test" {
		configFile = "config/" + fmt.Sprintf("settings.%s.yml", os.Args[1])
	} else if strings.Contains(os.Args[1], "config") {
		configFile = os.Args[1]
	} else {
		configFile = "config/" + configFile
	}

	// 加载配置文件
	if err := setting.Init(configFile); err != nil {
		fmt.Printf("load config failed, err:%v\n", err)
		return
	}

	consts.Init(setting.Config.VectorDbConfig)

	if err := logger.Init(setting.Config.LogConfig, setting.Config.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}

	if err := dao.Init(setting.Config.VectorDbConfig); err != nil {
		fmt.Printf("init vector db failed, err:%v\n", err)
		return
	}

	redis.Init(setting.Config.RedisConfig)

	r := router.Setup(setting.Config.Mode)
	err := r.Run(fmt.Sprintf(":%d", setting.Config.Port))
	if err != nil {
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}
}
