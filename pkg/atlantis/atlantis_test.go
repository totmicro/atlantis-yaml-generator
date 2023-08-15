package atlantis

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterAtlantisProjects(t *testing.T) {
	atlantisProjects := []Project{
		{Name: "project1"},
		{Name: "project2"},
		{Name: "excluded_project"},
	}

	testCases := []struct {
		name             string
		excludedProjects string
		includedProjects string
		expectedProjects []Project
		expectedError    bool
	}{
		{
			name:             "IncludeAllProjects",
			includedProjects: ".*",
			expectedProjects: atlantisProjects,
		},
		{
			name:             "IncludeOneProject",
			includedProjects: "project1",
			expectedProjects: []Project{
				{Name: "project1"},
			},
		},
		{
			name:             "ExcludeProject",
			excludedProjects: "excluded.*",
			expectedProjects: []Project{
				{Name: "project1"},
				{Name: "project2"},
			},
		},
		{
			name:             "IncludeOneExcludeOne",
			includedProjects: "project1",
			excludedProjects: "project2",
			expectedProjects: []Project{
				{Name: "project1"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filteredProjects, err := filterAtlantisProjects(tc.excludedProjects, tc.includedProjects, atlantisProjects)

			if tc.expectedError && err == nil {
				t.Error("Expected an error, but got nil")
			} else if !tc.expectedError && err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}

			if !reflect.DeepEqual(filteredProjects, tc.expectedProjects) {
				t.Errorf("Expected projects: %+v, but got: %+v", tc.expectedProjects, filteredProjects)
			}
		})
	}
}

func TestGenProjectName(t *testing.T) {
	testCases := []struct {
		name         string
		path         string
		workspace    string
		expectedName string
	}{
		{
			name:         "DefaultWorkspace",
			path:         "project/path",
			workspace:    "default",
			expectedName: "project-path",
		},
		{
			name:         "NonDefaultWorkspace",
			path:         "project/path",
			workspace:    "test",
			expectedName: "project-path-test",
		},
		{
			name:         "EmptyPath",
			path:         "",
			workspace:    "test",
			expectedName: "-test",
		},
		{
			name:         "EmptyWorkspace",
			path:         "project/path",
			workspace:    "",
			expectedName: "project-path-",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := genProjectName(tc.path, tc.workspace)
			if result != tc.expectedName {
				t.Errorf("Expected %s, but got %s", tc.expectedName, result)
			}
		})
	}
}

func TestGenerateAtlantisConfig(t *testing.T) {
	parameters := Parameters{
		Automerge:     "true",
		ParallelApply: "false",
		ParallelPlan:  "true",
		WhenModified:  []string{"*.tf", "*.yaml"},
	}

	projects := []Project{
		{
			Name:      "project1",
			Workspace: "default",
			Workflow:  "workflow1",
			Dir:       "dir1",
		},
		{
			Name:      "project2",
			Workspace: "test",
			Workflow:  "workflow2",
			Dir:       "dir2",
		},
	}

	expectedConfig := Config{
		Version:       3,
		Automerge:     true,
		ParallelApply: false,
		ParallelPlan:  true,
		Projects: []Project{
			{
				Name:      "project1",
				Workspace: "default",
				Workflow:  "workflow1",
				Dir:       "dir1",
				Autoplan: struct {
					Enabled      bool     `yaml:"enabled"`
					WhenModified []string `yaml:"when_modified"`
				}{
					Enabled:      true,
					WhenModified: []string{"*.tf", "*.yaml"},
				},
			},
			{
				Name:      "project2",
				Workspace: "test",
				Workflow:  "workflow2",
				Dir:       "dir2",
				Autoplan: struct {
					Enabled      bool     `yaml:"enabled"`
					WhenModified []string `yaml:"when_modified"`
				}{
					Enabled:      true,
					WhenModified: []string{"*.tf", "*.yaml"},
				},
			},
		},
	}

	config, err := generateAtlantisConfig(parameters, projects)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if !reflect.DeepEqual(config, expectedConfig) {
		t.Errorf("Expected config: %+v, but got: %+v", expectedConfig, config)
	}
}

