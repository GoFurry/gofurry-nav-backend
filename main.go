package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/GoFurry/gofurry-nav-backend/apps/schedule"
	"github.com/GoFurry/gofurry-nav-backend/common"
	cs "github.com/GoFurry/gofurry-nav-backend/common/service"
	"github.com/GoFurry/gofurry-nav-backend/middleware"
	"github.com/GoFurry/gofurry-nav-backend/roof/env"
	routers "github.com/GoFurry/gofurry-nav-backend/router"
	"github.com/gofiber/fiber/v2/log"
	"github.com/kardianos/service"
)

//@title GoFurry-Fiber
//@version v1.0.0
//@description GoFurry-Fiber

var (
	errChan = make(chan error)
)

func main() {
	dir, _ := os.Getwd()

	svcConfig := &service.Config{
		Name:        common.COMMON_PROJECT_NAME,
		DisplayName: "gf-nav",
		Description: "gf-nav",
		Option: service.KeyValue{
			"SystemdScript": `[Unit]
Description=gf-nav (自定义配置)
After=network.target
Requires=network.target

[Service]
Type=simple
WorkingDirectory=` + dir + `/
ExecStart=` + dir + `/gf-nav
Restart=always
RestartSec=30
LogOutput=true
LogDirectory=/var/log/gf-nav
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target`,
		},
	}
	prg := &goFurry{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Error(err)
	}

	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			s.Install()
			log.Info("服务安装成功.")
			return
		}

		if os.Args[1] == "uninstall" {
			s.Uninstall()
			log.Info("服务卸载成功.")
			return
		}
	}

	// 内存限制和 GC 策略
	debug.SetGCPercent(1000)
	debug.SetMemoryLimit(int64(env.GetServerConfig().Server.MemoryLimit << 30))

	// 初始化系统服务
	InitOnStart()
	// 启动系统
	err = s.Run()
	if err != nil {
		log.Error(err)
	}
}

type goFurry struct{}

func InitOnStart() {
	// 初始化 GeoIP 中间件
	middleware.InitGeoIP()
	// 初始化 redis
	cs.InitRedisOnStart()
	// 初始化时间调度
	cs.InitTimeWheelOnStart()

	// 初始化定时任务
	schedule.InitScheduleOnStart()
}

func (gf *goFurry) Start(s service.Service) error {
	go gf.run()
	return nil
}

func (gf *goFurry) run() {
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	// 启动 web
	go func() {

		app := routers.Router.Init()
		addr := env.GetServerConfig().Server.IPAddress + ":" + env.GetServerConfig().Server.Port
		//pem := env.GetServerConfig().Key.TlsPem
		//key := env.GetServerConfig().Key.TlsKey
		//if err := app.ListenTLS(addr, pem, key); err != nil {
		//	fmt.Println(err)
		//	errChan <- err
		//}
		if err := app.Listen(addr); err != nil {
			fmt.Println(err)
			errChan <- err
		}
	}()
	if err := <-errChan; err != nil {
		log.Error(err)
	}
}

func (gf *goFurry) Stop(s service.Service) error {
	return nil
}
