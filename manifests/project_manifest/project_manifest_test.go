package project_manifest

import "os"
import "io/ioutil"
import "testing"
import "github.com/gbdubs/ecology/output"

func Test_SaveOnUndefinedManifestPathFails(t *testing.T) {
	o := output.NewForTesting()
	dirForTesting, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dirForTesting)

	projectName := "NewProjectForTest"
	projectManifest := ProjectManifest{
		ProjectName: projectName,
	}
	actualError := projectManifest.Save(o)

	expectedErrorMessage := "open : no such file or directory"
	if actualError.Error() != expectedErrorMessage {
		t.Errorf("ERROR: Errors don't match': Actual [%s] Expected [%s]", actualError.Error(), expectedErrorMessage)
	}
}

func Test_SaveAndGetWorkTogetherToPersistProjectManifests(t *testing.T) {
	o := output.NewForTesting()
	dirForTesting, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dirForTesting)
	projectName := "NewProjectForTest"
	projectManifestPath := dirForTesting + "/a/b/ProjectManifestPath/ecology.ecology"

	projectManifest := ProjectManifest{
		ProjectName:         projectName,
		ProjectManifestPath: projectManifestPath,
	}
	actualError := projectManifest.Save(o)
	if actualError != nil {
		t.Errorf("ERROR: Expected no error but was %v", actualError)
	}

	actualPM, actualError := GetProjectManifest(projectManifestPath, o)

	if actualError != nil {
		t.Errorf("ERROR: Expected no error but was %v", actualError)
	}
	if actualPM.ProjectManifestPath != projectManifestPath {
		t.Errorf("ERROR: Strings don't match': Actual [%s] Expected [%s]", actualPM.ProjectManifestPath, projectManifestPath)
	}
	if actualPM.ProjectName != projectName {
		t.Errorf("ERROR: Strings don't match': Actual [%s] Expected [%s]", actualPM.ProjectName, projectName)
	}
}
