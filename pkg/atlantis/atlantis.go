package atlantis

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/totmicro/atlantis-yaml-generator/pkg/github"
	"github.com/totmicro/atlantis-yaml-generator/pkg/helpers"
	"gopkg.in/yaml.v3"
)

// DefaultWhenModified is the default list of file extensions to trigger a plan
var DefaultWhenModified = []string{
	"**/*.tf",
	"**/*.tfvars",
	"**/*.json",
	"**/*.tpl",
	"**/*.tmpl",
	"**/*.xml",
}

// Default values for PatternDetector when not defined
var WorkflowPatternDetectorMap = map[string]string{
	"single-workspace": "main.tf",
	"multi-workspace":  "workspace_vars",
}

// Default values when not defined
const (
	DefaultOutputFile        = "atlantis.yaml"
	DefaultAutomerge         = "false"
	DefaultParallelPlan      = "true"
	DefaultParallelApply     = "true"
	DefaultAtlantisTfRootDir = "./"
)

type Config struct {
	Version       int       `yaml:"version"`
	Automerge     bool      `yaml:"automerge"`
	ParallelApply bool      `yaml:"parallel_apply"`
	ParallelPlan  bool      `yaml:"parallel_plan"`
	Projects      []Project `yaml:"projects"`
}

type Project struct {
	Name      string `yaml:"name"`
	Workspace string `yaml:"workspace"`
	Workflow  string `yaml:"workflow"`
	Dir       string `yaml:"dir"`
	Autoplan  struct {
		Enabled      bool     `yaml:"enabled"`
		WhenModified []string `yaml:"when_modified"`
	} `yaml:"autoplan"`
}

type Parameters struct {
	Automerge        string
	ParallelApply    string
	ParallelPlan     string
	TfRootDir        string
	OutputFile       string
	Workflow         string
	WhenModified     []string
	ExcludedProjects string
	IncludedProjects string
	PatternDetector  string
}

type ProjectFolder struct {
	Path          string
	WorkspaceList []string
}

// GenerateAtlantisYAML generates the atlantis.yaml file
func GenerateAtlantisYAML(gh github.GithubRequest, at Parameters) error {
	// Get the changed files from the PR
	prChangedFiles, err := github.GetChangedFiles(gh)
	if err != nil {
		return err
	}
	// Scan folders to detect projects
	projectFoldersList, err := scanProjectFolders(at.TfRootDir, at.Workflow, at.PatternDetector, prChangedFiles)
	if err != nil {
		return err
	}
	// Detect project workspaces
	projectFoldersListWithWorkspaces, err := detectProjectWorkspaces(projectFoldersList, at.Workflow, at.PatternDetector, prChangedFiles)
	if err != nil {
		return err
	}
	// Generate atlantis projects
	atlantisProjects, err := generateAtlantisProjects(at.Workflow, projectFoldersListWithWorkspaces)
	if err != nil {
		return err
	}
	// Filter atlantis projects with included and excluded regex rules
	filteredAtlantisProjects, err := filterAtlantisProjects(at.ExcludedProjects, at.IncludedProjects, atlantisProjects)
	if err != nil {
		return err
	}
	// Generate atlantis config to later render the atlantis.yaml file
	atlantisConfig, err := generateAtlantisConfig(at, filteredAtlantisProjects)
	if err != nil {
		return err
	}
	// Generate atlantis.yaml file
	err = generateOutputYAML(&atlantisConfig, at.OutputFile)
	if err != nil {
		return err
	}
	return nil
}

func scanProjectFolders(basePath, workflow, patternDetector string, changedFiles []string) (projectFolders []ProjectFolder, err error) {
	// Scan folders for projects and apply filters
	err = filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil {
			return err
		}
		relPath, _ := filepath.Rel(basePath, filepath.Dir(path))
		// Detect projects folders based on the workflow
		workflowFilterResult := workflowFilter(info, path, workflow, patternDetector)
		// Filter projects based on the PR changed files
		prFilterResult := prFilter(relPath, changedFiles)
		if workflowFilterResult && prFilterResult {
			projectFolders = append(projectFolders, ProjectFolder{
				Path: relPath,
			})
		}
		return nil
	})
	return projectFolders, err
}

