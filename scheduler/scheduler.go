package scheduler

import (
	"github.com/jasonlvhit/gocron"
)

func Init() {
	cron := gocron.NewScheduler()

	// 每秒执行
	cron.Every(1).Second().Do(taskFunction)

	// 每分钟执行
	cron.Every(1).Minute().Do(taskFunction)

	// 每小时执行
	cron.Every(1).Hour().Do(taskFunction)

	// 每天固定时间执行（例如每天10:30）
	cron.Every(1).Day().At("10:30").Do(taskFunction)

	// 每周特定星期几执行（例如每周一10:30）
	cron.Every(1).Monday().At("10:30").Do(taskFunction)

	// 每月特定日期执行（例如每月1号10:30）
	cron.Every(30).Days().At("10:30").Do(taskFunction)

	// 启动定时任务
	go cron.Start()
}

// 示例任务函数
func taskFunction() {
	// 这里写定时执行的逻辑
}
