package lambda_manifest

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/gbdubs/ecology/manifests/role_manifest"
	"github.com/gbdubs/ecology/util/file_hash"
	"github.com/gbdubs/ecology/util/output"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

type LambdaManifest struct {
	LambdaName               string
	LambdaNameLowercase      string
	FullyQualifiedLambdaName string
	LambdaCodeFolderPath     string
	LambdaCodePath           string
	LambdaCodeLastPushedHash string
	Region                   string
	ExistsOnPlatform         bool
	ExecutorRoleManifest     role_manifest.RoleManifest
}

func New(projectDir string, projectName string, lambdaName string, region string, o *output.Output) (lm *LambdaManifest, err error) {
	fullyQualifiedLambdaName := projectName + "-" + lambdaName
	erm := role_manifest.New(fullyQualifiedLambdaName + "-executor")
	lambdaNameLower := strings.ToLower(lambdaName)
	lambdaCodeFolderPath := fmt.Sprintf("%s/lambda/%s", projectDir, lambdaNameLower)
	lambdaCodePath := fmt.Sprintf("%s/%s.go", lambdaCodeFolderPath, lambdaName)
	lambdaManifest := LambdaManifest{
		LambdaName:               lambdaName,
		LambdaNameLowercase:      lambdaNameLower,
		FullyQualifiedLambdaName: fullyQualifiedLambdaName,
		LambdaCodeFolderPath:     lambdaCodeFolderPath,
		LambdaCodePath:           lambdaCodePath,
		LambdaCodeLastPushedHash: "",
		Region:                   region,
		ExistsOnPlatform:         false,
		ExecutorRoleManifest:     erm,
	}
	lm = &lambdaManifest
	return
}

func (lm *LambdaManifest) WriteInitialFile(o *output.Output) (err error) {
	o.Info("Writing first version of %s to disk", lm.FullyQualifiedLambdaName)
	filePath := lm.LambdaCodePath
	err = os.MkdirAll(lm.LambdaCodeFolderPath, 0777)
	if err == nil {
		err = ioutil.WriteFile(filePath, []byte(
			fmt.Sprintf(initialLambdaFileContents, lm.LambdaName, lm.LambdaName, lm.FullyQualifiedLambdaName)), 0777)
	}
	if err != nil {
		o.Error(err)
	} else {
		o.Dedent().Done()
	}
	return
}

func (lm *LambdaManifest) PackageToDeploy(o *output.Output) (err error) {
	o.Info("Building Lambda").Indent()
	lambdaFolder := lm.LambdaCodePath[:strings.LastIndex(lm.LambdaCodePath, "/")]
	buildArgs := strings.Split(fmt.Sprintf("GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o %s/%s %s/%s.go", lambdaFolder, lm.FullyQualifiedLambdaName, lambdaFolder, lm.LambdaName), " ")

	ctx, cancelBuild := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelBuild()

	if result, err := exec.CommandContext(ctx, "env", buildArgs...).CombinedOutput(); err != nil {
		o.Failure(string(result))
		o.Error(err)
		return err
	}
	o.Dedent().Done()

	o.Info("Zipping Result").Indent()
	zipArgs := strings.Split(fmt.Sprintf("-j %s/%s.zip %s/%s", lambdaFolder, lm.FullyQualifiedLambdaName, lambdaFolder, lm.FullyQualifiedLambdaName), " ")
	ctx, cancelZip := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelZip()
	if result, err := exec.CommandContext(ctx, "zip", zipArgs...).CombinedOutput(); err != nil {
		o.Failure(string(result))
		o.Error(err)
		return err
	}
	o.Dedent().Done()
	return
}

func (lm *LambdaManifest) PushToPlatform(o *output.Output) (err error) {
	lambdaFolder := lm.LambdaCodePath[:strings.LastIndex(lm.LambdaCodePath, "/")]

	o.Info("Pushing Lambda %s to Platform", lm.FullyQualifiedLambdaName).Indent()

	err = lm.ExecutorRoleManifest.PushToPlatform(o)
	if err != nil {
		o.Error(err)
		return
	}

	var currentCodeHash string
	if lm.ExistsOnPlatform {
		o.Info("Verifying Lambda is Up To Date...")
		currentCodeHash, err = file_hash.ComputeFileHash(lm.LambdaCodePath, o)
		if err != nil {
			return err
		}
		if currentCodeHash == lm.LambdaCodeLastPushedHash {
			o.Info("Code hasn't changed since last push, no push needed.'")
			return nil
		}
		o.Info("Code has changed since last push.")
	} else {
		o.Info("Lambda has never been deployed.")
	}

	lm.PackageToDeploy(o)

	o.Info("Reading Deployable Lambda from Disk...").Indent()
	zipPath := fmt.Sprintf("%s/%s.zip", lambdaFolder, lm.FullyQualifiedLambdaName)
	zipBytes, err := ioutil.ReadFile(zipPath)
	if err != nil {
		o.Error(err)
		return
	}
	o.Dedent().Done()

	svc := lambda.New(session.New(), aws.NewConfig().WithRegion(lm.Region))

	o.Info("Creating Lambda on AWS...").Indent()
	createFunctionRequest := &lambda.CreateFunctionInput{
		Code: &lambda.FunctionCode{
			ZipFile: zipBytes,
		},
		Description:  aws.String(fmt.Sprintf("Ecology-Generated Lambda %s.", lm.FullyQualifiedLambdaName)),
		FunctionName: aws.String(lm.FullyQualifiedLambdaName),
		Handler:      aws.String(lm.FullyQualifiedLambdaName),
		Publish:      aws.Bool(true),
		Role:         aws.String(lm.ExecutorRoleManifest.Arn),
		Runtime:      aws.String("go1.x"),
	}
	_, err = svc.CreateFunction(createFunctionRequest)
	if err != nil {
		o.Error(err)
		return err
	}
	o.Dedent().Done()

	lm.LambdaCodeLastPushedHash = currentCodeHash
	lm.ExistsOnPlatform = true

	return
}

const initialLambdaFileContents = `
package main

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
