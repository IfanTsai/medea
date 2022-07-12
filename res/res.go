package res

import (
	"embed"
)

type TypeFramework string

const (
	TypeFrameworkGin TypeFramework = "gin"
)

type ProjectInfo struct {
	Framework    string `survey:"framework"`
	ProjectName  string `survey:"project-name"`
	ModuleName   string `survey:"module-name"`
	Path         string `survey:"path"`
	ProjectPath  string
	GoModulePath string
}

type ProjectTemplate struct {
	Framework                 TypeFramework
	EmbedFs                   *embed.FS
	RootPath                  string
	DescriptionAfterGenerated string
}

//go:embed gin-template/*
var ginTemplate embed.FS

var ProjectTemplateMap = map[TypeFramework]*ProjectTemplate{
	TypeFrameworkGin: {
		Framework: TypeFrameworkGin,
		EmbedFs:   &ginTemplate,
		RootPath:  "gin-template",
		DescriptionAfterGenerated: `
Scaffolding project in {{ .ProjectPath }}...

Done. Now run:

  cd {{ .ProjectPath }}
  go mod tidy
  ENV=dev go run main.go

`,
	},
}
