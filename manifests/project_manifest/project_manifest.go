package project_manifest

import (
	"encoding/json"
	"github.com/gbdubs/ecology/manifests/lambda_manifest"
	"github.com/gbdubs/ecology/manifests/ecology_manifest"
	"github.com/gbdubs/ecology/util/output"
	"io/ioutil"
	"os"
	"fmt"
	"strings"
	"errors"
)

type ProjectManifest struct {
	ProjectManifestPath string
	ProjectName         string
	LambdaManifests     []lambda_manifest.LambdaManifest
}

func GetProjectManifestFromEcologyManifest(project string, em *ecology_manifest.EcologyManifest, o *output.Output) (projectManifest *ProjectManifest, err error) {
  projectManifestPath := em.EcologyProjectsDirectoryPath + "/" + project + "/ecology.ecology";
  return GetProjectManifest(projectManifestPath, o)
}

func GetProjectManifest(projectManifestPath string, o *output.Output) (projectManifest *ProjectManifest, err error) {
	o.Info("Reading Project Manifest from %s...", projectManifestPath).Indent()
	data, err := ioutil.ReadFile(projectManifestPath)
	if err == nil {
		err = json.Unmarshal(data, &projectManifest)
	}
	if err != nil {
		o.Error(err)
	}
	o.Dedent().Done()
	return
}

func (pm *ProjectManifest) GetLambdaManifest(lambdaName string, o *output.Output) (*lambda_manifest.LambdaManifest, error) {
  // TRICKSY POINTERSES! FILTHY TRICKSY POINTERSESSESSS!
  for i, l := range pm.LambdaManifests {
    if l.LambdaName == lambdaName {
      return &pm.LambdaManifests[i], nil
    }
  }
  return nil, errors.New(fmt.Sprintf("No Lambda named %s in Project %s", lambdaName, pm.ProjectName))
}


func (projectManifest *ProjectManifest) Save(o *output.Output) (err error) {
	o.Info("Writing Project Manifest to %s", projectManifest.ProjectManifestPath).Indent()
	contents, _ := json.MarshalIndent(projectManifest, "", "  ")
	filePath := projectManifest.ProjectManifestPath
	if strings.Index(filePath, "/") > -1 {
		parentDir := filePath[:strings.LastIndex(filePath, "/")]
		err = os.MkdirAll(parentDir, 0777)
	}
	if err == nil {
		err = ioutil.WriteFile(filePath, contents, 0777)
	}
	if err != nil {
		o.Error(err)
	}
	o.Dedent().Done()
	return
}

func (pm *ProjectManifest) PushToPlatform(o *output.Output) (err error) {
	o.Info("Pushing Project %s to Platform", pm.ProjectName).Indent()
	o.Info("Pushing Lambdas")
	for _, lm := range pm.LambdaManifests {
		err = lm.PushToPlatform(o)
		if err != nil {
			o.Error(err)
			return err
		}
	}
	o.Dedent().Done()
	o.Dedent().Done()
	return
}
