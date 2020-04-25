package project_manifest

import (
	"encoding/json"
	"github.com/gbdubs/ecology/output"
	"io/ioutil"
	"os"
	"strings"
)

type ProjectManifest struct {
	ProjectManifestPath string
	ProjectName         string
	LambdaManifests     []lambda_manifest.LambdaManifest
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
	for _, lm := range pm.LambdaManifest {
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
