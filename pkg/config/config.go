package config

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/totmicro/atlantis-yaml-generator/pkg/helpers"
)

type Config struct {
	Parameters map[string]string
}

var GlobalConfig Config

type Parameter struct {
	Name         string
	Description  string
	Required     bool
	Dependencies DependentParameters
	DefaultValue string
	Shorthand    string
	Value        string
}

type DependentParameters struct {
	WhenParentParameterIs string
	ParameterList         []string
}

// ParameterList is the source of truth list for all parameters
var ParameterList = []Parameter{
	{
		Name:         "automerge",
		Description:  "Atlantis automerge config value.",
		Required:     false,
		DefaultValue: "false",
		Shorthand:    "",
	},
	{
		Name:         "parallel-apply",
		Description:  "Atlantis parallel apply config value.",
		Required:     false,
		DefaultValue: "true",
		Shorthand:    "",
	},
	{
		Name:         "parallel-plan",
		Description:  "Atlantis parallel plan config value.",
		Required:     false,
		DefaultValue: "true",
		Shorthand:    "",
	},
	{
		Name:         "terraform-base-dir",
		Description:  "Basedir for terraform resources.",
		Required:     false,
		DefaultValue: "./",
		Shorthand:    "",
	},
	{
		Name:         "output-file",
		Description:  "Atlantis YAML output file name.",
		Required:     false,
		DefaultValue: "atlantis.yaml",
		Shorthand:    "f",
	},
	{
		Name:         "output-type",
		Description:  "Atlantis YAML output type. [file|stdout].",
		Required:     false,
		DefaultValue: "file",
		Shorthand:    "e",
	},
	{
		Name:         "workflow",
		Description:  "Atlantis Workflow to be used.",
		Required:     false,
		DefaultValue: "",
		Shorthand:    "w",
	},
	{
		Name:         "when-modified",
		Description:  "Atlantis will trigger an autoplan when these modifications occur (list of strings).",
		Required:     false,
		DefaultValue: "**/*.tf,**/*.tfvars,**/*.json,**/*.tpl,**/*.tmpl,**/*.xml",
		Shorthand:    "m",
	},
	{
		Name:         "excluded-projects",
		Description:  "Atlantis regex filter to exclude projects.",
		Required:     false,
		DefaultValue: "",
		Shorthand:    "x",
	},
	{
		Name:         "included-projects",
		Description:  "Atlantis regex filter to only include projects.",
		Required:     false,
		DefaultValue: "",
		Shorthand:    "z",
	},
	{
		Name:         "pattern-detector",
		Description:  "discover projects based on files or directories names.",
		Required:     false,
		DefaultValue: "main.tf",
		Shorthand:    "q",
	},
	{
		Name:         "discovery-mode",
		Description:  "mode used to discover projects [single-workspace|multi-workspace]",
		Required:     false,
		DefaultValue: "single-workspace",
		Shorthand:    "d",
	},
	{
		Name:         "pull-num",
		Description:  "Github Pull Request Number to check diffs.",
		Required:     false,
		DefaultValue: "",
		Shorthand:    "p",
	},
	{
		Name:         "base-repo-name",
		Description:  "Github Repo Name.",
		Required:     false,
		DefaultValue: "",
		Shorthand:    "r",
	},
	{
		Name:         "base-repo-owner",
		Description:  "Github Repo Owner Name.",
		Required:     false,
		DefaultValue: "",
		Shorthand:    "o",
	},
	{
		Name:         "gh-token",
		Description:  "Specify the GitHub token when automatic detection is not possible.",
		Required:     false,
		DefaultValue: "",
		Shorthand:    "t",
	},
	{
		Name:        "pr-filter",
		Description: "Filter projects based on the PR changes (Only for github SCM).",
		Required:    false,
		Dependencies: DependentParameters{
			WhenParentParameterIs: "true",
			ParameterList:         []string{"pull-num", "base-repo-name", "base-repo-owner"}},
		DefaultValue: "false",
		Shorthand:    "u",
	},
}

// Init generates the config Parameters object and checks for missing required parameters
func Init(ccmd *cobra.Command) (err error) {
	GlobalConfig.Parameters = make(map[string]string)
	for i, param := range ParameterList {
		GlobalConfig.Parameters[param.Name] = getFlagOrEnv(ccmd, ParameterList[i].Name,
			generateEnvVarName(ParameterList[i].Name), ParameterList[i].DefaultValue)
	}
	err = CheckRequiredParameters(ParameterList)
	return err
}

// GetFlagOrEnv gets the value of a flag or environment variable and if both are empty, returns the default value
func getFlagOrEnv(ccmd *cobra.Command, flagName, envVar string, defaultValue string) string {
	val, _ := ccmd.Flags().GetString(flagName)
	if val != "" {
		return val
	}
	val = helpers.LookupEnvString(envVar)
	if val != "" {
		return val
	} else {
		return defaultValue
	}
}

func CheckRequiredParameters(parameterList []Parameter) error {
	missingParams := []string{}

	// Iterate through the list of parameters
	for _, param := range parameterList {
		if param.Required && GlobalConfig.Parameters[param.Name] == "" {
			missingParams = append(missingParams, param.Name)
		} else if len(param.Dependencies.ParameterList) > 0 {
			// Check if the parameter has dependent parameters
			parentValue := GlobalConfig.Parameters[param.Name]
			expectedParentValue := param.Dependencies.WhenParentParameterIs

			if parentValue == expectedParentValue {
				// Iterate through the dependent parameters
				for _, dependentParameter := range param.Dependencies.ParameterList {
					if GlobalConfig.Parameters[dependentParameter] == "" {
						missingParams = append(missingParams, dependentParameter)
					}
				}
			}
		}
	}

	// Check if any missing parameters were found
	if len(missingParams) > 0 {
		return fmt.Errorf("Missing required parameters: %s", strings.Join(missingParams, ", "))
	}

	return nil
}

// GenerateDescription generates the description for a parameter
func GenerateDescription(param string, description string) (updatedDescription string) {
	envVar := generateEnvVarName(param)
	updatedDescription = fmt.Sprintf("%s (Equivalent envVar: [%s])", description, envVar)
	return updatedDescription
}

// generateEnvVarName generates the environment variable name for a parameter
func generateEnvVarName(param string) string {
	param = strings.ReplaceAll(param, "-", "_")
	param = strings.ToUpper(param)
	return param
}
