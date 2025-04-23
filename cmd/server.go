package cmd

import (
	"context"
	"fmt"
	"go-project/consts"
	"go-project/logger"
	"go-project/redis"
	"go-project/router"
	"go-project/setting"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

func Execute() {
	// 默认配置文件路径，可以通过命令行参数指定
	configFile := "settings.yml"

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

	consts.Init(&setting.Config.VerctorDB)

	if err := logger.Init(&setting.Config.Log, setting.Config.Application.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}

	// 向量数据库初始化
	// if err := dao.Init(&setting.Config.VerctorDB); err != nil {
	// 	fmt.Printf("init vector db failed, err:%v\n", err)
	// 	return
	// }

	redis.Init(&setting.Config.Redis)

	r := router.Init(setting.Config.Application.Mode)

	// Gin框架的Run方法会启动一个HTTP服务器，监听指定的端口，并处理所有的HTTP请求
	// err := r.Run(fmt.Sprintf(":%d", setting.Config.Application.Port))
	// if err != nil {
	// 	fmt.Printf("run server failed, err:%v\n", err)
	// 	return
	// }

	// 创建一个HTTP服务器实例，用于监听指定的端口，并处理所有的HTTP请求
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.Config.Application.Port),
		Handler:        r,
		ReadTimeout:    setting.Config.Application.ReadTimeout,
		WriteTimeout:   setting.Config.Application.WriteTimeout,
		MaxHeaderBytes: 1 << 20, // 1MB,
	}

	// 使用goroutine启动服务器，以便在主goroutine中继续执行其他操作
	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	// 使用通道来等待中断信号，以便在接收到中断信号时优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	// 接收中断信号（Ctrl+C）
	signal.Notify(quit, os.Interrupt)
	// 阻塞，直到接收到中断信号
	<-quit

	fmt.Println("Shutdown Server ...")

	// 创建一个超时上下文，用于控制关闭服务器的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅地关闭服务器
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Server Shutdown:", err)
	}

	// 打印服务器退出的消息，用于确认服务器已成功关闭
	fmt.Println("Server exiting")
}
