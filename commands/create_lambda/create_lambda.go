package create_lambda

import (
	"github.com/gbdubs/ecology/manifests/ecology_manifest"
	"github.com/gbdubs/ecology/manifests/lambda_manifest"
	"github.com/gbdubs/ecology/util/flag_validation"
	"github.com/gbdubs/ecology/util/output"
	"strings"
)

type CreateLambdaCommand struct {
	EcologyManifest ecology_manifest.EcologyManifest
	Platform        string
	Region          string
	Project         string
	Lambda          string
}

func (clc CreateLambdaCommand) Execute(o *output.Output) error {
	em := &clc.EcologyManifest
	pm, err := em.GetProjectManifest(clc.Project)
	err = flag_validation.ValidateAll(
		flag_validation.Platform(clc.Platform),
		flag_validation.Region(clc.Region),
		flag_validation.Project(clc.Project),
		flag_validation.ProjectExists(clc.Project, em),
		flag_validation.Lambda(clc.Lambda),
		flag_validation.LambdaDoesNotExist(clc.Lambda, pm),
		err)
	if err != nil {
		o.Error(err)
		return err
	}

	o.Info("CreateLambdaCommand - LambdaManifest.New").Indent()
	projectRootDir := pm.ProjectManifestPath[:strings.LastIndex(pm.ProjectManifestPath, "/")]
	lm, err := lambda_manifest.New(
		projectRootDir,
		clc.Project,
		clc.Lambda,
		clc.Region,
		o)
	if err != nil {
		o.Error(err)
		return err
	}
	o.Dedent().Done()

	o.Info("CreateLambdaCommand - %s.Save", clc.Project).Indent()
	pm.LambdaManifests = append(pm.LambdaManifests, *lm)
	err = pm.Save(o)
	if err != nil {
		o.Error(err)
		return err
	}
	o.Dedent().Done()
	return nil
}
