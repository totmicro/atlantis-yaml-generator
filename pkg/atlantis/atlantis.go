package atlantis

// TODO
// improve error handling
// Add comments
// Add tests

import (
	"errors"
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

var (
	//changedFiles     []string
	resources        []Resource
	atlantisProjects []Project
)

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

type Resource struct {
	Path          string
	Project       string
	WorkspaceList []string
}

func GenAtlantisYaml(gh github.GithubRequest, at Parameters) error {
	err := errors.New("")

	//workflow = at.Workflow
	//whenModified = at.WhenModified
	//projectPatternDetector = at.ProjectsPatternDetector

	//projectFilter = helpers.ProjectRegexFilter{
	//	Includes: at.IncludedProjects,
	//	Excludes: at.ExcludedProjects,
	//}

	// Get the changed files from the PR
	changedFiles, err := github.GetChangedFiles(gh)
	if err != nil {
		return err
	}

	// Check the workflow type
	switch at.Workflow {
	case "multi-workspace":
		err = multiWorkspace(at, changedFiles)
		if err != nil {
			return err
		}
	case "single-workspace":
		err = singleWorkspace(at, changedFiles)
		if err != nil {
			return err
		}
	default:
		err = fmt.Errorf("workflow %s not supported", at.Workflow)
		return err
	}
	return err
}

func generateConfig(at Parameters) (Config, error) {
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
	for _, info := range atlantisProjects {
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

func scanFolders(at Parameters, changedFiles []string) error {
	err := filepath.Walk(at.TfRootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if workflowFilter(info, path, at.Workflow) {
			project, err := helpers.GetRelativePath(path, at.TfRootDir)
			if err != nil {
				return err
			}
			if prChangesFilter(project, changedFiles) {
				addResource(filepath.Dir(path), project, at.TfRootDir, at.Workflow, changedFiles)
			}
		}
		return nil
	})
	return err
}

func workflowFilter(info os.FileInfo, path string, workflow string) bool {
	switch workflow {
	case "multi-workspace":
		return multiWorkspaceFilter(info, path)

	case "single-workspace":
		return singleWorkspaceFilter(info, path)
	}
	return false
}

func prChangesFilter(project string, changedFiles []string) bool {
	for _, file := range changedFiles {
		if strings.HasPrefix(file, fmt.Sprintf("%s/", project)) {
			return true
		}
	}
	return false
}

func addResource(path, project, tfRootDir string, workflow string, changedFiles []string) error {
	switch workflow {
	case "multi-workspace":
		return multiWorkspaceAddResource(path, project, tfRootDir, changedFiles)
	case "single-workspace":
		return singleWorkspaceAddResource(path, project, tfRootDir, changedFiles)
	}
	return nil
}

func genAtlantisProjects(workflow string, projectFilter helpers.ProjectRegexFilter) error {
	for _, folder := range resources {
		for _, workspace := range folder.WorkspaceList {
			name := fmt.Sprintf("%s-%s", strings.Replace(folder.Project, "/", "-", 1), workspace)
			if helpers.ProjectFilter(name, projectFilter) {
				atlantisProjects = append(atlantisProjects, Project{
					Name:      name,
					Dir:       folder.Path,
					Workspace: workspace,
					Workflow:  workflow,
				})
			}
		}
	}
	return nil
}

func genOutput(config *Config, outputFile string) error {
	yamlBytes, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	err = helpers.WriteFile(string(yamlBytes), outputFile)
	return err
}
