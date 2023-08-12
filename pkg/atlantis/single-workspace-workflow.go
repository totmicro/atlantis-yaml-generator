package atlantis

import (
	"os"
	"strings"

	"github.com/totmicro/atlantis-yaml-generator/pkg/helpers"
)

var singleWorkspacePatternDetector = "main.tf"

func singleWorkspace(at Parameters, changedFiles []string) error {
	// Set the project pattern detector
	if at.ProjectsPatternDetector != "" {
		singleWorkspacePatternDetector = at.ProjectsPatternDetector
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

func singleWorkspaceFilter(info os.FileInfo, path string) bool {
	return !info.IsDir() && info.Name() == singleWorkspacePatternDetector && !strings.Contains(path, ".terraform")
}

func singleWorkspaceAddResource(path string, project string, tfRootDir string, changedFiles []string) error {
	resources = append(resources, Resource{
		Path:          path,
		Project:       project,
		WorkspaceList: []string{"default"},
	})
	return nil
}
