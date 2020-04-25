package delete_project

import (
	"github.com/gbdubs/ecology/manifests/ecology_manifest"
	"github.com/gbdubs/ecology/util/output"
)

type DeleteProjectCommand struct {
	EcologyManifest ecology_manifest.EcologyManifest
	Project         string
}

func (dpc DeleteProjectCommand) Execute(o *output.Output) (err error) {
	o.Info("Deleting Project").Indent()
	o.Dedent().Done()
	return nil
}
