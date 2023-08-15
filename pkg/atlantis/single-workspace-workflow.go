package atlantis

import (
	"os"
	"strings"
)

func singleWorkspaceWorkflowFilter(info os.FileInfo, path, patternDetector string) bool {
	return !info.IsDir() &&
		info.Name() == patternDetector &&
		!strings.Contains(path, ".terraform")
}

func singleWorkspaceDetectProjectWorkspaces(foldersList []ProjectFolder) ([]ProjectFolder, error) {
	for i := range foldersList {
		foldersList[i].WorkspaceList = []string{"default"}
	}
	return foldersList, nil

}
