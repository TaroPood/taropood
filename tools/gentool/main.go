package main

import (
	"os"

	"github.com/TaroPood/taropood/internal/repository/postgres/model"
	"gorm.io/gen"
)

func main() {
	outPath := "./internal/repository/postgres/query"
	if len(os.Args) > 1 {
		outPath = os.Args[1]
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: outPath,
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.ApplyBasic(
		&model.RuleModel{},
		&model.ActionModel{},
	)

	g.Execute()
}
