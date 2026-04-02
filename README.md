# ACMBot

ACMBot 是一个面向算法竞赛场景的 QQ 机器人，聚合 AtCoder、Codeforces 等平台的数据，支持用户资料查询、Codeforces rating 曲线渲染，以及近期比赛查询。

## 功能

- AtCoder 用户资料卡片查询
- Codeforces 用户资料卡片查询
- Codeforces rating 历史曲线渲染
- 聚合近期比赛信息，支持按平台筛选
- 使用定时任务同步比赛数据
- 使用 PostgreSQL 缓存用户与比赛数据

## 技术栈

- Go
- ZeroBot
- PostgreSQL
- sqlc
- Viper
- Playwright
- Colly
- req/v3

## 支持的查询

机器人当前内置了以下消息指令：

- `at <用户名>`: 查询 AtCoder 用户资料图片
- `cf <用户名>`: 查询 Codeforces 用户资料图片
- `rt <用户名>` 或 `rating <用户名>`: 查询 Codeforces rating 曲线图
- `race`: 查询近期比赛
- `近期比赛`: 查询近期比赛
- `近期<平台>`: 查询指定平台近期比赛，例如 `近期cf`、`近期atcoder`、`近期牛客`
- `近期比赛 <页码>`: 按页查看比赛列表

支持的平台别名包括：

- Codeforces: `codeforces` / `cf`
- AtCoder: `atcoder` / `at`
- LeetCode: `leetcode` / `lc` / `力扣`
- Luogu: `lg` / `洛谷`
- NowCoder: `nk` / `牛客`

## 工作原理

- 启动时读取 `config.toml`
- 初始化日志、数据库、API 客户端、渲染器和调度器
- 调度器周期性拉取 Clist 比赛数据并写入数据库
- 收到聊天消息后，ZeroBot 按正则规则路由到对应处理逻辑
- 用户资料与 rating 数据按需从远端 API 拉取，并缓存到 PostgreSQL
- 图片由 Playwright 加载 HTML 模板后截图生成

## 环境要求

- Go 1.24.4
- PostgreSQL
- 可运行 Chromium 的环境
- OneBot 兼容实现，并提供 WebSocket Client 连接入口

首次启动渲染模块时会自动安装 Playwright 所需的 Chromium。

## 快速开始

### 1. 安装依赖

确保本地已经准备好：

- Go
- PostgreSQL
- OneBot 服务端

### 2. 配置数据库

创建一个可访问的 PostgreSQL 数据库，并准备连接串。

默认配置中的 DSN 为：

```toml
postgres://postgres:postgres@localhost:5432/postgres
```

### 3. 配置机器人

程序首次运行时，如果仓库根目录不存在 `config.toml`，会自动生成默认配置并退出。

按需填写以下关键配置：

```toml
[api]
codeforces_key = "<your-codeforces-key>"
codeforces_secret = "<your-codeforces-secret>"
clist_authenticated = "<your-clist-authorization-header>"

[database]
dsn = "postgres://postgres:postgres@localhost:5432/postgres"

[zerobot]
command_prefix = "/"
host = "localhost"
port = 15630
token = ""

[render]
headless = true
poolSize = 8
```

说明：

- `api.codeforces_key` 和 `api.codeforces_secret` 用于访问 Codeforces API
- `api.clist_authenticated` 用于访问 Clist API
- `zerobot.host`、`zerobot.port`、`zerobot.token` 需要与 OneBot 服务端配置一致

## 运行

```bash
go run .
```

或先构建再运行：

```bash
go build -o acmbot .
./acmbot
```

## 项目结构

```text
.
├── api/        # 第三方平台接口封装
├── config/     # 配置加载与默认配置
├── database/   # 数据访问层与 sqlc 生成代码
├── handler/    # 业务处理逻辑
├── render/     # HTML 模板渲染与截图
├── scheduler/  # 定时任务调度
├── util/       # 通用工具
├── main.go     # 程序入口
└── zerobot.go  # QQ 消息规则与机器人启动
```

## 注意事项

- 项目依赖外部平台接口，网络异常时部分功能可能不可用
- Codeforces 与 AtCoder 用户数据会被缓存，避免重复请求
- 比赛查询依赖定时任务同步数据库中的比赛数据
- 若首次渲染较慢，通常是 Playwright 浏览器安装或启动造成的

## 开发

运行测试：

```bash
go test ./...
```

如需更新 sqlc 生成代码，请根据仓库中的 `sqlc.yaml` 重新生成。
