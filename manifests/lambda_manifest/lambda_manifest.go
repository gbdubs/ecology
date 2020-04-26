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

type LambdaCodeInfo struct {
	FolderPath string
	CodePath   string
	BuiltPath  string
	ZippedPath string
}

type LambdaDeployInfo struct {
	Region           string
	LastDeployedHash string
	Arn              string
}

type LambdaManifest struct {
	LambdaName               string
	FullyQualifiedLambdaName string
	CodeInfo                 LambdaCodeInfo
	DeployInfo               LambdaDeployInfo
	ExecutorRoleManifest     role_manifest.RoleManifest
}

func New(projectDir string, projectName string, lambdaName string, region string, o *output.Output) (lm *LambdaManifest, err error) {
	fullyQualifiedLambdaName := projectName + "-" + lambdaName
	erm := role_manifest.New(fullyQualifiedLambdaName + "-executor")
	codeInfoFolderPath := fmt.Sprintf("%s/lambda/%s", projectDir, lambdaName)
	codeInfoCodePath := fmt.Sprintf("%s/%s.go", codeInfoFolderPath, lambdaName)
	codeInfoBuiltPath := fmt.Sprintf("%s/%s", codeInfoFolderPath, fullyQualifiedLambdaName)
	codeInfoZippedPath := fmt.Sprintf("%s/%s.zip", codeInfoFolderPath, fullyQualifiedLambdaName)
	lambdaManifest := LambdaManifest{
		LambdaName:               lambdaName,
		FullyQualifiedLambdaName: fullyQualifiedLambdaName,
		CodeInfo: LambdaCodeInfo{
			FolderPath: codeInfoFolderPath,
			CodePath:   codeInfoCodePath,
			BuiltPath:  codeInfoBuiltPath,
			ZippedPath: codeInfoZippedPath,
		},
		DeployInfo: LambdaDeployInfo{
			Region:           region,
			LastDeployedHash: "",
			Arn:              "",
		},
		ExecutorRoleManifest: erm,
	}
	lm = &lambdaManifest
	err = lm.writeInitialFile(o)
	return
}

func (lm *LambdaManifest) writeInitialFile(o *output.Output) (err error) {
	filePath := lm.CodeInfo.CodePath
	err = os.MkdirAll(lm.CodeInfo.FolderPath, 0777)
	if err == nil {
		err = ioutil.WriteFile(filePath, []byte(
			fmt.Sprintf(initialLambdaFileContents, lm.LambdaName, lm.LambdaName, lm.FullyQualifiedLambdaName)), 0777)
	}
	return
}

func (lm *LambdaManifest) packageToDeploy(o *output.Output) (err error) {
	o.Info("LambdaManifest - packageToDeploy - Build").Indent()
	buildArgs := strings.Split(fmt.Sprintf("GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o %s %s", lm.CodeInfo.BuiltPath, lm.CodeInfo.CodePath), " ")
	ctx, cancelBuild := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelBuild()
	if result, err := exec.CommandContext(ctx, "env", buildArgs...).CombinedOutput(); err != nil {
		o.Failure(string(result))
		return err
	}
	o.Dedent().Done()

	o.Info("LambdaManifest - packageToDeploy - Zip").Indent()
	zipArgs := strings.Split(fmt.Sprintf("-j %s/%s.zip %s/%s", lm.CodeInfo.ZippedPath, lm.CodeInfo.BuiltPath), " ")
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
	o.Info("LambdaManifest - %s.PushToPlatform", lm.LambdaName).Indent()

	err = lm.ExecutorRoleManifest.PushToPlatform(o)
	if err != nil {
		o.Error(err)
		return
	}

	var currentCodeHash string
	if lm.DeployInfo.LastDeployedHash != "" {
		o.Info("LambdaManifest - PushToPlatform - Is Push Necessary?").Indent()
		currentCodeHash, err = file_hash.ComputeFileHash(lm.CodeInfo.CodePath, o)
		o.Info("Old Hash: %s", lm.DeployInfo.LastDeployedHash)
		o.Info("Current Hash: %s", currentCodeHash)
		if err != nil {
			return err
		}
		if currentCodeHash == lm.DeployInfo.LastDeployedHash {
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
	zipBytes, err := ioutil.ReadFile(lm.CodeInfo.ZippedPath)
	if err != nil {
		return
	}
	o.Dedent().Done()

	svc := lambda.New(session.New(), aws.NewConfig().WithRegion(lm.DeployInfo.Region))

	o.Info("LambdaManifest - PushToPlatform - Check If Lambda Exists").Indent()
	_, err = svc.GetFunction(&lambda.GetFunctionInput{
		FunctionName: aws.String(lm.FullyQualifiedLambdaName),
	})
	var arn string
	if err == nil {
		o.Warning("Lambda Already Exists.").Dedent().Done()
		o.Info("LambdaManifest - PushToPlatform - Update Lambda").Indent()
		updateFunctionRequest := &lambda.UpdateFunctionCodeInput{
			ZipFile:      zipBytes,
			FunctionName: aws.String(lm.FullyQualifiedLambdaName),
			Publish:      aws.Bool(true),
		}
		updateResult, err := svc.UpdateFunctionCode(updateFunctionRequest)
		if err != nil {
			return err
		}
		arn = *updateResult.FunctionArn
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
		createResult, err := svc.CreateFunction(createFunctionRequest)
		if err != nil {
			return err
		}
		arn = *createResult.FunctionArn
		o.Dedent().Done()
	}
	lm.DeployInfo.LastDeployedHash = currentCodeHash
	lm.DeployInfo.Arn = arn
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
	svc := lambda.New(session.New(), aws.NewConfig().WithRegion(lm.DeployInfo.Region))
	_, err = svc.DeleteFunction(deleteFunctionRequest)
	if err != nil {
		o.Error(err)
		return err
	}
	lm.DeployInfo.Arn = ""
	lm.DeployInfo.LastDeployedHash = ""
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
