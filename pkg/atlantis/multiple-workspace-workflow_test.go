package atlantis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiWorkspaceGetProjectScope(t *testing.T) {
	changedFiles := []string{"mockproject/multiworkspace/workspace_vars/file1.tfvars"}

	scope := multiWorkspaceGetProjectScope("mockproject/multiworkspace", "workspace_vars", changedFiles)
	assert.Equal(t, "workspace", scope)

	scope = multiWorkspaceGetProjectScope("project2", "workspace_vars", changedFiles)
	assert.Equal(t, "workspace", scope)

	changedFiles = append(changedFiles, "mockproject/multiworkspace/main.tf")
	scope = multiWorkspaceGetProjectScope("mockproject/multiworkspace", "workspace_vars", changedFiles)
	assert.Equal(t, "crossWorkspace", scope)
}

func TestMultiWorkspaceGenWorkspaceList(t *testing.T) {
	changedFiles := []string{"mockproject/multiworkspace/workspace_vars/test1.tfvars"}

	workspaceList, err := multiWorkspaceGenWorkspaceList("mockproject/multiworkspace", changedFiles, "workspace")
	assert.NoError(t, err)
	assert.Equal(t, []string{"test1"}, workspaceList)

}
