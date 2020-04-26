package delete_lambda

import (
	"github.com/gbdubs/ecology/manifests/ecology_manifest"
	"github.com/gbdubs/ecology/manifests/project_manifest"
	"github.com/gbdubs/ecology/util/output"
)

type DeleteLambdaCommand struct {
	EcologyManifest ecology_manifest.EcologyManifest
	Project         string
	Lambda          string
}

func (dlc DeleteLambdaCommand) Execute(o *output.Output) (err error) {
	o.Info("DeleteLambdaCommand - Get Project Manifest %s", dlc.Project).Indent()
	pm, err := project_manifest.GetProjectManifestFromEcologyManifest(dlc.Project, &dlc.EcologyManifest, o)
	if err != nil {
		o.Error(err)
		return
	}
	o.Dedent().Done()

	o.Info("DeleteLambdaCommand - Get Lambda Manifest %s", dlc.Lambda).Indent()
	lm, err := pm.GetLambdaManifest(dlc.Lambda, o)
	if err != nil {
		o.Error(err)
		return
	}
	o.Dedent().Done()

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
