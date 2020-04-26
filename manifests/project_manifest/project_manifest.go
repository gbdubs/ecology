package project_manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gbdubs/ecology/manifests/lambda_manifest"
	"github.com/gbdubs/ecology/util/output"
	"io/ioutil"
	"os"
	"strings"
)

type ProjectManifest struct {
	ProjectManifestPath string
	ProjectName         string
	LambdaManifests     []lambda_manifest.LambdaManifest
}

func GetProjectManifestFromFile(projectManifestPath string) (projectManifest *ProjectManifest, err error) {
	data, err := ioutil.ReadFile(projectManifestPath)
	if err == nil {
		err = json.Unmarshal(data, &projectManifest)
	}
	return
}

func (pm *ProjectManifest) GetLambdaManifest(lambdaName string) (*lambda_manifest.LambdaManifest, error) {
	// TRICKSY POINTERSES! FILTHY TRICKSY POINTERSESSESSS!
	for i, l := range pm.LambdaManifests {
		if l.LambdaName == lambdaName {
			return &pm.LambdaManifests[i], nil
		}
	}
	return nil, errors.New(fmt.Sprintf("No Lambda named %s in Project %s", lambdaName, pm.ProjectName))
}

func (pm *ProjectManifest) RemoveLambdaManifest(ptr *lambda_manifest.LambdaManifest) error {
	indexToRemove := -1
	for i, _ := range pm.LambdaManifests {
		if &pm.LambdaManifests[i] == ptr {
			indexToRemove = i
		}
	}
	if indexToRemove == -1 {
		return errors.New("No Lambda with the given pointer was present")
	}
	pm.LambdaManifests = append(pm.LambdaManifests[:indexToRemove], pm.LambdaManifests[indexToRemove+1:]...)
	return nil
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
	o.Info("Pushing Lambdas").Indent()
	for i, _ := range pm.LambdaManifests {
		lm := &pm.LambdaManifests[i]
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

func (pm *ProjectManifest) DeleteFromPlatform(o *output.Output) (err error) {
	o.Info("Deleting Project %s", pm.ProjectName).Indent()
	o.Info("Deleting Lambdas").Indent()
	for i, _ := range pm.LambdaManifests {
		lm := &pm.LambdaManifests[i]
		err = lm.DeleteFromPlatform(o)
		if err != nil {
			o.Error(err)
			pm.Save(o) // Saves partial deletion progress in case we fail midway.
			return
		}
		err = pm.RemoveLambdaManifest(lm)
	}
	o.Dedent().Done()
	o.Dedent().Done()
	return
}
