package flag_validation

import (
	"errors"
	"fmt"
	"github.com/gbdubs/ecology/manifests/ecology_manifest"
	"github.com/gbdubs/ecology/manifests/project_manifest"
	"io/ioutil"
	"regexp"
)

const alphanumericRegex = "^[a-zA-Z0-9]+$"
const alphanumericWithSlashesRegex = "^[a-zA-Z0-9/]+$"

// Whether the given platform is currently supported
var platforms = map[string]bool{"GCP": false, "AWS": true}

var knownRegions = []string{"us-west-1", "us-west-2", "us-east-1", "us-east-2"}

func ValidateAll(errs ...error) error {
	nonNilErrs := []error{}
	for _, e := range errs {
		if e != nil {
			nonNilErrs = append(nonNilErrs, e)
		}
	}
	if len(nonNilErrs) == 0 {
		return nil
	}
	resultErrStr := ""
	for _, err := range nonNilErrs {
		resultErrStr = resultErrStr + "\n" + fmt.Sprintf("%v", err)
	}
	return errors.New(resultErrStr)
}

func Platform(platform string) error {
	if platform == "" {
		return errors.New("Must set --platform")
	}
	_, isEnumeratedPlatform := platforms[platform]
	if !isEnumeratedPlatform {
		return errors.New("--platform should be one of AWS or GCP")
	}
	if !platforms[platform] {
		return errors.New("--platform=GCP is not yet supported")
	}
	return nil
}

func Region(region string) error {
	if region == "" {
		return errors.New("Must set --region")
	}
	found := false
	for _, r := range knownRegions {
		if r == region {
			found = true
		}
	}
	if !found {
		return errors.New("--region was not recognized")
	}
	return nil
}

func Project(project string) error {
	if project == "" {
		return errors.New("Must set --project")
	}
	match, _ := regexp.MatchString(alphanumericRegex, project)
	if !match {
		return errors.New("--project can only contain alphanumeric characters")
	}
	return nil
}

func Path(path string) error {
	if path == "" {
		return errors.New("Must set --path")
	}
	_, err := ioutil.ReadDir(path)
	if err == nil {
		return errors.New("Folder already exists: " + path)
	}
	return nil
}

func ProjectExists(project string, em *ecology_manifest.EcologyManifest) error {
  if Project(project) != nil {
    return nil
  }
	if !projectExists(project, em) {
		return errors.New(fmt.Sprintf("--project=%s doesn't exist", project))
	}
	return nil
}

func ProjectDoesNotExist(project string, em *ecology_manifest.EcologyManifest) error {
  if Project(project) != nil {
    return nil
  }
	if projectExists(project, em) {
		return errors.New(fmt.Sprintf("--project=%s doesn't exist", project))
	}
	return nil
}

func projectExists(project string, em *ecology_manifest.EcologyManifest) bool {
	_, err := em.GetProjectManifest(project)
	return err == nil
}

func Lambda(lambda string) error {
	if lambda == "" {
		return errors.New("Must set --lambda")
	}
	match, _ := regexp.MatchString(alphanumericRegex, lambda)
	if !match {
		return errors.New("--lambda can only contain alphanumeric characters")
	}
	return nil
}

func LambdaExists(lambda string, pm *project_manifest.ProjectManifest) error {
if     Lambda(lambda) != nil {
    return nil
  }
	if pm == nil {
		return errors.New("Couldn't find a Project Manifest")
	}
	if !lambdaExists(lambda, pm) {
		return errors.New(fmt.Sprintf("--lambda=%s doesn't exist", lambda))
	}
	return nil
}

func LambdaDoesNotExist(lambda string, pm *project_manifest.ProjectManifest) error {
  if     Lambda(lambda) != nil {
    return nil
  }
	if pm == nil {
		return errors.New("Couldn't find a Project Manifest")
	}
	if lambdaExists(lambda, pm) {
		return errors.New(fmt.Sprintf("--lambda=%s already exists", lambda))
	}
	return nil
}

func lambdaExists(lambda string, pm *project_manifest.ProjectManifest) bool {
	_, err := pm.GetLambdaManifest(lambda)
	return err == nil
}
