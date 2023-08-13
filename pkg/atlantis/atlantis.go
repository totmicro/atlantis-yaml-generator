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
	Automerge               string
	ParallelApply           string
	ParallelPlan            string
	TfRootDir               string
	OutputFile              string
	Workflow                string
	WhenModified            []string
	ExcludedProjects        string
	IncludedProjects        string
	ProjectsPatternDetector string
}

type ProjectFolder struct {
	Path          string
	WorkspaceList []string
}

func GenerateAtlantisYAML(gh github.GithubRequest, at Parameters) error {
	projectFoldersList, err := scanProjectFolders(at.TfRootDir, at.Workflow, gh)
	if err != nil {
		return err
	}
	projectFoldersListWithWorkspaces, err := detectProjectWorkspaces(projectFoldersList, at.Workflow, gh)
	if err != nil {
		return err
	}
	atlantisProjects, err := generateAtlantisProjects(at.Workflow, projectFoldersListWithWorkspaces)
	if err != nil {
		return err
	}
	filteredAtlantisProjects, err := filterAtlantisProjects(at.ExcludedProjects, at.IncludedProjects, atlantisProjects)
	if err != nil {
		return err
	}
	atlantisConfig, err := generateAtlantisConfig(at, filteredAtlantisProjects)
	if err != nil {
		return err
	}
	err = generateOutputYAML(&atlantisConfig, at.OutputFile)
	if err != nil {
		return err
	}
	return nil
}

func filterAtlantisProjects(excludedProjects, includedProjects string, atlantisProjects []Project) (filteredAtlantisProjects []Project, err error) {
	projectFilter := helpers.ProjectRegexFilter{
		Includes: includedProjects,
		Excludes: excludedProjects,
	}
	for _, project := range atlantisProjects {
		projectFilterResult, err := helpers.ProjectFilter(project.Name, projectFilter)
		if err != nil {
			return nil, err
		}

		if projectFilterResult {
			filteredAtlantisProjects = append(filteredAtlantisProjects, project)
		}
	}
	return filteredAtlantisProjects, nil
}

func detectProjectWorkspaces(foldersList []ProjectFolder, workflow string, gh github.GithubRequest) ([]ProjectFolder, error) {
	switch workflow {
	case "single-workspace":
		foldersList, _ = singleWorkspaceDetectProjectWorkspaces(foldersList)
	case "multi-workspace":
		foldersList, _ = multiWorkspaceDetectProjectWorkspaces(gh, foldersList)
	}
	return foldersList, nil
}

func scanProjectFolders(basePath, workflow string, gh github.GithubRequest) (projectFolders []ProjectFolder, err error) {
	changedFiles, err := github.GetChangedFiles(gh)
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil {
			return err
		}

		relPath, _ := filepath.Rel(basePath, filepath.Dir(path))
		workflowFilterResult := workflowFilter(info, path, workflow)
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

func workflowFilter(info os.FileInfo, path, workflow string) bool {
	switch workflow {
	case "single-workspace":
		return singleWorkspaceWorkflowFilter(info, path)
	case "multi-workspace":
		return multiWorkspaceWorkflowFilter(info, path)
	}
	return true
}

func prFilter(relPath string, changedFiles []string) bool {
	for _, file := range changedFiles {
		if strings.HasPrefix(file, fmt.Sprintf("%s/", relPath)) {
			return true
		}
	}
	return false
}

func generateAtlantisProjects(workflow string, projectFolderList []ProjectFolder) (projects []Project, err error) {
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

func genProjectName(path, workspace string) string {
	if workspace != "default" {
		return fmt.Sprintf("%s-%s", strings.Replace(path, "/", "-", 1), workspace)
	}
	return strings.Replace(path, "/", "-", 1)
}

func generateAtlantisConfig(at Parameters, projects []Project) (Config, error) {
	// Parse atlantis parameters
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
	yamlBytes, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	err = helpers.WriteFile(string(yamlBytes), outputFile)
	return err
}
