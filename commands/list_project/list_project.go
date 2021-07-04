package list_project

import (
	"github.com/gbdubs/ecology/manifests/ecology_manifest"
	"github.com/gbdubs/ecology/util/output"
)

type ListProjectCommand struct {
	EcologyManifest ecology_manifest.EcologyManifest
	Verbose         bool
}

func (lpc ListProjectCommand) Execute(o *output.Output) error {
	o.Info("Available Projects:").Indent()
	for project, projectPath := range lpc.EcologyManifest.ProjectManifestPaths {
		if lpc.Verbose {
			o.Success(project)
		} else {
			o.Success(project + " - " + projectPath)
		}
	}
	o.Dedent().Done()
	return nil
}
