package project_manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gbdubs/ecology/manifests/api_manifest"
	"github.com/gbdubs/ecology/manifests/lambda_manifest"
	"github.com/gbdubs/ecology/util/output"
	"io/ioutil"
	"os"
	"strings"
)

type ProjectConfigInfo struct {
	Name         string
	ManifestPath string
}

type ProjectDeployInfo struct {
	Platform string
	Region   string
}

type ProjectManifest struct {
	Config          ProjectConfigInfo
	Deploy          ProjectDeployInfo
	LambdaManifests []lambda_manifest.LambdaManifest
	ApiManifest     api_manifest.ApiManifest
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
		if l.Config.Name == lambdaName {
			return &pm.LambdaManifests[i], nil
		}
	}
	return nil, errors.New(fmt.Sprintf("No Lambda named %s in Project %s", lambdaName, pm.Config.Name))
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

func (pm *ProjectManifest) Save(o *output.Output) (err error) {
	o.Info("Writing Project Manifest to %s", pm.Config.ManifestPath).Indent()
	contents, _ := json.MarshalIndent(pm, "", "  ")
	filePath := pm.Config.ManifestPath
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
	o.Info("Pushing Project %s to Platform", pm.Config.Name).Indent()
	err = pm.pushLambdas(o)
	o.Dedent().Done()
	return
}

func (pm *ProjectManifest) pushLambdas(o *output.Output) (err error) {
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
	return
}

func (pm *ProjectManifest) DeleteFromPlatform(o *output.Output) (err error) {
	o.Info("Deleting Project %s", pm.Config.Name).Indent()
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
