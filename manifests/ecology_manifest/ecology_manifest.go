package ecology_manifest

import (
	"encoding/json"
	"github.com/gbdubs/ecology/manifests/project_manifest"
	"github.com/gbdubs/ecology/util/output"
	"io/ioutil"
)

type EcologyManifest struct {
	EcologyManifestFilePath      string
	EcologyProjectsDirectoryPath string
}

const ecologyManifestFilePath = "/Users/gradyward/go/ecology/ecology.ecology"
const defaultEcologyProjectsDirectoryPath = "/Users/gradyward/go/src/ecology"

func GetEcologyManifest(o *output.Output) (ecologyManifest EcologyManifest, err error) {
	data, err := ioutil.ReadFile(ecologyManifestFilePath)
	if err != nil {
		o.Warning("Ecology Manifest Not Found. Creating New Ecology Manifest.").Indent()
		ecologyManifest = EcologyManifest{
			EcologyManifestFilePath:      ecologyManifestFilePath,
			EcologyProjectsDirectoryPath: defaultEcologyProjectsDirectoryPath,
		}
		o.Dedent().Done()
	} else {
		err = json.Unmarshal(data, &ecologyManifest)
		if err != nil {
			return
		}
	}
	ecologyManifest.Save(o)
	return ecologyManifest, err
}

func (ecologyManifest EcologyManifest) Save(o *output.Output) (err error) {
	o.Info("Writing Ecology Manifest to %s...", ecologyManifestFilePath).Indent()
	file, err := json.MarshalIndent(ecologyManifest, "", "  ")
	if err == nil {
		err = ioutil.WriteFile(ecologyManifestFilePath, file, 0777)
	}
	if err != nil {
		o.Error(err)
	}
	o.Dedent().Done()
	return
}

func (em *EcologyManifest) GetProjectManifestPath(project string) string {
	return em.EcologyProjectsDirectoryPath + "/" + project + "/ecology.ecology"
}

func (em *EcologyManifest) GetProjectManifest(project string) (*project_manifest.ProjectManifest, error) {
	return project_manifest.GetProjectManifestFromFile(em.GetProjectManifestPath(project))
}
