package push_lambda

import (
	"github.com/gbdubs/ecology/manifests/ecology_manifest"
	"github.com/gbdubs/ecology/manifests/project_manifest"
	"github.com/gbdubs/ecology/util/output"
)

type PushLambdaCommand struct {
	EcologyManifest ecology_manifest.EcologyManifest
	Project         string
	Lambda          string
}

func (plc PushLambdaCommand) Execute(o *output.Output) (err error) {
	o.Info("PushLambdaCommand - Get Project Manifest %s", plc.Project).Indent()
	pm, err := project_manifest.GetProjectManifestFromEcologyManifest(plc.Project, &plc.EcologyManifest, o)
	if err != nil {
		o.Error(err)
		return
	}
	o.Dedent().Done()

	o.Info("PushLambdaCommand - %s.GetLambdaManifest", plc.Project).Indent()
	lm, err := pm.GetLambdaManifest(plc.Lambda, o)
	if err != nil {
		o.Error(err)
		return
	}
	o.Dedent().Done()

	o.Info("PushLambdaCommand - %s.PushToPlatform", plc.Lambda).Indent()
	err = lm.PushToPlatform(o)
	if err != nil {
		o.Error(err)
		return
	}
	o.Dedent().Done()

	o.Info("PushLambdaCommand - %s.Save", plc.Project).Indent()
	err = pm.Save(o)
	if err != nil {
		o.Error(err)
		return
	}
	o.Dedent().Done()
	return nil
}
