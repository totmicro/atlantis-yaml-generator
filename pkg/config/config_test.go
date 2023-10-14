package config

import (
	"fmt"
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

func TestCheckRequiredParameters(t *testing.T) {
	// Create a sample GlobalConfig for testing
	GlobalConfig.Parameters = map[string]string{
		"automerge":       "true",
		"parallel-apply":  "false",
		"pull-num":        "",
		"base-repo-name":  "myrepo",
		"base-repo-owner": "myowner",
		"pr-filter":       "true",
	}

	tests := []struct {
		name       string
		parameters []Parameter
		expected   error
	}{
		{
			name: "All required parameters set",
			parameters: []Parameter{
				{
					Name:     "automerge",
					Required: true,
				},
				{
					Name:     "parallel-apply",
					Required: true,
				},
			},
			expected: nil,
		},
		{
			name: "Required parameters missing",
			parameters: []Parameter{
				{
					Name:     "automerge",
					Required: true,
				},
				{
					Name:     "pull-num",
					Required: true,
				},
			},
			expected: fmt.Errorf("Missing required parameters: pull-num"),
		},
		{
			name: "Dependent parameters not triggered",
			parameters: []Parameter{
				{
					Name:         "pr-filter",
					Required:     false,
					Dependencies: DependentParameters{WhenParentParameterIs: "true", ParameterList: []string{"base-repo-name", "base-repo-owner"}},
				},
			},
			expected: nil,
		},
		{
			name: "Dependent parameters triggered",
			parameters: []Parameter{
				{
					Name:         "pr-filter",
					Required:     false,
					Dependencies: DependentParameters{WhenParentParameterIs: "true", ParameterList: []string{"mydep"}},
				},
				{
					Name:     "base-repo-name",
					Required: false,
				},
			},
			expected: fmt.Errorf("Missing required parameters: mydep"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CheckRequiredParameters(test.parameters)
			if err != nil {
				if err.Error() != test.expected.Error() {
					t.Errorf("Expected error: %s, but got: %s", test.expected, err)
				}
			} else if test.expected != nil {
				t.Errorf("Expected error: %s, but got nil", test.expected)
			}
		})
	}
}