func TestGenerateAtlantisProjects(t *testing.T) {
	projectFolders := []ProjectFolder{
		{
			Path:          "project1",
			WorkspaceList: []string{"default", "test"},
		},
		{
			Path:          "project2",
			WorkspaceList: []string{"dev", "prod"},
		},
	}

	workflow := "myWorkflow"

	expectedProjects := []Project{
		{
			Name:      "project1",
			Dir:       "project1",
			Workspace: "default",
			Workflow:  "myWorkflow",
		},
		{
			Name:      "project1-test",
			Dir:       "project1",
			Workspace: "test",
			Workflow:  "myWorkflow",
		},
		{
			Name:      "project2-dev",
			Dir:       "project2",
			Workspace: "dev",
			Workflow:  "myWorkflow",
		},
		{
			Name:      "project2-prod",
			Dir:       "project2",
			Workspace: "prod",
			Workflow:  "myWorkflow",
		},
	}

	projects, err := generateAtlantisProjects(workflow, projectFolders)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if !reflect.DeepEqual(projects, expectedProjects) {
		t.Errorf("Expected projects: %+v, but got: %+v", expectedProjects, projects)
	}
}

func TestSingleWorkspaceDetectProjectWorkspaces(t *testing.T) {
	projectFolders := []ProjectFolder{
		{
			Path:          "project1",
			WorkspaceList: []string{},
		},
		{
			Path:          "project2",
			WorkspaceList: []string{},
		},
	}

	foldersList, err := singleWorkspaceDetectProjectWorkspaces(projectFolders)

	assert.NoError(t, err)

	for _, folder := range foldersList {
		assert.Equal(t, []string{"default"}, folder.WorkspaceList)
	}
}
func TestMultiWorkspaceDetectProjectWorkspaces(t *testing.T) {
	foldersList := []ProjectFolder{
		{
			Path:          "mockproject",
			WorkspaceList: []string{},
		},
	}

	changedFiles := []string{"mockproject/multiworkspace/workspace_vars/test1.tfvars", "project2/file2.tfvars"}

	foldersList, err := detectProjectWorkspaces(foldersList, "multi-workspace", "workspace_vars", changedFiles)

	assert.NoError(t, err)

	expectedFoldersList := []ProjectFolder{
		{
			Path:          "mockproject",
			WorkspaceList: []string{"test1"},
		},
	}

	assert.Equal(t, expectedFoldersList, foldersList)
}

func TestMultiWorkspaceGetProjectScope(t *testing.T) {
	changedFiles := []string{"mockproject/multiworkspace/workspace_vars/file1.tfvars"}

	scope := multiWorkspaceGetProjectScope("mockproject/multiworkspace", "workspace_vars", changedFiles)
	assert.Equal(t, "workspace", scope)

	scope = multiWorkspaceGetProjectScope("project2", "workspace_vars", changedFiles)
	assert.Equal(t, "workspace", scope)

	changedFiles = append(changedFiles, "mockproject/multiworkspace/main.tf")
	scope = multiWorkspaceGetProjectScope("mockproject/multiworkspace", "workspace_vars", changedFiles)
	assert.Equal(t, "crossWorkspace", scope)
}

func TestMultiWorkspaceGenWorkspaceList(t *testing.T) {
	changedFiles := []string{"mockproject/multiworkspace/workspace_vars/test1.tfvars"}

	workspaceList, err := multiWorkspaceGenWorkspaceList("mockproject/multiworkspace", changedFiles, "workspace")
	assert.NoError(t, err)
	assert.Equal(t, []string{"test1"}, workspaceList)

}
func TestPrFilter(t *testing.T) {
	changedFiles := []string{
		"project1/file1.tf",
		"project1/file2.tf",
		"project2/file1.tf",
	}

	relPath := "project1"

	result := prFilter(relPath, changedFiles)
	assert.True(t, result, "Expected true, but got false")

	relPath = "project3"

	result = prFilter(relPath, changedFiles)
	assert.False(t, result, "Expected false, but got true")
}

func TestScanProjectFolders(t *testing.T) {
	tests := []struct {
		basePath              string
		workflow              string
		patternDetector       string
		changedFiles          []string
		expectedProjectFolder []ProjectFolder
	}{
		{
			basePath:        "mockproject",
			workflow:        "multi-workspace",
			patternDetector: "workspace_vars",
			changedFiles: []string{
				"multiworkspace/workspace_vars/test1.tfvars",
				"multiworkspace2/workspace_vars/test1.tfvars"},
			expectedProjectFolder: []ProjectFolder{
				{Path: "multiworkspace"},
				{Path: "multiworkspace2"},
			},
		},
		{
			basePath:        "mockproject",
			workflow:        "single-workspace",
			patternDetector: "main.tf",
			changedFiles:    []string{"singleworkspace/main.tf", "file2"},
			expectedProjectFolder: []ProjectFolder{
				{Path: "singleworkspace"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.workflow, func(t *testing.T) {
			projectFolders, err := scanProjectFolders(test.basePath, test.workflow, test.patternDetector, test.changedFiles)
			assert.NoError(t, err)
			assert.NotEmpty(t, projectFolders) // Assuming there are valid project folders
			assert.Equal(t, test.expectedProjectFolder, projectFolders)
		})
	}
}
