package ecology_manifest

import (
	"encoding/json"
	"github.com/gbdubs/ecology/manifests/project_manifest"
	"github.com/gbdubs/ecology/util/output"
	"io/ioutil"
	"os"
	"strings"
)

type EcologyManifest struct {
	ManifestPath         string
	ProjectManifestPaths map[string]string
}

const defaultEcologyManifestFilePath = "/Users/gradyward/.ecology/ecology.json"

func Get(o *output.Output) (ecologyManifest EcologyManifest, err error) {
	data, err := ioutil.ReadFile(defaultEcologyManifestFilePath)
	if err != nil {
		o.Warning("Ecology Manifest Not Found. Creating New Ecology Manifest.").Indent()
		ecologyManifest = EcologyManifest{
			ManifestPath:         defaultEcologyManifestFilePath,
			ProjectManifestPaths: make(map[string]string),
		}
		o.Dedent().Done()
		err = nil
	} else {
		err = json.Unmarshal(data, &ecologyManifest)
		if err != nil {
			return
		}
	}
	return ecologyManifest, err
}

func (em *EcologyManifest) Save(o *output.Output) (err error) {
	o.Info("Writing Ecology Manifest to %s...", em.ManifestPath).Indent()

	file, err := json.MarshalIndent(em, "", "  ")
	if err == nil {
		if strings.Index(em.ManifestPath, "/") > -1 {
			parentDir := em.ManifestPath[:strings.LastIndex(em.ManifestPath, "/")]
			err = os.MkdirAll(parentDir, 0777)
			if err != nil {
				o.Error(err)
				return
			}
		}
		err = ioutil.WriteFile(em.ManifestPath, file, 0777)
	}
	if err != nil {
		o.Error(err)
	}
	o.Dedent().Done()
	return
}

func (em *EcologyManifest) GetProjectManifest(project string) (*project_manifest.ProjectManifest, error) {
	return project_manifest.GetProjectManifestFromFile(em.ProjectManifestPaths[project])
}
