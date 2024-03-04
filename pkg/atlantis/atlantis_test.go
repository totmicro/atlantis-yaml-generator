package atlantis

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/totmicro/atlantis-yaml-generator/pkg/config"
	"github.com/totmicro/atlantis-yaml-generator/pkg/helpers"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestApplyProjectFilter(t *testing.T) {
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
		{
			name:             "InvalidIncludeRegex",
			includedProjects: "[a-zA-Z", //Invalid regex
			excludedProjects: "project2",
			expectedProjects: nil,
			expectedError:    true,
		},
		{
			name:             "InvalidExcludeRegex",
			includedProjects: "project1", //Invalid regex
			excludedProjects: "[a-zA-Z",
			expectedProjects: nil,
			expectedError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filteredProjects, err := applyProjectFilter(tc.excludedProjects, tc.includedProjects, atlantisProjects)

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
			name:         "NonDefaultWorkspace",
			path:         "project/path/subpath",
			workspace:    "test",
			expectedName: "project-path-subpath-test",
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
	tests := []struct {
		name           string
		autoMerge      string
		parallelApply  string
		parallelPlan   string
		whenModified   string
		projects       []Project
		expectedConfig Config
		expectedError  bool
	}{
		{
			name:          "Unparsable autoMerge",
			autoMerge:     "truee",
			parallelApply: "false",
			parallelPlan:  "true",
			whenModified:  "*.tf,*.yaml",
			projects:      []Project{},
			expectedError: true,
		},
		{
			name:          "Unparsable parallelApply",
			autoMerge:     "true",
			parallelApply: "falsee",
			parallelPlan:  "true",
			whenModified:  "*.tf,*.yaml",
			projects:      []Project{},
			expectedError: true,
		},
		{
			name:          "Unparsable parallelPlan",
			autoMerge:     "true",
			parallelApply: "false",
			parallelPlan:  "truee",
			whenModified:  "*.tf,*.yaml",
			projects:      []Project{},
			expectedError: true,
		},
		{
			name:          "Default Values",
			autoMerge:     "true",
			parallelApply: "false",
			parallelPlan:  "true",
			whenModified:  "*.tf,*.yaml",
			projects: []Project{
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
			},
			expectedConfig: Config{
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
			},
		},
		// Add more test cases here with different values for automerge, parallelapply, parallelplan, and whenmodified
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, err := generateAtlantisConfig(test.autoMerge, test.parallelApply, test.parallelPlan, test.whenModified, test.projects)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, reflect.DeepEqual(config, test.expectedConfig))
			}
		})
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
	//TODO: improve test to better test workflow
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

func TestApplyPRFilter(t *testing.T) {
	// Define some sample data
	projectFolders := []ProjectFolder{
		{Path: "project1"},
		{Path: "project2"},
		{Path: "project3"},
	}

	changedFiles := []string{"project1/main.tf", "project2/main.tf"}

	// Test case 1: Filtered projects
	filteredProjectFolders, err := applyPRFilter(projectFolders, changedFiles)
	if err != nil {
		t.Errorf("Expected no error, but got an error: %v", err)
	}

	expectedFilteredProjects := []ProjectFolder{
		{Path: "project1"},
		{Path: "project2"},
	}

	if len(filteredProjectFolders) != len(expectedFilteredProjects) {
		t.Errorf("Expected %d filtered projects, but got %d", len(expectedFilteredProjects), len(filteredProjectFolders))
	}

	for i, project := range filteredProjectFolders {
		if project.Path != expectedFilteredProjects[i].Path {
			t.Errorf("Expected project path %s, but got %s", expectedFilteredProjects[i].Path, project.Path)
		}
	}

}

