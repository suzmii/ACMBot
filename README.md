# ACMBot

![Badge](https://img.shields.io/badge/OneBot-v11-black)
![Badge](https://img.shields.io/badge/go-%3E%3D1.20-30dff3?logo=go)

## 项目介绍

这是一个使用 Go 语言，基于 onebot11 协议开发的 QQBot 项目，主要提供比赛查询，个人信息查询，群友排行等功能

## TODO

### 个人信息展示

- [x] CodeForces | usage: `cf [username]`
- [x] CodeForces Rating 曲线图 | Usage: `rating [username]`
- [ ] AtCoder
- [ ] NowCoder

### 近期比赛

- [x] CodeForces | usage: `近期cf`
- [x] AtCoder | usage: `近期比赛`
- [x] NowCoder | usage: `近期比赛`
- [x] Luogu | usage: `近期比赛`

### 其他

- [ ] 群内排行
- [ ] ...

## 如何运行

```shell
git clone https://github.com/suzmii/ACMBot
cd ACMBot
go mod tidy
go run ./main.go
```

第一次启动会自动生成配置文件，填写好相关内容之后启动即可正常运行

## 提示

用了部分`Gensokyo`扩展的 api，非`Gensokyo`协议端可能不支持

[TODO] 后期会在配置文件里加开关的
