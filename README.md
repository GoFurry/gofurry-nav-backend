# GoFurry-Nav 后端部分

### 项目地址

GoFurry导航站 go-furry.com APII地址 api.go-furry.com

### 环境依赖

|     Go版本     | 数据库 | 缓存中间件 |
| :------------: | :----: | :--------: |
| GO SDK 1.23.0+ | PGSQL  |   Redis    |

### 部署步骤

1. go build 获取二进制文件 gf-nav
2. chmod +777 ./gf-nav
3. ./gf-nav install 将服务注册到 systemd
4. systemctl enable gf-nav 这里的 gf-nav 可在 main.go 中进行修改
5. systemctl start gf-nav

### 目录结构描述

│  go.mod
│  go.sum
│  main.go							// 程序入口
├─apps						// MVC结构
│  ├─nav
│  │  ├─navPage
│  │  └─sitePage
│  ├─schedule				// 定时任务
│  │  │  schedule.go
│  │  └─task
│  │          statTask.go
│  └─system
│      ├─stat
│      └─user
│          ├─controller
│          ├─dao
│          ├─models
│          └─service
├─bin						// go build 的二进制文件
│      gf-nav
├─common					// 公共类库
│  │  constant.go
│  │  error.go
│  │  response.go
│  ├─abstract
│  │      dao.go
│  │      model.go
│  │      validate.go
│  ├─log
│  ├─models
│  │      email.go
│  │      jwtClaims.go
│  │      query.go
│  │      time.go
│  │      wsModel.go
│  ├─service
│  │      dataEventService.go
│  │      emailService.go
│  │      kafkaService.go
│  │      oauthService.go
│  │      redisService.go
│  │      timewheelService.go
│  └─util
│          function.go
│          http.go
│          parse.go
├─conf							// 配置文件
│      coraza.conf
│      server.build.yaml
│      server.yaml
│      tls.csr
│      tls.key
│      tls.pem
├─data							// GeoLite离线数据
│      GeoLite2-ASN.mmdb
│      GeoLite2-City.mmdb
│      GeoLite2-Country.mmdb
│      
├─docs							// swagger
│      docs.go
│      swagger.json
│      swagger.yaml
├─middleware
│      corazaWAF.go					// corazaWAF 中间件
│      geoip.go						// geoip 中间件
├─roof
│  ├─db							 // 数据库连接	
│  │      db.go
│  └─env							// 读取配置文件
│          config.go
└─router							// 路由层