func TestScanProjectFolders(t *testing.T) {
	tests := []struct {
		basePath              string
		discoveryMode         string
		patternDetector       string
		expectedProjectFolder []ProjectFolder
		expectedError         bool
	}{
		{
			basePath:        "mockproject",
			discoveryMode:   "multi-workspace",
			patternDetector: "workspace_vars",
			expectedProjectFolder: []ProjectFolder{
				{Path: "multiworkspace"},
				{Path: "multiworkspace2"},
			},
		},
		{
			basePath:        "mockproject",
			discoveryMode:   "single-workspace",
			patternDetector: "main.tf",
			expectedProjectFolder: []ProjectFolder{
				{Path: "singleworkspace"},
				{Path: "singleworkspace2"},
			},
		},
		{
			basePath:        "invalidpath", // Invalid path to check file walk error
			discoveryMode:   "single-workspace",
			patternDetector: "main.tf",
			expectedProjectFolder: []ProjectFolder{
				{Path: "singleworkspace"},
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.discoveryMode, func(t *testing.T) {
			projectFolders, err := scanProjectFolders(
				OSFileSystem{},
				test.basePath,
				test.discoveryMode,
				test.patternDetector)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, projectFolders) // Assuming there are valid project folders
				assert.Equal(t, test.expectedProjectFolder, projectFolders)
			}
		})
	}
}

func TestProjectFilter(t *testing.T) {
	testCases := []struct {
		name           string
		item           string
		excludes       string
		includes       string
		expectedResult bool
	}{
		{
			name:           "NoFilter",
			item:           "project123",
			excludes:       "",
			includes:       "",
			expectedResult: true,
		},
		{
			name:           "IncludeOnly",
			item:           "project123",
			excludes:       "",
			includes:       "project\\d+",
			expectedResult: true,
		},
		{
			name:           "ExcludeOnly",
			item:           "project123",
			excludes:       "project\\d+",
			includes:       "",
			expectedResult: false,
		},
		{
			name:           "IncludeAndExclude",
			item:           "project123",
			excludes:       "project123",
			includes:       "project\\d+",
			expectedResult: false,
		},
		{
			name:           "IncludeExcludeNotMatching",
			item:           "project123",
			excludes:       "proj\\d+",
			includes:       "proj\\d+",
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := projectFilter(tc.item, tc.excludes, tc.includes)
			if err != nil {
				t.Errorf("Error: %v", err)
			}
			if result != tc.expectedResult {
				t.Errorf("Expected %v, but got %v", tc.expectedResult, result)
			}
		})
	}
}

