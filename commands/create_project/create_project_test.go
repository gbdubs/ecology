package create_project

import "os"
import "io/ioutil"
import "testing"
import "github.com/gbdubs/ecology/ecology_manifest"
import "github.com/gbdubs/ecology/project_manifest"
import "github.com/gbdubs/ecology/output"

func Test_applyConstantTransformations_CreatesProjectManifestPathFromProjectPath(t *testing.T) {
	createProjectCommand := CreateProjectCommand {
	  EcologyManifest: ecology_manifest.EcologyManifest {
	    EcologyProjectsDirectoryPath: "/BD",
	  },
		ProjectPath: "a/b/c/d/project",
	}

	modifiedCreateProjectCommand := applyConstantTransformations(createProjectCommand)

	assertStringsEqual(t, modifiedCreateProjectCommand.ProjectManifestPath, "/BD/a/b/c/d/project/ecology.ecology")
}

func Test_applyConstantTransformations_CreatesProjectManifestPathFromProjectName(t *testing.T) {
	createProjectCommand := CreateProjectCommand {
	  EcologyManifest: ecology_manifest.EcologyManifest {
	    EcologyProjectsDirectoryPath: "/BD",
	  },
		ProjectSimpleName: "SimpleProjectName",
	}

	modifiedCreateProjectCommand := applyConstantTransformations(createProjectCommand)

	assertStringsEqual(t, modifiedCreateProjectCommand.ProjectManifestPath, "/BD/SimpleProjectName/ecology.ecology")
}

func Test_applyConstantTransformations_InfersSimpleProjectNameIfUnspecified(t *testing.T) {
	createProjectCommand := CreateProjectCommand {
	  EcologyManifest: ecology_manifest.EcologyManifest {},
		ProjectPath: "a/b/c/d/ProjectName",
	}

	modifiedCreateProjectCommand := applyConstantTransformations(createProjectCommand)

	assertStringsEqual(t, modifiedCreateProjectCommand.ProjectSimpleName, "ProjectName")
}

func Test_applyConstantTransformations_InfersProjectPathIfUnspecified(t *testing.T) {
	createProjectCommand := CreateProjectCommand {
	  EcologyManifest: ecology_manifest.EcologyManifest {},
		ProjectSimpleName: "ProjectName",
	}

	modifiedCreateProjectCommand := applyConstantTransformations(createProjectCommand)

	assertStringsEqual(t, modifiedCreateProjectCommand.ProjectPath, "ProjectName")
}

func Test_validate_ValidArguments(t *testing.T) {
	createProjectCommand := CreateProjectCommand {
	  EcologyManifest: ecology_manifest.EcologyManifest {
	    EcologyProjectsDirectoryPath: "/BD",
	    ProjectManifestPaths: []string {"/BD/P1/ecology.ecology", "/BD/P2/ecology.ecology"},
	  },
		ProjectPath: "Hello/World/P3",
		LambdaName:  "ValidLambdaName",
		Platform:    "AWS",
	}

	actualError := createProjectCommand.validate()

	if actualError != nil {
		t.Errorf("No validation error expected, but one was present %v", actualError)
	}
}

func Test_validate_ProjectName_MustBeNonEmpty(t *testing.T) {
	createProjectCommand := CreateProjectCommand{}

	actualError := createProjectCommand.validate()

	assertErrorsEqual(t, actualError, "create_project requires --project_name or --project_path")
}

func Test_validate_ProjectName_MustNotContainSpaces(t *testing.T) {
	createProjectCommand := CreateProjectCommand{
		ProjectSimpleName: "Hello World",
	}

	actualError := createProjectCommand.validate()

	assertErrorsEqual(t, actualError, "--project_name can only contain alphanumeric characters")
}

func Test_validate_ProjectName_MustNotContainSymbols(t *testing.T) {
	createProjectCommand := CreateProjectCommand{
		ProjectSimpleName: "HelloWorld!",
	}

	actualError := createProjectCommand.validate()

	assertErrorsEqual(t, actualError, "--project_name can only contain alphanumeric characters")
}

func Test_validate_ProjectPath_MustBeNonEmpty(t *testing.T) {
	createProjectCommand := CreateProjectCommand{}

	actualError := createProjectCommand.validate()

	assertErrorsEqual(t, actualError, "create_project requires --project_name or --project_path")
}

func Test_validate_ProjectPath_MustNotContainSpaces(t *testing.T) {
	createProjectCommand := CreateProjectCommand{
		ProjectPath: "Hello/World/Hello World",
	}

	actualError := createProjectCommand.validate()

	assertErrorsEqual(t, actualError, "--project_path can only contain alphanumeric characters or slashes")
}

