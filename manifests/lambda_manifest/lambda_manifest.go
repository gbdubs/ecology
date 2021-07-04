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
	"os/exec"
	"strings"
	"time"
)

type LambdaConfigInfo struct {
	Name               string
	FullyQualifiedName string
	FolderPath         string
	CodePath           string
	BuiltPath          string
	ZippedPath         string
}

type LambdaDeployInfo struct {
	Platform         string
	Region           string
	LastDeployedHash string
	Arn              string
}

type LambdaManifest struct {
	Config               LambdaConfigInfo
	Deploy               LambdaDeployInfo
	ExecutorRoleManifest role_manifest.RoleManifest
}

func New(projectDir string, projectName string, lambdaName string, platform string, region string, o *output.Output) (lm *LambdaManifest, err error) {
	fullyQualifiedLambdaName := projectName + "-" + lambdaName
	erm := role_manifest.New(fullyQualifiedLambdaName + "-executor")
	configInfoFolderPath := fmt.Sprintf("%s/lambda/%s", projectDir, lambdaName)
	configInfoCodePath := fmt.Sprintf("%s/%s.go", configInfoFolderPath, lambdaName)
	configInfoBuiltPath := fmt.Sprintf("%s/%s", configInfoFolderPath, fullyQualifiedLambdaName)
	configInfoZippedPath := fmt.Sprintf("%s/%s.zip", configInfoFolderPath, fullyQualifiedLambdaName)
	lambdaManifest := LambdaManifest{
		Config: LambdaConfigInfo{
			Name:               lambdaName,
			FullyQualifiedName: fullyQualifiedLambdaName,
			FolderPath:         configInfoFolderPath,
			CodePath:           configInfoCodePath,
			BuiltPath:          configInfoBuiltPath,
			ZippedPath:         configInfoZippedPath,
		},
		Deploy: LambdaDeployInfo{
			Platform:         platform,
			Region:           region,
			LastDeployedHash: "",
			Arn:              "",
		},
		ExecutorRoleManifest: erm,
	}
	lm = &lambdaManifest
	return
}

func (lm *LambdaManifest) packageToDeploy(o *output.Output) (err error) {
	o.Info("LambdaManifest - packageToDeploy - Build").Indent()
	buildArgs := strings.Split(fmt.Sprintf("GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o %s %s", lm.Config.BuiltPath, lm.Config.CodePath), " ")
	ctx, cancelBuild := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelBuild()
	if result, err := exec.CommandContext(ctx, "env", buildArgs...).CombinedOutput(); err != nil {
		o.Failure(string(result))
		return err
	}
	o.Dedent().Done()

	o.Info("LambdaManifest - packageToDeploy - Zip").Indent()
	zipArgs := strings.Split(fmt.Sprintf("-j %s %s", lm.Config.ZippedPath, lm.Config.BuiltPath), " ")
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
	o.Info("LambdaManifest - %s.PushToPlatform", lm.Config.Name).Indent()

	err = lm.ExecutorRoleManifest.PushToPlatform(o)
	if err != nil {
		o.Error(err)
		return
	}

	var currentCodeHash string
	if lm.Deploy.LastDeployedHash != "" {
		o.Info("LambdaManifest - PushToPlatform - Checking if Nescessary").Indent()
		currentCodeHash, err = file_hash.ComputeFileHash(lm.Config.CodePath)
		if err != nil {
			return err
		}
		o.Info("Old Code Hash: %s", lm.Deploy.LastDeployedHash)
		o.Info("New Code Hash: %s", currentCodeHash)
		if currentCodeHash == lm.Deploy.LastDeployedHash {
			o.Dedent().Success("Code hasn't changed since last push, no push needed.").Dedent().Done()
			return nil
		}
		o.Dedent().Warning("Code has changed since last push.")
	} else {
		o.Info("LambdaManifest - PushToPlatform - First Lambda Push")
		if currentCodeHash, err = file_hash.ComputeFileHash(lm.Config.CodePath); err != nil {
			return err
		}
	}

	err = lm.packageToDeploy(o)
	if err != nil {
		return err
	}

	o.Info("LambdaManifest - PushToPlatform - Read Zip").Indent()
	zipBytes, err := ioutil.ReadFile(lm.Config.ZippedPath)
	if err != nil {
		return
	}
	o.Dedent().Done()

	svc := lambda.New(session.New(), aws.NewConfig().WithRegion(lm.Deploy.Region))

	o.Info("LambdaManifest - PushToPlatform - Check If Lambda Exists").Indent()
	_, err = svc.GetFunction(&lambda.GetFunctionInput{
		FunctionName: aws.String(lm.Config.FullyQualifiedName),
	})
	var arn string
	if err == nil {
		o.Warning("Lambda Already Exists.").Dedent().Done()
		o.Info("LambdaManifest - PushToPlatform - Update Lambda").Indent()
		updateFunctionRequest := &lambda.UpdateFunctionCodeInput{
			ZipFile:      zipBytes,
			FunctionName: aws.String(lm.Config.FullyQualifiedName),
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
			Description:  aws.String(fmt.Sprintf("Ecology-Generated Lambda %s.", lm.Config.FullyQualifiedName)),
			FunctionName: aws.String(lm.Config.FullyQualifiedName),
			Handler:      aws.String(lm.Config.FullyQualifiedName),
			Publish:      aws.Bool(true),
			Role:         aws.String(lm.ExecutorRoleManifest.Deploy.Arn),
			Runtime:      aws.String("go1.x"),
		}
		createResult, err := svc.CreateFunction(createFunctionRequest)
		if err != nil {
			return err
		}
		arn = *createResult.FunctionArn
		o.Dedent().Done()
	}
	lm.Deploy.LastDeployedHash = currentCodeHash
	lm.Deploy.Arn = arn
	o.Dedent().Done()
	return nil
}

func (lm *LambdaManifest) DeleteFromPlatform(o *output.Output) (err error) {
	o.Info("LambdaManifest - DeleteFromPlatform - %s", lm.Config.FullyQualifiedName).Indent()

	err = lm.ExecutorRoleManifest.DeleteFromPlatform(o)
	if err != nil {
		o.Error(err)
		return
	}

	deleteFunctionRequest := &lambda.DeleteFunctionInput{
		FunctionName: aws.String(lm.Config.FullyQualifiedName),
	}
	svc := lambda.New(session.New(), aws.NewConfig().WithRegion(lm.Deploy.Region))
	_, err = svc.DeleteFunction(deleteFunctionRequest)
	if err != nil {
		o.Error(err)
		return err
	}
	lm.Deploy.Arn = ""
	lm.Deploy.LastDeployedHash = ""
	return nil
}
