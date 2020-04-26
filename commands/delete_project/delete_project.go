package delete_project

import (
	"github.com/gbdubs/ecology/manifests/ecology_manifest"
	"github.com/gbdubs/ecology/util/flag_validation"
	"github.com/gbdubs/ecology/util/output"
)

type DeleteProjectCommand struct {
	EcologyManifest ecology_manifest.EcologyManifest
	Project         string
}

func (dpc DeleteProjectCommand) Execute(o *output.Output) (err error) {
	em := &dpc.EcologyManifest
	err = flag_validation.ValidateAll(
		flag_validation.Project(dpc.Project),
		flag_validation.ProjectExists(dpc.Project, em),
		err)
	if err != nil {
		o.Error(err)
		return err
	}
	pm, err := em.GetProjectManifest(dpc.Project)

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
