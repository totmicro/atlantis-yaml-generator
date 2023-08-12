package atlantis

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/totmicro/atlantis-yaml-generator/pkg/helpers"
)

const tfvarsExtension = ".tfvars"

var multiWorkspacePatternDetector = "workspace_vars"

func multiWorkspace(at Parameters, changedFiles []string) error {
	// Set the project pattern detector
	if at.ProjectsPatternDetector != "" {
		multiWorkspacePatternDetector = at.ProjectsPatternDetector
	}

	// Scan folders generate the atlantis projects
	err := scanFolders(at, changedFiles)
	if err != nil {
		return err
	}
	// Generate atlantis projects
	err = genAtlantisProjects(at.Workflow, helpers.CreateProjectFilter(at.IncludedProjects, at.ExcludedProjects))
	if err != nil {
		return err
	}
	// Generate atlantis config struct
	atlantisConfig, err := generateConfig(at)
	if err != nil {
		return err
	}
	// Render atlantis yaml manifest
	err = genOutput(&atlantisConfig, at.OutputFile)
	if err != nil {
		return err
	}

	return err
}

func multiWorkspaceFilter(info os.FileInfo, path string) bool {
	return info.IsDir() &&
		info.Name() == multiWorkspacePatternDetector &&
		!strings.Contains(path, ".terraform")
}

func multiWorkspaceAddResource(path, project, tfRootDir string, changedFiles []string) error {
	workspaceList, err := multiWorkspaceGetWorkspaces(project, tfRootDir, changedFiles)
	if err != nil {
		return err
	}
	resources = append(resources, Resource{
		Path:          path,
		Project:       project,
		WorkspaceList: workspaceList,
	})
	return err
}

func multiWorkspaceGetWorkspaces(project string, rootDir string, changedFiles []string) ([]string, error) {
	var matchingWorkspaces []string
	scope := multiWorkspaceGetProjectScope(project, changedFiles)

	err := filepath.Walk(filepath.Join(rootDir, project), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), tfvarsExtension) {
			if helpers.IsStringInList(path, changedFiles) || (scope == "crossWorkspace") {
				matchingWorkspaces = append(matchingWorkspaces, helpers.TrimFileExtension(info.Name()))
			}
		}
		return nil
	})
	return matchingWorkspaces, err
}

func multiWorkspaceGetProjectScope(project string, changedFiles []string) string {
	for _, file := range changedFiles {
		if strings.HasPrefix(file, fmt.Sprintf("%s/", project)) &&
			!strings.Contains(file, multiWorkspacePatternDetector) {
			return "crossWorkspace"
		}
	}
	return "workspace"
}
