package schedule

import (
	"fmt"
	"time"

	"github.com/GoFurry/gofurry-nav-backend/apps/schedule/task"
	"github.com/GoFurry/gofurry-nav-backend/common/log"
	cs "github.com/GoFurry/gofurry-nav-backend/common/service"
)

// 初始化
func InitScheduleOnStart() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Sprintf("receive InitScheduleOnStart recover: %v", err))
		}
	}()
	log.Info("Schedule 模块初始化开始...")

	//初始化后执行一次 Schedule
	go Schedule()
	// 定时任务执行 Schedule
	cs.AddCronJob(10*time.Minute, Schedule)

	log.Info("Schedule 模块初始化结束...")
}

// 任务表
func Schedule() {
	// task 任务
	task.UpdateTopCountCache()
	task.UpdateLatestPingLog()
}
