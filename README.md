# 微信 IPAD 协议

## 环境要求

1. Linux 或者 Windows
2. 服务器配置任意，根据需求自行调整（推荐 224）
3. 软件：redis。 mysql。
4. 开发环境：GO

## 打包

#### --linux

```
brew install FiloSottile/musl-cross/musl-cross
//  打包   assets  static templates 目录要手动创建

CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc CGO_LDFLAGS="-static"  go build -a main.go

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
```

#### windows

```
go build -v
```

## 在线文档

[http://181.214.39.230:8888/02/](http://181.214.39.230:8888/02/)

## 配置文件

```
/assets/setting.json
```

## Mysql

```
"mySqlConnectStr": "账号:密码@tcp(IP:端口)/库?charset=utf8mb4&parseTime=true&loc=Local&time_zone='Asia%2Fshanghai'",
```
