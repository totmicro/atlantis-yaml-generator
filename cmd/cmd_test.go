package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/totmicro/atlantis-yaml-generator/pkg/atlantis"
)

func TestDefineWhenModifiedList(t *testing.T) {
	defaultList := atlantis.DefaultWhenModified
	customList := []string{"x", "y", "z"}

	assert.Equal(t, defaultList, defineWhenModifiedList(""))
	assert.Equal(t, customList, defineWhenModifiedList("x,y,z"))
}

func TestDefineProjectPatternDetector(t *testing.T) {
	testCases := []struct {
		name            string
		patternDetector string
		workflow        string
		expectedResult  string
	}{
		{
			name:            "Default for single-workspace",
			patternDetector: "",
			workflow:        "single-workspace",
			expectedResult:  "main.tf",
		},
		{
			name:            "Default for multi-workspace",
			patternDetector: "",
			workflow:        "multi-workspace",
			expectedResult:  "workspace_vars",
		},
		{
			name:            "Custom for multi-workspace",
			patternDetector: "mypattern",
			workflow:        "multi-workspace",
			expectedResult:  "mypattern",
		},
		{
			name:            "Custom for single-workspace",
			patternDetector: "mypattern",
			workflow:        "single-workspace",
			expectedResult:  "mypattern",
		},
		{
			name:            "Undefined workflow",
			patternDetector: "",
			workflow:        "undefined-workspace",
			expectedResult:  "Workflow undefined-workspace not found",
		},
		{
			name:            "Undefined workflow with custom pattern",
			patternDetector: "mypattern",
			workflow:        "undefined-workspace",
			expectedResult:  "Workflow undefined-workspace not found",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := defineProjectPatternDetector(tc.patternDetector, tc.workflow)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
