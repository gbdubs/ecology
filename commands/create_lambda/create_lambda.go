package create_lambda

import (
	"github.com/gbdubs/ecology/manifests/ecology_manifest"
	"github.com/gbdubs/ecology/manifests/project_manifest"
	"github.com/gbdubs/ecology/manifests/lambda_manifest"
	"github.com/gbdubs/ecology/util/output"
	"strings"
)

type CreateLambdaCommand struct {
	EcologyManifest ecology_manifest.EcologyManifest
	Project         string
	Lambda          string
	Platform        string
	Region          string
}

func (clc CreateLambdaCommand) Execute(o *output.Output) (err error) {
	o.Info("Looking up Project Manifest for %s", clc.Project).Indent()
	pm, err := project_manifest.GetProjectManifestFromEcologyManifest(clc.Project, &clc.EcologyManifest, o)
	if err != nil {
	  o.Error(err)
	  return
	}
	o.Dedent().Done()
	
	o.Info("Creating Lambda").Indent()
	projectRootDir := pm.ProjectManifestPath[:strings.LastIndex(pm.ProjectManifestPath, "/")]
	lm, err := lambda_manifest.New(
	  projectRootDir, 
	  clc.Project,
	  clc.Lambda,
	  clc.Region,
	  o);
	if err != nil {
	  o.Error(err)
	  return
	}
	pm.LambdaManifests = append(pm.LambdaManifests, *lm)
	pm.Save(o)
	o.Dedent().Done()
	return nil
}
