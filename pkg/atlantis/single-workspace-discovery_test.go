package atlantis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleWorkspaceDetectProjectWorkspaces(t *testing.T) {
	projectFolders := []ProjectFolder{
		{
			Path:          "project1",
			WorkspaceList: []string{},
		},
		{
			Path:          "project2",
			WorkspaceList: []string{},
		},
	}

	foldersList, err := singleWorkspaceDetectProjectWorkspaces(projectFolders)

	assert.NoError(t, err)

	for _, folder := range foldersList {
		assert.Equal(t, []string{"default"}, folder.WorkspaceList)
	}
}
