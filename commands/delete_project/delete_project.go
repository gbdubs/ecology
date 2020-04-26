package delete_project

import (
	"github.com/gbdubs/ecology/manifests/ecology_manifest"
	"github.com/gbdubs/ecology/manifests/project_manifest"
	"github.com/gbdubs/ecology/util/output"
)

type DeleteProjectCommand struct {
	EcologyManifest ecology_manifest.EcologyManifest
	Project         string
}

func (dpc DeleteProjectCommand) Execute(o *output.Output) (err error) {
	o.Info("DeleteProjectCommand - Get Project Manifest %s", dpc.Project).Indent()
	pm, err := project_manifest.GetProjectManifestFromEcologyManifest(dpc.Project, &dpc.EcologyManifest, o)
	if err != nil {
		o.Error(err)
		return
	}
	o.Dedent().Done()

	o.Info("DeleteProjectCommand - %s.DeleteFromPlatform", dpc.Project).Indent()
	err = pm.DeleteFromPlatform(o)
	if err != nil {
		o.Error(err)
		return
	}
	o.Dedent().Done()

	o.Info("DeleteProjectCommand - %s.Save", dpc.Project).Indent()
	err = pm.Save(o)
	if err != nil {
		o.Error(err)
		return
	}
	o.Dedent().Done()
	return nil
}
