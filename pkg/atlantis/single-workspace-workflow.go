package atlantis

import (
	"os"
	"strings"
)

var singleWorkspacePatternDetector = "main.tf"

func singleWorkspaceWorkflowFilter(info os.FileInfo, path string) bool {
	return !info.IsDir() &&
		info.Name() == singleWorkspacePatternDetector &&
		!strings.Contains(path, ".terraform")
}

func singleWorkspaceDetectProjectWorkspaces(foldersList []ProjectFolder) ([]ProjectFolder, error) {
	for i := range foldersList {
		foldersList[i].WorkspaceList = []string{"default"}
	}
	return foldersList, nil

}
