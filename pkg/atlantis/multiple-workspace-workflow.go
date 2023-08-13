package atlantis

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/totmicro/atlantis-yaml-generator/pkg/github"
	"github.com/totmicro/atlantis-yaml-generator/pkg/helpers"
)

const tfvarsExtension = ".tfvars"

var multiWorkspacePatternDetector = "workspace_vars"

func multiWorkspaceGetProjectScope(relPath string, changedFiles []string) string {
	for _, file := range changedFiles {
		if strings.HasPrefix(file, fmt.Sprintf("%s/", relPath)) &&
			!strings.Contains(file, multiWorkspacePatternDetector) {
			return "crossWorkspace"
		}
	}
	return "workspace"
}

func multiWorkspaceGenWorkspaceList(relPath string, changedFiles []string, scope string) (workspaceList []string, err error) {
	err = filepath.Walk(relPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), tfvarsExtension) {
			if helpers.IsStringInList(path, changedFiles) || (scope == "crossWorkspace") {
				workspaceList = append(workspaceList, helpers.TrimFileExtension(info.Name()))
			}
		}
		return nil
	})
	return workspaceList, err
}

func multiWorkspaceDetectProjectWorkspaces(gh github.GithubRequest, foldersList []ProjectFolder) ([]ProjectFolder, error) {
	changedFiles, _ := github.GetChangedFiles(gh)

	for i := range foldersList {
		scope := multiWorkspaceGetProjectScope(foldersList[i].Path, changedFiles)
		workspaceList, _ := multiWorkspaceGenWorkspaceList(foldersList[i].Path, changedFiles, scope)
		foldersList[i].WorkspaceList = workspaceList
	}
	return foldersList, nil
}

func multiWorkspaceWorkflowFilter(info os.FileInfo, path string) bool {
	return info.IsDir() &&
		info.Name() == "workspace_vars" &&
		!strings.Contains(path, ".terraform")

}
