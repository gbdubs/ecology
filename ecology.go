package main

import (
	"errors"
	"flag"
	"github.com/gbdubs/ecology/create_project"
	"github.com/gbdubs/ecology/ecology_manifest"
	"github.com/gbdubs/ecology/list_project"
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
	listProjectCommand := flag.NewFlagSet("list_project", flag.ExitOnError)

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

	verboseFlagKey := "verbose"
	verboseDefaultValue := false
	verboseHelpText := "Whether or not to be verbose in the resulting output."
	listProjectVerbosePtr := listProjectCommand.Bool(verboseFlagKey, verboseDefaultValue, verboseHelpText)

	illegalCommandNameError := errors.New("Invalid command. Implemented Commands:\n  create_project\n  list_project")
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
	case "list_project":
		listProjectCommand.Parse(os.Args[2:])
		list_project.ListProjectCommand{
			EcologyManifest: ecologyManifest,
			Verbose:         *listProjectVerbosePtr,
		}.Execute(o)
	default:
		o.Error(illegalCommandNameError)
	}
}
