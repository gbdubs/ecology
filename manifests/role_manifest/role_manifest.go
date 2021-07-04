package role_manifest

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/gbdubs/ecology/util/output"
)

type RoleConfigInfo struct {
	Name string
}

type RoleDeployInfo struct {
	RoleId           string
	Arn              string
	ExistsOnPlatform bool
}

type RoleManifest struct {
	Config RoleConfigInfo
	Deploy RoleDeployInfo
}

func New(roleName string) RoleManifest {
	rm := RoleManifest{
		Config: RoleConfigInfo{
			Name: roleName,
		},
		Deploy: RoleDeployInfo{
			ExistsOnPlatform: false,
			RoleId:           "",
			Arn:              "",
		},
	}
	return rm
}

func (rm *RoleManifest) PushToPlatform(o *output.Output) (err error) {
	o.Info("Pushing Role %s To Platform", rm.Config.Name).Indent()
	if rm.Deploy.ExistsOnPlatform {
		o.Info("No Push Needed.").Dedent().Done()
		return nil
	}
	svc := iam.New(session.New())
	o.Info("Checking to see if Role %s already exists...", rm.Config.Name).Indent()
	getRoleRequest := &iam.GetRoleInput{
		RoleName: aws.String(rm.Config.Name),
	}
	role, err := svc.GetRole(getRoleRequest)
	if err == nil {
		o.Info("Role already exists.").Dedent().Done().Dedent().Done()
		rm.Deploy.ExistsOnPlatform = true
		rm.Deploy.Arn = *role.Role.Arn
		rm.Deploy.RoleId = *role.Role.RoleId
		return
	} else {
		o.Warning("Role does not exist.").Dedent()
	}

	o.Info("Creating Role %s on Platform", rm.Config.Name).Indent()
	createRoleRequest := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(allowAmazonToRunLambdaPolicy),
		Path:                     aws.String("/"),
		RoleName:                 aws.String(rm.Config.Name),
	}
	result, err := svc.CreateRole(createRoleRequest)
	if err != nil {
		o.Error(err)
		return
	}
	rm.Deploy.Arn = *result.Role.Arn
	rm.Deploy.RoleId = *result.Role.RoleId
	rm.Deploy.ExistsOnPlatform = true
	o.Info("Role ARN = %s", rm.Deploy.Arn)
	o.Info("Role Id = %s", rm.Deploy.RoleId)
	o.Dedent().Done().Dedent().Done()
	return
}

func (rm *RoleManifest) DeleteFromPlatform(o *output.Output) (err error) {
	o.Info("Removing Role %s From Platform", rm.Config.Name).Indent()
	if !rm.Deploy.ExistsOnPlatform {
		o.Info("No Removal Needed.").Dedent().Done()
		return nil
	}
	svc := iam.New(session.New())
	o.Info("Deleting Role %s from Platform", rm.Config.Name)
	deleteRoleRequest := &iam.DeleteRoleInput{
		RoleName: aws.String(rm.Config.Name),
	}
	_, err = svc.DeleteRole(deleteRoleRequest)
	if err != nil {
		o.Error(err)
		return
	}
	rm.Deploy.Arn = ""
	rm.Deploy.RoleId = ""
	rm.Deploy.ExistsOnPlatform = false
	o.Success("Deleted Successfully.")
	o.Dedent().Done()
	return
}

const allowAmazonToRunLambdaPolicy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}`