func detectProjectWorkspaces(foldersList []ProjectFolder, workflow string, patternDetector string, changedFiles []string) (updatedFoldersList []ProjectFolder, err error) {
	// Detect project workspaces based on the workflow
	switch workflow {
	case "single-workspace":
		updatedFoldersList, err = singleWorkspaceDetectProjectWorkspaces(foldersList)
	case "multi-workspace":
		updatedFoldersList, err = multiWorkspaceDetectProjectWorkspaces(changedFiles, foldersList, patternDetector)
	}
	// You can add more workflows rules here if required
	return updatedFoldersList, err
}

func generateAtlantisProjects(workflow string, projectFolderList []ProjectFolder) (projects []Project, err error) {
	// Iterate over the project folders and generate atlantis projects
	for _, folder := range projectFolderList {
		for _, workspace := range folder.WorkspaceList {
			name := genProjectName(folder.Path, workspace)
			projects = append(projects, Project{
				Name:      name,
				Dir:       folder.Path,
				Workspace: workspace,
				Workflow:  workflow,
			})
		}
	}
	return projects, nil
}

func filterAtlantisProjects(excludedProjects, includedProjects string, atlantisProjects []Project) (filteredAtlantisProjects []Project, err error) {
	// Create project filter with included and excluded regex rules
	projectFilter := helpers.ProjectRegexFilter{
		Includes: includedProjects,
		Excludes: excludedProjects,
	}
	// Iterate over atlantis projects and filter them
	for _, project := range atlantisProjects {
		projectFilterResult, err := helpers.ProjectFilter(project.Name, projectFilter)
		if err != nil {
			return filteredAtlantisProjects, err
		}
		if projectFilterResult {
			filteredAtlantisProjects = append(filteredAtlantisProjects, project)
		}
	}
	return filteredAtlantisProjects, nil
}

func generateAtlantisConfig(at Parameters, projects []Project) (Config, error) {
	// Parse atlantis parameters to detect config values
	automerge, err := strconv.ParseBool(at.Automerge)
	if err != nil {
		return Config{}, err
	}
	parallelApply, _ := strconv.ParseBool(at.ParallelApply)
	if err != nil {
		return Config{}, err
	}
	parallelPlan, _ := strconv.ParseBool(at.ParallelPlan)
	if err != nil {
		return Config{}, err
	}

	// Generate the atlantis base config
	config := Config{
		Version:       3,
		Automerge:     automerge,
		ParallelApply: parallelApply,
		ParallelPlan:  parallelPlan,
	}
	// Append generated projects to the atlantis config
	for _, info := range projects {
		project := Project{
			Name:      info.Name,
			Workspace: info.Workspace,
			Workflow:  info.Workflow,
			Dir:       info.Dir,
			Autoplan: struct {
				Enabled      bool     `yaml:"enabled"`
				WhenModified []string `yaml:"when_modified"`
			}{
				Enabled:      true,
				WhenModified: at.WhenModified,
			},
		}
		config.Projects = append(config.Projects, project)
	}
	return config, err
}

func generateOutputYAML(config *Config, outputFile string) error {
	// Generate the atlantis.yaml file
	yamlBytes, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	err = helpers.WriteFile(string(yamlBytes), outputFile)
	return err
}

func workflowFilter(info os.FileInfo, path, workflow, patternDetector string) bool {
	// Detect projects folders based on the workflow
	// Each workflow has different rules to detect projects
	switch workflow {
	case "single-workspace":
		return singleWorkspaceWorkflowFilter(info, path, patternDetector)
	case "multi-workspace":
		return multiWorkspaceWorkflowFilter(info, path, patternDetector)
	}
	// You can add more workflows rules here if required
	return true
}

func prFilter(relPath string, changedFiles []string) bool {
	// Filter projects based on the PR changed files
	for _, file := range changedFiles {
		if strings.HasPrefix(file, fmt.Sprintf("%s/", relPath)) {
			return true
		}
	}
	return false
}

func genProjectName(path, workspace string) string {
	// Generate project name based on the path and workspace
	if workspace != "default" {
		return fmt.Sprintf("%s-%s", strings.Replace(path, "/", "-", 1), workspace)
	}
	return strings.Replace(path, "/", "-", 1)
}
