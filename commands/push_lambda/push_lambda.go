package push_lambda

import (
	"github.com/gbdubs/ecology/manifests/ecology_manifest"
	"github.com/gbdubs/ecology/util/flag_validation"
	"github.com/gbdubs/ecology/util/output"
)

type PushLambdaCommand struct {
	EcologyManifest ecology_manifest.EcologyManifest
	Project         string
	Lambda          string
}

func (plc PushLambdaCommand) Execute(o *output.Output) (err error) {
	em := &plc.EcologyManifest
	pm, err := em.GetProjectManifest(plc.Project)
	err = flag_validation.ValidateAll(
		flag_validation.Project(plc.Project),
		flag_validation.ProjectExists(plc.Project, em),
		flag_validation.Lambda(plc.Lambda),
		flag_validation.LambdaExists(plc.Lambda, pm),
		err)
	if err != nil {
		o.Error(err)
		return err
	}
	lm, err := pm.GetLambdaManifest(plc.Lambda)

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
