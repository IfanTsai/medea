package cmd

import (
	"res"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
)

var question = []*survey.Question{
	{
		Name: "framework",
		Prompt: &survey.Select{
			Message: "Select a framework:",
			Options: []string{string(res.TypeFrameworkGin)},
			Default: string(res.TypeFrameworkGin),
		},
	},
	{
		Name:     "project-name",
		Prompt:   &survey.Input{Message: "What is your project name ?"},
		Validate: survey.Required,
	},
	{
		Name:     "module-name",
		Prompt:   &survey.Input{Message: "What is your project module name?"},
		Validate: survey.Required,
	},
	{
		Name: "path",
		Prompt: &survey.Input{
			Message: "What is your project path?",
			Default: "./",
		},
		Validate: survey.Required,
	},
}

func Run() (*res.ProjectInfo, error) {
	var answer res.ProjectInfo
	if err := survey.Ask(question, &answer); err != nil {
		return nil, errors.Wrap(err, "failed to ask questions")
	}

	answer.ProjectPath = answer.Path + "/" + answer.ProjectName
	answer.GoModulePath = answer.ModuleName

	return &answer, nil
}
