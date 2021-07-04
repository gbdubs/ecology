package create_lambda

import (
	"github.com/gbdubs/ecology/manifests/ecology_manifest"
	"github.com/gbdubs/ecology/manifests/lambda_manifest"
	"github.com/gbdubs/ecology/util/flag_validation"
	"github.com/gbdubs/ecology/util/output"
	"strings"
	"io/ioutil"
	"os"
	"fmt"
)

type CreateLambdaCommand struct {
	EcologyManifest ecology_manifest.EcologyManifest
	Project         string
	Lambda          string
}

func (clc CreateLambdaCommand) Execute(o *output.Output) error {
	em := &clc.EcologyManifest
	pm, err := em.GetProjectManifest(clc.Project)
	err = flag_validation.ValidateAll(
		flag_validation.Project(clc.Project),
		flag_validation.ProjectExists(clc.Project, em),
		flag_validation.Lambda(clc.Lambda),
		flag_validation.LambdaDoesNotExist(clc.Lambda, pm),
		err)
	if err != nil {
		o.Error(err)
		return err
	}

	o.Info("CreateLambdaCommand - LambdaManifest.New").Indent()
	projectRootDir := pm.Config.ManifestPath[:strings.LastIndex(pm.Config.ManifestPath, "/")]
	lm, err := lambda_manifest.New(
		projectRootDir,
		clc.Project,
		clc.Lambda,
		pm.Deploy.Platform,
		pm.Deploy.Region,
		o)
	if err != nil {
		o.Error(err)
		return err
	}
	o.Dedent().Done()
	
	o.Info("CreateLambdaCommand - Create Initial Contents").Indent()
	err = os.MkdirAll(lm.Config.FolderPath, 0777)
	if err != nil {
	  o.Error(err)
	  return err
	}
	contents := fmt.Sprintf(initialLambdaFileContents, clc.Lambda, clc.Lambda, clc.Lambda)
  err = ioutil.WriteFile(lm.Config.CodePath, []byte(contents), 0777)
  if err != nil {
	  o.Error(err)
	  return err
	}
	o.Dedent().Done()

	o.Info("CreateLambdaCommand - %s.Save", clc.Project).Indent()
	pm.LambdaManifests = append(pm.LambdaManifests, *lm)
	err = pm.Save(o)
	if err != nil {
		o.Error(err)
		return err
	}
	o.Dedent().Done()
	return nil
}

const initialLambdaFileContents = `package main

import (
  "context"
  "github.com/aws/aws-lambda-go/lambda"
)

type %sRequest struct {
  Input string
}

func HandleRequest(ctx context.Context, request %sRequest) (string, error) {
  return "This is the lambda %s! request.Input=" + request.Input, nil
}

func main() {
  lambda.Start(HandleRequest)
}
`

