package create_project

import (
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
}

func (cpc *CreateProjectCommand) Execute(o *output.Output) (err error) {
	em := &cpc.EcologyManifest
	err = flag_validation.ValidateAll(
		flag_validation.Platform(cpc.Platform),
		flag_validation.Region(cpc.Region),
		flag_validation.Project(cpc.Project),
		flag_validation.ProjectDoesNotExist(cpc.Project, em))
	if err != nil {
		o.Error(err)
		return err
	}
	o.Info("CreateProjectCommand").Indent()
	manifest := project_manifest.ProjectManifest{
		ProjectName:         cpc.Project,
		ProjectManifestPath: em.EcologyProjectsDirectoryPath + "/" + cpc.Project + "/ecology.ecology",
	}

	err = manifest.Save(o)
	o.Dedent().Done()
	return
}
