package create_project

import (
	"encoding/json"
	"errors"
	"github.com/gbdubs/ecology/manifests/ecology_manifest"
	"github.com/gbdubs/ecology/manifests/project_manifest"
	"github.com/gbdubs/ecology/util/output"
	"regexp"
	"strings"
)

type CreateProjectCommand struct {
	EcologyManifest     ecology_manifest.EcologyManifest
	ProjectPath         string
	ProjectSimpleName   string
	ProjectManifestPath string
	Platform            string
	Region              string
}

func (cp CreateProjectCommand) Execute(o *output.Output) (err error) {
	o.Info("Inferring Unspecified Values for CreateProjectCommand").Indent()
	cp = applyConstantTransformations(cp)
	o.Dedent().Done()
	o.Info("Validating CreateProjectComand").Indent()
	err = cp.validate()
	if err != nil {
		o.Error(err)
		return
	}
	o.Dedent().Done()

	o.Info("Running CreateProjectCommand").Indent()

	manifest := project_manifest.ProjectManifest{
		ProjectManifestPath: cp.ProjectManifestPath,
		ProjectName:         cp.ProjectSimpleName,
	}
	err = manifest.Save(o)
	o.Dedent().Done()
	return
}

const alphanumericRegex = "^[a-zA-Z0-9]*$"
const alphanumericWithSlashesRegex = "^[a-zA-Z0-9/]*$"

// Whether the given platform is currently supported
var platforms = map[string]bool{"GCP": false, "AWS": true}

func applyConstantTransformations(cp CreateProjectCommand) CreateProjectCommand {
	if cp.ProjectSimpleName == "" && cp.ProjectPath != "" {
		splits := strings.Split(cp.ProjectPath, "/")
		cp.ProjectSimpleName = splits[len(splits)-1]
	} else if cp.ProjectPath == "" && cp.ProjectSimpleName != "" {
		cp.ProjectPath = cp.ProjectSimpleName
	}
	cp.ProjectManifestPath = cp.EcologyManifest.EcologyProjectsDirectoryPath + "/" + cp.ProjectPath + "/ecology.ecology"
	return cp
}

func (cp CreateProjectCommand) validate() (err error) {
	if cp.ProjectSimpleName == "" && cp.ProjectPath == "" {
		return errors.New("create_project requires --project_name or --project_path")
	}
	match, _ := regexp.MatchString(alphanumericRegex, cp.ProjectSimpleName)
	if !match {
		return errors.New("--project_name can only contain alphanumeric characters")
	}
	match, _ = regexp.MatchString(alphanumericWithSlashesRegex, cp.ProjectPath)
	if !match {
		return errors.New("--project_path can only contain alphanumeric characters or slashes")
	}

	_, isEnumeratedPlatform := platforms[cp.Platform]
	if !isEnumeratedPlatform {
		return errors.New("--platform should be one of AWS or GCP")
	}
	if !platforms[cp.Platform] {
		return errors.New("--platform=GCP is not yet supported for this command")
	}

	for _, otherProjectManifestPath := range cp.EcologyManifest.ProjectManifestPaths {
		if otherProjectManifestPath == cp.ProjectManifestPath {
			return errors.New("Already a project at " + cp.ProjectManifestPath)
		}
	}

	return nil
}

func (cp CreateProjectCommand) ToString() string {
	out, err := json.Marshal(cp)
	if err != nil {
		panic(err)
	}
	return string(out)
}
