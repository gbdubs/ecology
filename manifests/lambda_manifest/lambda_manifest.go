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
	err = lm.writeInitialFile(o)
	return
}

func (lm *LambdaManifest) writeInitialFile(o *output.Output) (err error) {
	filePath := lm.LambdaCodePath
	err = os.MkdirAll(lm.LambdaCodeFolderPath, 0777)
	if err == nil {
		err = ioutil.WriteFile(filePath, []byte(
			fmt.Sprintf(initialLambdaFileContents, lm.LambdaName, lm.LambdaName, lm.FullyQualifiedLambdaName)), 0777)
	}
	return
}

func (lm *LambdaManifest) packageToDeploy(o *output.Output) (err error) {
	o.Info("LambdaManifest - packageToDeploy - Build").Indent()
	lambdaFolder := lm.LambdaCodePath[:strings.LastIndex(lm.LambdaCodePath, "/")]
	buildArgs := strings.Split(fmt.Sprintf("GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o %s/%s %s/%s.go", lambdaFolder, lm.FullyQualifiedLambdaName, lambdaFolder, lm.LambdaName), " ")

	ctx, cancelBuild := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelBuild()

	if result, err := exec.CommandContext(ctx, "env", buildArgs...).CombinedOutput(); err != nil {
		o.Failure(string(result))
		return err
	}
	o.Dedent().Done()
  
  o.Info("LambdaManifest - packageToDeploy - Zip").Indent()
	zipArgs := strings.Split(fmt.Sprintf("-j %s/%s.zip %s/%s", lambdaFolder, lm.FullyQualifiedLambdaName, lambdaFolder, lm.FullyQualifiedLambdaName), " ")
	ctx, cancelZip := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelZip()
	if result, err := exec.CommandContext(ctx, "zip", zipArgs...).CombinedOutput(); err != nil {
		o.Failure(string(result))
		return err
	}
	o.Dedent().Done()
	return
}

func (lm *LambdaManifest) PushToPlatform(o *output.Output) (err error) {
	lambdaFolder := lm.LambdaCodePath[:strings.LastIndex(lm.LambdaCodePath, "/")]

	o.Info("LambdaManifest - %s.PushToPlatform", lm.LambdaName).Indent()

	err = lm.ExecutorRoleManifest.PushToPlatform(o)
	if err != nil {
		o.Error(err)
		return
	}

	var currentCodeHash string
	if lm.ExistsOnPlatform {
		o.Info("LambdaManifest - PushToPlatform - Is Push Necessary?").Indent()
		currentCodeHash, err = file_hash.ComputeFileHash(lm.LambdaCodePath, o)
		o.Info("Old Hash: %s", lm.LambdaCodeLastPushedHash)
		o.Info("Current Hash: %s", currentCodeHash)
		if err != nil {
			return err
		}
		if currentCodeHash == lm.LambdaCodeLastPushedHash {
			o.Dedent().Success("Code hasn't changed since last push, no push needed.").Dedent().Done()
			return nil
		}
		o.Dedent().Warning("Code has changed since last push.")
	} else {
		o.Info("LambdaManifest - PushToPlatform - First Lambda Push")
	}

  
	err = lm.packageToDeploy(o)
	if err != nil {
	  return err
	}

	o.Info("LambdaManifest - PushToPlatform - Read Zip").Indent()
	zipPath := fmt.Sprintf("%s/%s.zip", lambdaFolder, lm.FullyQualifiedLambdaName)
	zipBytes, err := ioutil.ReadFile(zipPath)
	if err != nil {
		return
	}
	o.Dedent().Done()

	svc := lambda.New(session.New(), aws.NewConfig().WithRegion(lm.Region))

	o.Info("LambdaManifest - PushToPlatform - Check If Lambda Exists").Indent()
	_, err = svc.GetFunction(&lambda.GetFunctionInput{
		FunctionName: aws.String(lm.FullyQualifiedLambdaName),
	})
	if err == nil {
		o.Warning("Lambda Already Exists.").Dedent().Done()
	  o.Info("LambdaManifest - PushToPlatform - Update Lambda").Indent()
		updateFunctionRequest := &lambda.UpdateFunctionCodeInput{
			ZipFile:      zipBytes,
			FunctionName: aws.String(lm.FullyQualifiedLambdaName),
			Publish:      aws.Bool(true),
		}
		_, err = svc.UpdateFunctionCode(updateFunctionRequest)
		if err != nil {
			return err
		}
		o.Dedent().Done()
	} else {
		o.Warning("Lambda Does Not Exist.").Dedent().Done()
	  o.Info("LambdaManifest - PushToPlatform - Create Lambda").Indent()
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
			return err
		}
		o.Dedent().Done()
	}
	lm.LambdaCodeLastPushedHash = currentCodeHash
	lm.ExistsOnPlatform = true
	return nil
}

func (lm *LambdaManifest) DeleteFromPlatform(o *output.Output) (err error) {
	o.Info("LambdaManifest - DeleteFromPlatform - %s", lm.FullyQualifiedLambdaName).Indent()

	err = lm.ExecutorRoleManifest.DeleteFromPlatform(o)
	if err != nil {
		o.Error(err)
		return
	}

	deleteFunctionRequest := &lambda.DeleteFunctionInput{
		FunctionName: aws.String(lm.FullyQualifiedLambdaName),
	}
	svc := lambda.New(session.New(), aws.NewConfig().WithRegion(lm.Region))
	_, err = svc.DeleteFunction(deleteFunctionRequest)
	if err != nil {
		o.Error(err)
		return err
	}
	lm.ExistsOnPlatform = false
	lm.LambdaCodeLastPushedHash = ""
	return nil
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
