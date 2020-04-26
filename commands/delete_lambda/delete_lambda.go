package delete_lambda

import (
	"github.com/gbdubs/ecology/manifests/ecology_manifest"
	"github.com/gbdubs/ecology/util/flag_validation"
	"github.com/gbdubs/ecology/util/output"
)

type DeleteLambdaCommand struct {
	EcologyManifest ecology_manifest.EcologyManifest
	Project         string
	Lambda          string
}

func (dlc DeleteLambdaCommand) Execute(o *output.Output) (err error) {
	em := &dlc.EcologyManifest
	pm, err := em.GetProjectManifest(dlc.Project)
	err = flag_validation.ValidateAll(
		flag_validation.Project(dlc.Project),
		flag_validation.ProjectExists(dlc.Project, em),
		flag_validation.Lambda(dlc.Lambda),
		flag_validation.LambdaExists(dlc.Lambda, pm),
		err)
	if err != nil {
		o.Error(err)
		return err
	}
	lm, err := pm.GetLambdaManifest(dlc.Lambda)

	o.Info("DeleteLambdaCommand - %s.DeleteFromPlatform", dlc.Lambda).Indent()
	err = lm.DeleteFromPlatform(o)
	if err != nil {
		o.Error(err)
		return
	}
	o.Dedent().Done()

	o.Info("DeleteLambdaCommand - %s.Save", dlc.Project).Indent()
	err = pm.RemoveLambdaManifest(lm)
	if err != nil {
		o.Error(err)
		return
	}
	err = pm.Save(o)
	if err != nil {
		o.Error(err)
		return
	}
	o.Dedent().Done()
	return nil
}
