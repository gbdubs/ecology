package list_project

import (
	"github.com/gbdubs/ecology/manifests/ecology_manifest"
	"github.com/gbdubs/ecology/util/output"
	"os"
	"path/filepath"
	"strings"
)

type ListProjectCommand struct {
	EcologyManifest ecology_manifest.EcologyManifest
	Verbose         bool
}

func (lpc ListProjectCommand) Execute(o *output.Output) error {
	pms, err := listAllEcologyManifestsInDirectory(lpc.EcologyManifest.EcologyProjectsDirectoryPath)
	if err != nil {
		o.Error(err)
		return err
	}
	o.Info("Available Projects:").Indent()
	for _, manifestPath := range pms {
		if lpc.Verbose {
			o.Success(manifestPath)
		} else {
			o.Success(
				strings.Replace(strings.Replace(manifestPath, lpc.EcologyManifest.EcologyProjectsDirectoryPath+"/", "", 1), "/ecology.ecology", "", 1))
		}
	}
	o.Dedent().Done()
	return nil
}

func listAllEcologyManifestsInDirectory(root string) (files []string, err error) {
	err = filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == ".ecology" {
				files = append(files, path)
			}
			return nil
		})
	return files, err
}
