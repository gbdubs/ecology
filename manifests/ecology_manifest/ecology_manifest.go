package ecology_manifest

import (
	"encoding/json"
	"github.com/gbdubs/ecology/output"
	"io/ioutil"
	"os"
	"path/filepath"
)

type EcologyManifest struct {
	EcologyManifestFilePath      string
	EcologyProjectsDirectoryPath string
	InvokingWorkingDirectoryPath string
	ProjectManifestPaths         []string
}

const ecologyManifestFilePath = "/Users/gradyward/go/ecology/ecology.ecology"
const defaultEcologyProjectsDirectoryPath = "/Users/gradyward/go/src/ecology"

func GetEcologyManifest(o *output.Output) (ecologyManifest EcologyManifest, err error) {
	o.Info("Reading Ecology Manifest from %s...", ecologyManifestFilePath)
	data, err := ioutil.ReadFile(ecologyManifestFilePath)
	if err != nil {
		o.Indent().Warning("Ecology Manifest Not Found...").Indent()
		ecologyManifest = createEcologyManifest(o)
		o.Dedent().Done().Dedent()
	} else {
		err = json.Unmarshal(data, &ecologyManifest)
		if err != nil {
			return
		}
	}
	o.Done()
	o.Info("Updating Manifest with data from disk...").Indent()
	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	ecologyManifest.InvokingWorkingDirectoryPath = pwd
	ecologyManifest.ProjectManifestPaths, err = listAllEcologyManifestsInDirectory(ecologyManifest.EcologyProjectsDirectoryPath)
	o.Dedent().Done()
	ecologyManifest.Save(o)
	return ecologyManifest, err
}

func createEcologyManifest(o *output.Output) EcologyManifest {
	o.Info("Creating New Ecology Manifest at %s", ecologyManifestFilePath)
	return EcologyManifest{
		EcologyManifestFilePath:      ecologyManifestFilePath,
		EcologyProjectsDirectoryPath: defaultEcologyProjectsDirectoryPath,
		InvokingWorkingDirectoryPath: "",  // Will be filled out by GetEcologyManifest
		ProjectManifestPaths:         nil, // Will be filled out by GetEcologyManifest
	}
}

func listAllEcologyManifestsInDirectory(root string) (files []string, err error) {
	err = filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == ".ecology" {
				files = append(files, path)
			}
			return nil
		})
	return files, err
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

func (cp EcologyManifest) ToString() string {
	out, err := json.Marshal(cp)
	if err != nil {
		panic(err)
	}
	return string(out)
}
