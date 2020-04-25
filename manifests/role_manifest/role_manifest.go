package role_manifest

import (
	"github.com/gbdubs/ecology/output"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type RoleManifest struct {
	RoleName         string
	RoleId           string
	Arn              string
	ExistsOnPlatform    bool
}

func New(roleName string) RoleManifest {
  rm := RoleManifest {
    RoleName: roleName,
    ExistsOnPlatform: false,
    RoleId: "",
    Arn: "",
  };
  return rm;
}

func (rm *RoleManifest) PushToPlatform(o *output.Output) (err error) {
  o.Info("Pushing Role %s To Platform", rm.RoleName).Indent()
  if rm.ExistsOnPlatform {
    o.Info("No Push Needed.").Dedent().Done()
    return nil
  }
  svc := iam.New(session.New())
  o.Info("Creating Role %s on Platform", rm.RoleName)
  createRoleRequest := &iam.CreateRoleInput{
    AssumeRolePolicyDocument: aws.String(allowAmazonToRunLambdaPolicy),
    Path:                     aws.String("/"),
    RoleName:                 aws.String(rm.RoleName),
  }
  result, err := svc.CreateRole(createRoleRequest)
  if (err != nil) {
    o.Error(err);
    return
  }
  rm.Arn = *result.Role.Arn
  rm.RoleId = *result.Role.RoleId
  rm.ExistsOnPlatform = true
  o.Info("Role ARN = %s", rm.Arn)
  o.Info("Role Id = %s", rm.RoleId)
  o.Dedent().Done()
  return
}

func (rm *RoleManifest) RemoveFromPlatform(o *output.Output) (err error) {
  o.Info("Removing Role %s From Platform", rm.RoleName).Indent()
  if !rm.ExistsOnPlatform {
    o.Info("No Removal Needed.").Dedent().Done()
    return nil
  }
  svc := iam.New(session.New())
  o.Info("Deleting Role %s from Platform", rm.RoleName)
  deleteRoleRequest := &iam.DeleteRoleInput{
    RoleName: aws.String(rm.RoleName),
  }
  _, err = svc.DeleteRole(deleteRoleRequest)
  if (err != nil) {
    o.Error(err);
    return
  }
  rm.Arn = ""
  rm.RoleId = ""
  rm.ExistsOnPlatform = false
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