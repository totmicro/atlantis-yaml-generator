package atlantis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterAtlantisProjects(t *testing.T) {
	excludedProjects := "excluded\\d"
	includedProjects := "included\\d"

	atlantisProjects := []Project{
		{Name: "project1"},
		{Name: "excluded1"},
		{Name: "project2"},
		{Name: "excluded2"},
		{Name: "included1"},
	}

	filteredProjects, err := filterAtlantisProjects(excludedProjects, includedProjects, atlantisProjects)

	assert.NoError(t, err, "FilterAtlantisProjects should not return an error")
	assert.Len(t, filteredProjects, 2, "Number of filtered projects should match")

	expectedFilteredProjects := []Project{
		{Name: "project1"},
		{Name: "project2"},
	}

	assert.Equal(t, expectedFilteredProjects, filteredProjects)
}
