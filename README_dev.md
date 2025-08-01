# 开发者文档

首先，非常感谢，也非常欢迎，你能够点看本项目的源码，以及这个说明文档。

这个文档主要是概括一下项目的大致设计思路，如果你有什么建议，无论是架构设计，代码风格，实现逻辑，产品效果，以及任何你想说的，都欢迎你提出宝贵的意见！

也非常欢迎你加入我们：<a href="https://qm.qq.com/cgi-bin/qm/qr?k=nceseONriNuKcB8yyyIdHLYzv7PKdGMB&jump_from=webapi&authKey=r3OeD12gQueeER8tvH5dp7Sx1DcIzb0pgxu6iUgLo3HP3AnnqS71oslV4v6fBnKv">
ACMBot交流会</a>

### 正题

#### 1. 怎么把项目跑起来？

首先 clone 代码，同步依赖，然后运行生成示例配置文件

```shell
git clone https://github.com/suzmii/ACMBot
cd ACMBot
go mod tidy
go run main.go
```

然后编写配置文件内的信息，再次运行应该就能正常运行了

这里有个