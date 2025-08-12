package main

import (
	"github.com/suzmii/ACMBot/internal/database"
	codeforcesrepo "github.com/suzmii/ACMBot/internal/database/dbmodel"
	"go.uber.org/mock/mockgen/model"
	"gorm.io/gen"
)

func main() {
	// 连接数据库（换成你自己的）

	db := database.Init()

	// 初始化 Generator
	g := gen.NewGenerator(gen.Config{
		OutPath:      "internal/database/gen", // 输出路径
		ModelPkgPath: "internal/database/gen", // core 所在包（必须指向已有结构体）
		Mode:         gen.WithoutContext | gen.WithDefaultQuery,
	})

	g.UseDB(db)

	// 只针对已有的结构体生成代码
	g.ApplyInterface(func(model.Method) {}, codeforcesrepo.CodeforcesModels...) // 你可以添加多个 struct
	// 确保 races 也纳入生成，避免手改 gen 代码
	g.ApplyBasic(new(codeforcesrepo.Races))

	// 生成代码
	g.Execute()
}
