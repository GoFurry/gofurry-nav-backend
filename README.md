# GoFurry-Nav 后端部分

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