func Test_validate_ProjectPath_MustNotContainSymbols(t *testing.T) {
	createProjectCommand := CreateProjectCommand{
		ProjectPath: "Hello/world/HelloWorld!",
	}

	actualError := createProjectCommand.validate()

	assertErrorsEqual(t, actualError, "--project_path can only contain alphanumeric characters or slashes")
}


func Test_validate_LambdaName_MustBeNonEmpty(t *testing.T) {
	createProjectCommand := CreateProjectCommand{
		ProjectSimpleName: "HelloWorld",
		LambdaName:  "",
	}

	actualError := createProjectCommand.validate()

	assertErrorsEqual(t, actualError, "create_project requires --lambda_name")
}

func Testvalidate_LambdaName_MustNotContainSpaces(t *testing.T) {
	createProjectCommand := CreateProjectCommand{
		ProjectSimpleName: "HelloWorld",
		LambdaName:  "Hello World",
	}

	actualError := createProjectCommand.validate()

	assertErrorsEqual(t, actualError, "--lambda_name can only contain alphanumeric characters")
}

func Test_validate_LambdaName_MustNotContainSymbols(t *testing.T) {
	createProjectCommand := CreateProjectCommand{
		ProjectSimpleName: "HelloWorld",
		LambdaName:  "HelloWorld!",
	}

	actualError := createProjectCommand.validate()

	assertErrorsEqual(t, actualError, "--lambda_name can only contain alphanumeric characters")
}

func Test_validate_Platform_MustBeOnEnumeratedList(t *testing.T) {
	createProjectCommand := CreateProjectCommand{
		ProjectSimpleName: "HelloWorld",
		LambdaName:  "HelloWorld",
		Platform:    "ABC",
	}

	actualError := createProjectCommand.validate()

	assertErrorsEqual(t, actualError, "--platform should be one of AWS or GCP")
}

func Test_validate_Platform_MustBeSupportedOnEnumeratedList(t *testing.T) {
	createProjectCommand := CreateProjectCommand{
		ProjectSimpleName: "HelloWorld",
		LambdaName:  "HelloWorld",
		Platform:    "GCP",
	}

	actualError := createProjectCommand.validate()

	assertErrorsEqual(t, actualError, "--platform=GCP is not yet supported for this command")
}

func Test_validate_DoesntAllowCreationOfDuplicateProjects(t *testing.T) {
	createProjectCommand := CreateProjectCommand {
	  EcologyManifest: ecology_manifest.EcologyManifest {
	    EcologyProjectsDirectoryPath: "/BD",
	    ProjectManifestPaths: []string {"/BD/P1/ecology.ecology", "/BD/P2/ecology.ecology"},
	  },
		ProjectPath: "P1",
		ProjectManifestPath: "/BD/P1/ecology.ecology",
		LambdaName:  "ValidLambdaName",
		Platform:    "AWS",
	}

	actualError := createProjectCommand.validate()

	assertErrorsEqual(t, actualError, "Already a project at /BD/P1/ecology.ecology")
}

func Test_Execute_CreatesProjectManifest(t *testing.T) {
  o := output.NewForTesting()
  dir, err := ioutil.TempDir("", "")
  if err != nil {
    panic(err)
  }
  defer os.RemoveAll(dir)
  projectName := "NewProjectForTest"
  createProjectCommand := CreateProjectCommand {
	  EcologyManifest: ecology_manifest.EcologyManifest {
	    EcologyProjectsDirectoryPath: dir,
	  },
		ProjectSimpleName: projectName,
		LambdaName:  "ValidLambdaName",
		Platform:    "AWS",
	}
  err = createProjectCommand.Execute(o)
  if err != nil { panic (err) }
  
  expectedProjectManifestPath := dir + "/"+projectName+"/ecology.ecology"
  actualPM, err := project_manifest.GetProjectManifest(expectedProjectManifestPath, o)
  if err != nil {
    panic(err)
  }
  assertStringsEqual(t, actualPM.ProjectManifestPath, expectedProjectManifestPath)
  assertStringsEqual(t, actualPM.ProjectName, projectName)
}

func assertStringsEqual(t *testing.T, actual string, expected string) {
  if actual != expected {
		t.Errorf("ERROR: Strings don't match': Actual [%s] Expected [%s]", actual, expected)
  }
}

func assertErrorsEqual(t *testing.T, actualError error, expectedErrorMessage string) {
	if actualError == nil {
		t.Errorf("ERROR: Expected error, but got nil")
	} else if actualError.Error() != expectedErrorMessage {
		t.Errorf("ERROR: Errors don't match': Actual [%s] Expected [%s]", actualError.Error(), expectedErrorMessage)
	}
}
