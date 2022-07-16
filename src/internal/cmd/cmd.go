package cmd

import (
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"res"
	"strings"

	"github.com/IfanTsai/go-lib/utils/byteutils"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
)

var gitRepoURLRegexp = regexp.MustCompile(`((git|ssh|http(s)?)|(git@([\w.]+)))(:(//)?)([\w.@:/\-~]+)(\.git)`)

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
	{
		Name: "git-repo-url",
		Prompt: &survey.Input{
			Message: "What is your git repo url? (optional, enter to skip)",
			Default: "",
		},
		Validate: func(val interface{}) error {
			if str, ok := val.(string); ok {
				if len(str) == 0 {
					return nil
				}

				if !gitRepoURLRegexp.Match(byteutils.S2B(str)) {
					return errors.New("invalid git repo url")
				}
			} else {
				return errors.Errorf("cannot convert the type %v to a string", reflect.TypeOf(val).Name())
			}

			return nil
		},
	},
}

func Run() (*res.ProjectInfo, error) {
	var answer res.ProjectInfo
	if err := survey.Ask(question, &answer); err != nil {
		return nil, errors.Wrap(err, "failed to ask questions")
	}

	answer.ProjectPath = filepath.Clean(answer.Path + "/" + answer.ProjectName)
	answer.ProjectName = strings.ReplaceAll(answer.ProjectName, "-", "_")
	answer.GoModulePath = answer.ModuleName

	if len(answer.GitRepoURL) != 0 {
		matches := gitRepoURLRegexp.FindStringSubmatch(answer.GitRepoURL)
		if strings.HasPrefix(matches[1], "http") {
			answer.GoModulePath = matches[len(matches)-2]
		} else {
			strs := strings.Split(matches[1], "@")
			if len(strs) == 2 {
				answer.GoModulePath = path.Join(strs[1], matches[len(matches)-2])
			}
		}
	}

	return &answer, nil
}
