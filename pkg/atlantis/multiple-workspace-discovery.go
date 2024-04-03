package atlantis

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/totmicro/atlantis-yaml-generator/pkg/helpers"
)

func multiWorkspaceGetProjectScope(relPath, patternDetector string, changedFiles []string) string {
	for _, file := range changedFiles {
		if strings.HasPrefix(file, fmt.Sprintf("%s/", relPath)) &&
			!helpers.MatchesPattern(patternDetector, file) {
			return "crossWorkspace"
		}
	}
	return "workspace"
}

func multiWorkspaceGenWorkspaceList(relPath string, changedFiles []string, enablePRFilter bool, scope string) (workspaceList []string, err error) {
	err = filepath.Walk(relPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(),
			tfvarsExtension) &&
			!strings.Contains(path, ".terraform") {

			if helpers.IsStringInList(path, changedFiles) ||
				(scope == "crossWorkspace" || !enablePRFilter) {
				workspaceList = append(workspaceList, helpers.TrimFileExtension(info.Name()))
			}
		}
		return nil
	})
	return workspaceList, err
}

func multiWorkspaceDetectProjectWorkspaces(changedFiles []string, enablePRFilter bool,
	foldersList []ProjectFolder, patternDetector string) (updatedFolderList []ProjectFolder, err error) {

	for i := range foldersList {
		scope := multiWorkspaceGetProjectScope(foldersList[i].Path, patternDetector, changedFiles)
		workspaceList, err := multiWorkspaceGenWorkspaceList(
			fmt.Sprintf("%s/%s", foldersList[i].Path, patternDetector),
			changedFiles,
			enablePRFilter,
			scope)
		if err != nil {
			return foldersList, err
		}
		foldersList[i].WorkspaceList = workspaceList
	}
	return foldersList, nil
}

func multiWorkspaceDiscoveryFilter(info os.FileInfo, path, patternDetector string) bool {
	return info.IsDir() &&
		helpers.MatchesPattern(patternDetector, info.Name()) &&
		!strings.Contains(path, ".terraform")
}