func TestDiscoveryFilter(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		discoveryMode   string
		patternDetector string
		expectedResult  bool
	}{
		{
			name:            "single-workspace-true",
			path:            "mockproject/singleworkspace/main.tf",
			discoveryMode:   "single-workspace",
			patternDetector: "main.tf",
			expectedResult:  true,
		},
		{
			name:            "single-workspace-false",
			path:            "mockproject/singleworkspace/dummy.tf",
			discoveryMode:   "single-workspace",
			patternDetector: "main.tf",
			expectedResult:  false,
		},
		{
			name:            "multi-workspace",
			path:            "mockproject/multiworkspace/workspace_vars",
			discoveryMode:   "multi-workspace",
			patternDetector: "workspace_vars",
			expectedResult:  true,
		},
		{
			name:            "undefined-mode",
			path:            "mockproject/multiworkspace/workspace_vars",
			discoveryMode:   "",
			patternDetector: "",
			expectedResult:  true,
		},

		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			info, err := os.Stat(test.path)
			if err == nil {
				result := discoveryFilter(info, test.path, test.discoveryMode, test.patternDetector)
				assert.Equal(t, test.expectedResult, result)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestDetectProjectWorkspaces(t *testing.T) {
	tests := []struct {
		name                string
		foldersList         []ProjectFolder
		discoveryMode       string
		patternDetector     string
		changedFiles        []string
		enablePRFilter      bool
		expectedFoldersList []ProjectFolder
		expectedErr         bool
	}{
		{
			name:            "single-workspace",
			foldersList:     []ProjectFolder{{Path: "mockproject/singleworkspace"}},
			discoveryMode:   "single-workspace",
			patternDetector: "main.tf",
			changedFiles:    []string{"mockproject/singleworkspace/main.tf"},
			enablePRFilter:  true,
			expectedFoldersList: []ProjectFolder{
				{Path: "mockproject/singleworkspace", WorkspaceList: []string{"default"}},
			},
			expectedErr: false,
		},
		{
			name:            "single-workspace",
			foldersList:     []ProjectFolder{{Path: "mockproject/singleworkspace"}, {Path: "mockproject/singleworkspace2"}},
			discoveryMode:   "single-workspace",
			patternDetector: "main.tf",
			changedFiles:    []string{"mockproject/singleworkspace/main.tf"},
			enablePRFilter:  false,
			expectedFoldersList: []ProjectFolder{
				{Path: "mockproject/singleworkspace", WorkspaceList: []string{"default"}},
				{Path: "mockproject/singleworkspace2", WorkspaceList: []string{"default"}},
			},
			expectedErr: false,
		},
		{
			name:            "multi-workspace-invalid-path",
			foldersList:     []ProjectFolder{{Path: "invalidpath"}},
			discoveryMode:   "multi-workspace",
			patternDetector: "workspace_vard",
			changedFiles:    []string{"invalid/workspace_vars/test1.tfvars"},
			enablePRFilter:  true,
			expectedFoldersList: []ProjectFolder{
				{Path: "invalidpath", WorkspaceList: []string{"default"}},
			},
			expectedErr: true,
		},
		{
			name:            "multi-workspace",
			foldersList:     []ProjectFolder{{Path: "mockproject/multiworkspace"}},
			discoveryMode:   "multi-workspace",
			patternDetector: "workspace_vars",
			changedFiles:    []string{"mockproject/multiworkspace/workspace_vars/test1.tfvars"},
			enablePRFilter:  true,
			expectedFoldersList: []ProjectFolder{
				{
					Path:          "mockproject/multiworkspace",
					WorkspaceList: []string{"test1"},
				},
			},
			expectedErr: false,
		},
		{
			name:            "multi-workspace",
			foldersList:     []ProjectFolder{{Path: "mockproject/multiworkspace"}, {Path: "mockproject/multiworkspace2"}},
			discoveryMode:   "multi-workspace",
			patternDetector: "workspace_vars",
			changedFiles: []string{"mockproject/multiworkspace/workspace_vars/test1.tfvars",
				"mockproject/multiworkspace2/workspace_vars/test2.tfvars",
			},
			enablePRFilter: true,
			expectedFoldersList: []ProjectFolder{
				{
					Path:          "mockproject/multiworkspace",
					WorkspaceList: []string{"test1"},
				},
				{
					Path:          "mockproject/multiworkspace2",
					WorkspaceList: []string{"test2"},
				},
			},
			expectedErr: false,
		},
		{
			name:            "multi-workspace",
			foldersList:     []ProjectFolder{{Path: "mockproject/multiworkspace"}, {Path: "mockproject/multiworkspace2"}},
			discoveryMode:   "multi-workspace",
			patternDetector: "workspace_vars",
			changedFiles: []string{"mockproject/multiworkspace/workspace_vars/test1.tfvars",
				"mockproject/multiworkspace2/workspace_vars/test2.tfvars",
			},
			enablePRFilter: false,
			expectedFoldersList: []ProjectFolder{
				{
					Path:          "mockproject/multiworkspace",
					WorkspaceList: []string{"test1"},
				},
				{
					Path:          "mockproject/multiworkspace2",
					WorkspaceList: []string{"test1", "test2"},
				},
			},
			expectedErr: false,
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			updatedFoldersList, err := detectProjectWorkspaces(test.foldersList, test.discoveryMode, test.patternDetector, test.changedFiles, test.enablePRFilter)
			if test.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedFoldersList, updatedFoldersList)
			}
		})
	}
}

