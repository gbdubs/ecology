package flag_validation

import (
	"errors"
	"fmt"
	"github.com/gbdubs/ecology/manifests/ecology_manifest"
	"github.com/gbdubs/ecology/manifests/project_manifest"
	"regexp"
)

const alphanumericRegex = "^[a-zA-Z0-9]*$"
const alphanumericWithSlashesRegex = "^[a-zA-Z0-9/]*$"

// Whether the given platform is currently supported
var platforms = map[string]bool{"GCP": false, "AWS": true}

var knownRegions = []string{"us-west-1", "us-west-2", "us-east-1", "us-east-2"}

func ValidateAll(errs ...error) error {
	nonNilErrs := []error{}
	for i := range errs {
		if errs[i] != nil {
			nonNilErrs = append(nonNilErrs, errs[i])
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
	match, _ := regexp.MatchString(alphanumericRegex, project)
	if !match {
		return errors.New("--project can only contain alphanumeric characters")
	}
	return nil
}

func ProjectExists(project string, em *ecology_manifest.EcologyManifest) error {
	if !projectExists(project, em) {
		return errors.New(fmt.Sprintf("--project=%s doesn't exist", project))
	}
	return nil
}

func ProjectDoesNotExist(project string, em *ecology_manifest.EcologyManifest) error {
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
	match, _ := regexp.MatchString(alphanumericRegex, lambda)
	if !match {
		return errors.New("--lambda can only contain alphanumeric characters")
	}
	return nil
}

func LambdaExists(lambda string, pm *project_manifest.ProjectManifest) error {
	if pm == nil {
		return errors.New("Couldn't find a Project Manifest")
	}
	if !lambdaExists(lambda, pm) {
		return errors.New(fmt.Sprintf("--lambda=%s doesn't exist", lambda))
	}
	return nil
}

func LambdaDoesNotExist(lambda string, pm *project_manifest.ProjectManifest) error {
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
