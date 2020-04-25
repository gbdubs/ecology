package push_project

import (
	"github.com/gbdubs/ecology/manifests/ecology_manifest"
	"github.com/gbdubs/ecology/util/output"
)

type PushProjectCommand struct {
	EcologyManifest ecology_manifest.EcologyManifest
	Project         string
}

func (plc PushProjectCommand) Execute(o *output.Output) (err error) {
	o.Info("Pushing Lambda").Indent()
	o.Dedent().Done()
	return nil
}
