package cmd

import (
	"fmt"
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

}
