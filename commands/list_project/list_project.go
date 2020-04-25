package list_project

import (
	"github.com/gbdubs/ecology/ecology_manifest"
	"github.com/gbdubs/ecology/output"
	"strings"
)

type ListProjectCommand struct {
	EcologyManifest     ecology_manifest.EcologyManifest
	Verbose bool
}

func (lpc ListProjectCommand) Execute(o *output.Output) (err error) {
	o.Info("Available Projects:").Indent()
	for _, manifestPath := range lpc.EcologyManifest.ProjectManifestPaths {
	  if lpc.Verbose {
	    o.Success(manifestPath); 
	  } else {
	    o.Success(
	      strings.Replace(strings.Replace(manifestPath, lpc.EcologyManifest.EcologyProjectsDirectoryPath + "/", "", 1), "/ecology.ecology", "", 1)) 
	  }
	}
	o.Dedent().Done();
	return nil
}