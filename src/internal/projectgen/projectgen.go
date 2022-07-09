package projectgen

import (
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"res"
	"strings"

	"github.com/IfanTsai/go-lib/utils/byteutils"
	"github.com/pkg/errors"
)

func GenerateProject(projectTemplate *res.ProjectTemplate, projectInfo *res.ProjectInfo) {
	embedFs, err := fs.Sub(projectTemplate.EmbedFs, projectTemplate.RootPath)
	if err != nil {
		log.Fatalln(err)
	}

	if err := fs.WalkDir(embedFs, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		dstPath := fmt.Sprintf("%s/%s", projectInfo.ProjectPath, path)
		dstPath = strings.ReplaceAll(dstPath, "@", "")
		dstPath, err = generateString(dstPath, projectInfo)
		if err != nil {
			return err
		}

		return copyFile(embedFs, projectInfo, path, dstPath)
	}); err != nil {
		log.Fatalln(err)
	}

	if err := executeCmdInProject(projectInfo); err != nil {
		os.RemoveAll(projectInfo.Path)
		log.Fatalln(err)
	}
}

func copyFile(embedFs fs.FS, projectInfo *res.ProjectInfo, srcPath, dstPath string) error {
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}

	f, err := embedFs.Open(srcPath)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}

	srcBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return errors.Wrap(err, "failed to read file")
	}

	srcData, err := generateString(byteutils.B2S(srcBytes), projectInfo)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dstPath, byteutils.S2B(srcData), 0644)
}

func generateString(str string, data interface{}) (string, error) {
	tmpl := template.Must(template.New("").Parse(str))

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", errors.Wrap(err, "failed to execute template")
	}

	return buf.String(), nil
}

func executeCmdInProject(projectInfo *res.ProjectInfo) error {
	script := strings.Join([]string{
		fmt.Sprintf("ln -sf cmd/%s/run.go main.go", projectInfo.ModuleName),
	}, "\n")

	cmd := exec.Command("/bin/sh", "-c", script)
	cmd.Dir = projectInfo.ProjectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "failed to run script: %s", script)
	}

	return nil
}
