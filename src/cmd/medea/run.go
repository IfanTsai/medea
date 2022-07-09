package main

import (
	"log"
	"res"

	"medea/internal/cmd"
	"medea/internal/projectgen"
)

func main() {
	projectInfo, err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}

	projectTemplate := res.ProjectTemplateMap[res.TypeFramework(projectInfo.Framework)]
	projectgen.GenerateProject(projectTemplate, projectInfo)
}
