package main

import (
	"errors"
	"flag"
	"github.com/gbdubs/ecology/create_project"
	"github.com/gbdubs/ecology/ecology_manifest"
	"github.com/gbdubs/ecology/output"
	"os"
)

func main() {
	o := output.New()
	o.Info("Reading Ecology Manifest...").Indent()
	ecologyManifest, err := ecology_manifest.GetEcologyManifest(o)
	o.Dedent().Done()

	if err != nil {
		o.Error(err)
	}

	// Subcommands
	createProjectCommand := flag.NewFlagSet("create_project", flag.ExitOnError)

	// Common Flag Arguments
	projectPathFlagKey := "project_path"
	projectPathDefaultValue := ""
	projectPathHelpText := "The path to the base directory of the project that this command should operate over."
	createProjectProjectPathPtr := createProjectCommand.String(projectPathFlagKey, projectPathDefaultValue, projectPathHelpText)

	projectNameFlagKey := "project_name"
	projectNameDefaultValue := ""
	projectNameHelpText := "The name of the project that this command should operate over."
	createProjectProjectNamePtr := createProjectCommand.String(projectNameFlagKey, projectNameDefaultValue, projectNameHelpText)

	lambdaNameFlagKey := "lambda_name"
	lambdaNameDefaultValue := ""
	lambdaNameHelpText := "The name of the lambda that this command should operate over."
	createProjectLambdaNamePtr := createProjectCommand.String(lambdaNameFlagKey, lambdaNameDefaultValue, lambdaNameHelpText)

	platformFlagKey := "platform"
	platformDefaultValue := ""
	platformHelpText := "The name of the platform that should be used for this command, AWS or GCP."
	createProjectPlatformPtr := createProjectCommand.String(platformFlagKey, platformDefaultValue, platformHelpText)

	illegalCommandNameError := errors.New("Invalid command. Implemented Commands:\n create_project")
	if len(os.Args) < 2 {
		o.Error(illegalCommandNameError)
	}

	switch os.Args[1] {
	case "create_project":
		createProjectCommand.Parse(os.Args[2:])
		create_project.CreateProjectCommand{
			EcologyManifest:   ecologyManifest,
			ProjectSimpleName: *createProjectProjectNamePtr,
			ProjectPath:       *createProjectProjectPathPtr,
			LambdaName:        *createProjectLambdaNamePtr,
			Platform:          *createProjectPlatformPtr,
		}.Execute(o)
	default:
		o.Error(illegalCommandNameError)
	}
}
