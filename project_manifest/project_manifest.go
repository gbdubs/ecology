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
}

func GetProjectManifest(projectManifestPath string, o *output.Output) (projectManifest ProjectManifest, err error) {
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

func (projectManifest ProjectManifest) Save(o *output.Output) (err error) {
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

func (pm ProjectManifest) ToString() string {
	out, err := json.Marshal(pm)
	if err != nil {
		panic(err)
	}
	return string(out)
}
