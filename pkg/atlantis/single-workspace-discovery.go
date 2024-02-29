package atlantis

import (
	"github.com/totmicro/atlantis-yaml-generator/pkg/helpers"
	"os"
	"strings"
)

func singleWorkspaceDiscoveryFilter(info os.FileInfo, path, patternDetector string) bool {
	return !info.IsDir() &&
		helpers.MatchesPattern(patternDetector, info.Name()) &&
		!strings.Contains(path, ".terraform")
}

func singleWorkspaceDetectProjectWorkspaces(foldersList []ProjectFolder) ([]ProjectFolder, error) {
	for i := range foldersList {
		foldersList[i].WorkspaceList = []string{"default"}
	}
	return foldersList, nil

}
