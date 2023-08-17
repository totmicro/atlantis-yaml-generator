package config

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	cmd := &cobra.Command{}
	// Check for required parameters
	requiredParams := []string{}
	for _, param := range ParameterList {
		if param.Required {
			requiredParams = append(requiredParams, param.Name)
		}
	}
	err := Init(cmd)
	if len(requiredParams) > 0 {
		assert.Error(t, err, "Expected error for missing required parameters: %v", requiredParams)
	} else {
		assert.NoError(t, err, "Expected no error, but got: %v", err)
	}
	// Add cmd flags for required parameters
	for _, param := range requiredParams {
		cmd.Flags().StringP(param, "", "dummy", "dummy")
	}
	err = Init(cmd)
	assert.NoError(t, err, "Expected no error, but got: %v", err)
	// Remove cmd Flags and add environment variables for required parameters
	for _, param := range requiredParams {
		cmd.ResetFlags()
		os.Setenv(generateEnvVarName(param), "dummy")
	}
	err = Init(cmd)
	assert.NoError(t, err, "Expected no error, but got: %v", err)
}

func TestGenerateDescription(t *testing.T) {
	testCases := []struct {
		param       string
		description string
		expected    string
	}{
		{
			param:       "automerge",
			description: "Atlantis automerge config value.",
			expected:    "Atlantis automerge config value. (Equivalent envVar: [AUTOMERGE])",
		},
		{
			param:       "parallel-apply",
			description: "Atlantis parallel apply config value.",
			expected:    "Atlantis parallel apply config value. (Equivalent envVar: [PARALLEL_APPLY])",
		},
		// Add more test cases as needed
	}
	for _, tc := range testCases {
		result := GenerateDescription(tc.param, tc.description)
		assert.Equal(t, tc.expected, result, "Generated description doesn't match expected")
	}
}
