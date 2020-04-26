package push_project

import (
	"github.com/gbdubs/ecology/manifests/ecology_manifest"
	"github.com/gbdubs/ecology/util/flag_validation"
	"github.com/gbdubs/ecology/util/output"
)

type PushProjectCommand struct {
	EcologyManifest ecology_manifest.EcologyManifest
	Project         string
}

func (ppc PushProjectCommand) Execute(o *output.Output) (err error) {
	em := &ppc.EcologyManifest
	err = flag_validation.ValidateAll(
		flag_validation.Project(ppc.Project),
		flag_validation.ProjectExists(ppc.Project, em),
		err)
	if err != nil {
		o.Error(err)
		return err
	}
	pm, err := em.GetProjectManifest(ppc.Project)

	o.Info("PushProjectCommand - %s.PushToPlatform", ppc.Project).Indent()
	err = pm.PushToPlatform(o)
	if err != nil {
		o.Error(err)
		return
	}
	o.Dedent().Done()

	o.Info("PushProjectCommand - %s.Save", ppc.Project).Indent()
	err = pm.Save(o)
	if err != nil {
		o.Error(err)
		return
	}
	o.Dedent().Done()
	return nil
}