func TestGenerateOutputYAML(t *testing.T) {
	// Define a sample configuration
	config := &Config{
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
			// Add more projects as needed
		},
	}

	// Define a sample output file path
	outputFile := "/tmp/test_output.yaml"
	outputType := "file"

	// Call the function and generate the YAML
	err := generateOutputYAML(config, outputFile, outputType)
	if err != nil {
		t.Errorf("Error generating output YAML: %v", err)
	}
	// Read the generated YAML and unmarshal it for comparison
	generatedYAML, err := helpers.ReadFile(outputFile)
	if err != nil {
		t.Errorf("Error reading generated YAML: %v", err)
	}

	var generatedConfig Config
	err = yaml.Unmarshal([]byte(generatedYAML), &generatedConfig)
	if err != nil {
		t.Errorf("Error unmarshaling generated YAML: %v", err)
	}

	// Compare the generated config with the original config
	assert.Equal(t, config, &generatedConfig)
}

func TestGenerateAtlantisYAML(t *testing.T) {
	// Create a temporary test directory
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "/atlantis.yaml")

	//prChangedFiles := []string{"mockproject/single-workspace/main.tf"}
	// Define a sample configuration
	config.GlobalConfig.Parameters = make(map[string]string)
	config.GlobalConfig.Parameters["gh-token"] = "test-token"
	config.GlobalConfig.Parameters["base-repo-owner"] = "test-owner"
	config.GlobalConfig.Parameters["base-repo-name"] = "test-repo"
	config.GlobalConfig.Parameters["pull-num"] = "55"
	config.GlobalConfig.Parameters["workflow"] = "workflow1"
	config.GlobalConfig.Parameters["discoveryMode"] = "single-workspace"
	config.GlobalConfig.Parameters["pattern-detector"] = "main.tf"
	config.GlobalConfig.Parameters["terraform-base-dir"] = "mockproject"
	config.GlobalConfig.Parameters["output-file"] = tempFile
	config.GlobalConfig.Parameters["output-type"] = "file"
	config.GlobalConfig.Parameters["parallel-apply"] = "true"
	config.GlobalConfig.Parameters["parallel-plan"] = "true"
	config.GlobalConfig.Parameters["automerge"] = "true"
	config.GlobalConfig.Parameters["pr-filter"] = "false"

	err := GenerateAtlantisYAML()
	assert.NoError(t, err)

	os.Remove(tempFile)

	config.GlobalConfig.Parameters["output-type"] = "undefined"

	err = GenerateAtlantisYAML()
	assert.Error(t, err)

	config.GlobalConfig.Parameters["output-type"] = "stdout"

	err = GenerateAtlantisYAML()
	assert.NoError(t, err)

}

// Tests multiple project hits returns unique projects only.
// e.g. if we scan for *.tf the same project isn't hit twice.
func TestScanProjectFoldersUniques(t *testing.T) {
	memFS := afero.NewMemMapFs()
	fs := afero.Afero{Fs: memFS}
	// Create directories and files
	// project3 has multiple hits
	afero.WriteFile(fs, "projects_root/project1/main.tf", []byte("content"), 0644)
	afero.WriteFile(fs, "projects_root/project2/main.tf", []byte("content"), 0644)
	afero.WriteFile(fs, "projects_root/project3/main.tf", []byte("content"), 0644)
	afero.WriteFile(fs, "projects_root/project3/outputs.tf", []byte("content"), 0644)

	// Use the fs (implementing Walkable) to call scanProjectFolders
	projectFolders, err := scanProjectFolders(fs, "projects_root", "single-workspace", `.*\.tf`)
	if err != nil {
		t.Errorf("scanProjectFolders returned an error: %v", err)
	}

	// Verify that 3 project folders were returned
	if len(projectFolders) != 3 {
		t.Errorf("Expected 3 project folders, got %d. Projects %v", len(projectFolders), projectFolders)
	}
}
