package main

import (
	"errors"
	"flag"
	"github.com/gbdubs/ecology/commands/create_lambda"
	"github.com/gbdubs/ecology/commands/create_project"
	"github.com/gbdubs/ecology/commands/delete_lambda"
	"github.com/gbdubs/ecology/commands/delete_project"
	"github.com/gbdubs/ecology/commands/list_project"
	"github.com/gbdubs/ecology/commands/push_lambda"
	"github.com/gbdubs/ecology/commands/push_project"
	"github.com/gbdubs/ecology/manifests/ecology_manifest"
	"github.com/gbdubs/ecology/util/output"
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
	pushProjectCommand := flag.NewFlagSet("push_project", flag.ExitOnError)
	deleteProjectCommand := flag.NewFlagSet("delete_project", flag.ExitOnError)

	createLambdaCommand := flag.NewFlagSet("create_lambda", flag.ExitOnError)
	pushLambdaCommand := flag.NewFlagSet("push_lambda", flag.ExitOnError)
	deleteLambdaCommand := flag.NewFlagSet("delete_lambda", flag.ExitOnError)

	// Common Flag Arguments

	platformFlagKey := "platform"
	platformDefaultValue := "AWS"
	platformHelpText := "The name of the platform that should be used for this command, AWS or GCP."
	// create_project.platform
	createProjectPlatformPtr := createProjectCommand.String(platformFlagKey, platformDefaultValue, platformHelpText)
	// create_lambda.platform
	createLambdaPlatformPtr := createLambdaCommand.String(platformFlagKey, platformDefaultValue, platformHelpText)

	regionFlagKey := "region"
	regionDefaultValue := "us-west-2"
	regionHelpText := "The name of the region that new resources should be created in"
	// create_project.region
	createProjectRegionPtr := createProjectCommand.String(regionFlagKey, regionDefaultValue, regionHelpText)
	// create_lambda.region
	createLambdaRegionPtr := createLambdaCommand.String(regionFlagKey, regionDefaultValue, regionHelpText)

	projectFlagKey := "project"
	projectDefaultValue := ""
	projectHelpText := "The name of the project that this command should operate over."
	// create_project.project
	createProjectProjectPtr := createProjectCommand.String(projectFlagKey, projectDefaultValue, projectHelpText)
	// push_project.project
	pushProjectProjectPtr := pushProjectCommand.String(projectFlagKey, projectDefaultValue, projectHelpText)
	// delete_project.project
	deleteProjectProjectPtr := deleteProjectCommand.String(projectFlagKey, projectDefaultValue, projectHelpText)
	// create_lambda.project
	createLambdaProjectPtr := createLambdaCommand.String(projectFlagKey, projectDefaultValue, projectHelpText)
	// push_lambda.project
	pushLambdaProjectPtr := pushLambdaCommand.String(projectFlagKey, projectDefaultValue, projectHelpText)
	// delete_lambda.project
	deleteLambdaProjectPtr := deleteLambdaCommand.String(projectFlagKey, projectDefaultValue, projectHelpText)

	lambdaFlagKey := "lambda"
	lambdaDefaultValue := ""
	lambdaHelpText := "The name of the lambda that this command should operate over."
	// create_lambda.lambda
	createLambdaLambdaPtr := createLambdaCommand.String(lambdaFlagKey, lambdaDefaultValue, lambdaHelpText)
	// push_lambda.lambda
	pushLambdaLambdaPtr := pushLambdaCommand.String(lambdaFlagKey, lambdaDefaultValue, lambdaHelpText)
	// delete_lambda.lambda
	deleteLambdaLambdaPtr := deleteLambdaCommand.String(lambdaFlagKey, lambdaDefaultValue, lambdaHelpText)

	verboseFlagKey := "verbose"
	verboseDefaultValue := false
	verboseHelpText := "Whether or not to be verbose in the resulting output."
	// list_project.verbose
	listProjectVerbosePtr := listProjectCommand.Bool(verboseFlagKey, verboseDefaultValue, verboseHelpText)

	illegalCommandNameError := errors.New(`Invalid commands - implemented commands:
	help
	
	create_project
	list_project
	push_project
	delete_project
	
	create_lambda
	push_lambda
	delete_lambda`)

	if len(os.Args) < 2 {
		o.Error(illegalCommandNameError)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create_project":
		createProjectCommand.Parse(os.Args[2:])
		cpc := &create_project.CreateProjectCommand{
			EcologyManifest: ecologyManifest,
			Platform:        *createProjectPlatformPtr,
			Region:          *createProjectRegionPtr,
			Project:         *createProjectProjectPtr,
		}
		cpc.Execute(o)
	case "list_project":
		listProjectCommand.Parse(os.Args[2:])
		list_project.ListProjectCommand{
			EcologyManifest: ecologyManifest,
			Verbose:         *listProjectVerbosePtr,
		}.Execute(o)
	case "push_project":
		pushProjectCommand.Parse(os.Args[2:])
		push_project.PushProjectCommand{
			EcologyManifest: ecologyManifest,
			Project:         *pushProjectProjectPtr,
		}.Execute(o)
	case "delete_project":
		deleteProjectCommand.Parse(os.Args[2:])
		delete_project.DeleteProjectCommand{
			EcologyManifest: ecologyManifest,
			Project:         *deleteProjectProjectPtr,
		}.Execute(o)
	case "create_lambda":
		createLambdaCommand.Parse(os.Args[2:])
		create_lambda.CreateLambdaCommand{
			EcologyManifest: ecologyManifest,
			Project:         *createLambdaProjectPtr,
			Lambda:          *createLambdaLambdaPtr,
			Platform:        *createLambdaPlatformPtr,
			Region:          *createLambdaRegionPtr,
		}.Execute(o)
	case "push_lambda":
		pushLambdaCommand.Parse(os.Args[2:])
		push_lambda.PushLambdaCommand{
			EcologyManifest: ecologyManifest,
			Project:         *pushLambdaProjectPtr,
			Lambda:          *pushLambdaLambdaPtr,
		}.Execute(o)
	case "delete_lambda":
		deleteLambdaCommand.Parse(os.Args[2:])
		delete_lambda.DeleteLambdaCommand{
			EcologyManifest: ecologyManifest,
			Project:         *deleteLambdaProjectPtr,
			Lambda:          *deleteLambdaLambdaPtr,
		}.Execute(o)
	default:
		o.Error(illegalCommandNameError)
	}
}
