package create_project

import (
	"github.com/gbdubs/ecology/manifests/lambda_manifest"
	"github.com/gbdubs/ecology/manifests/ecology_manifest"
	"github.com/gbdubs/ecology/manifests/project_manifest"
	"github.com/gbdubs/ecology/util/flag_validation"
	"github.com/gbdubs/ecology/util/output"
)

type CreateProjectCommand struct {
	EcologyManifest ecology_manifest.EcologyManifest
	Platform        string
	Region          string
	Project         string
	Path            string
}

func (cpc *CreateProjectCommand) Execute(o *output.Output) (err error) {
	em := &cpc.EcologyManifest
	err = flag_validation.ValidateAll(
		flag_validation.Platform(cpc.Platform),
		flag_validation.Region(cpc.Region),
		flag_validation.Project(cpc.Project),
		flag_validation.ProjectDoesNotExist(cpc.Project, em),
		flag_validation.Path(cpc.Path))
	if err != nil {
		o.Error(err)
		return err
	}
	o.Info("CreateProjectCommand").Indent()
	manifest := project_manifest.ProjectManifest{
		Config: project_manifest.ProjectConfigInfo{
			Name:         cpc.Project,
			ManifestPath: cpc.Path + "/project.ecology.json",
		},
		Deploy: project_manifest.ProjectDeployInfo{
			Region:   cpc.Region,
			Platform: cpc.Platform,
		},
		LambdaManifests: make([]lambda_manifest.LambdaManifest, 0),
	}
	err = manifest.Save(o)
	if err != nil {
		o.Error(err)
		return err
	}
	em.ProjectManifestPaths[cpc.Project] = manifest.Config.ManifestPath
	em.Save(o)
	if err != nil {
		o.Error(err)
		return err
	}
	o.Dedent().Done()
	return
}
