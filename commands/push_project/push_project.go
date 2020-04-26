package push_project

import (
	"github.com/gbdubs/ecology/manifests/ecology_manifest"
	"github.com/gbdubs/ecology/manifests/project_manifest"
	"github.com/gbdubs/ecology/util/output"
)

type PushProjectCommand struct {
	EcologyManifest ecology_manifest.EcologyManifest
	Project         string
}

func (plc PushProjectCommand) Execute(o *output.Output) (err error) {
	o.Info("PushProjectCommand - Get Project Manifest %s", plc.Project).Indent()
	pm, err := project_manifest.GetProjectManifestFromEcologyManifest(plc.Project, &plc.EcologyManifest, o)
	if err != nil {
		o.Error(err)
		return
	}
	o.Dedent().Done()

	o.Info("PushProjectCommand - %s.PushToPlatform", plc.Project).Indent()
	err = pm.PushToPlatform(o)
	if err != nil {
		o.Error(err)
		return
	}
	o.Dedent().Done()

	o.Info("PushProjectCommand - %s.Save", plc.Project).Indent()
	err = pm.Save(o)
	if err != nil {
		o.Error(err)
		return
	}
	o.Dedent().Done()
	return nil
}
