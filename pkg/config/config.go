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
	DefaultValue string
	Shorthand    string
	Value        string
}

// ParameterList is the source of truth list for all parameters
var ParameterList = []Parameter{
	{
		Name:         "automerge",
		Description:  "Atlantis automerge config value.",
		Required:     false,
		DefaultValue: "true",
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
		Description:  "Atlantis output file name.",
		Required:     false,
		DefaultValue: "atlantis.yaml",
		Shorthand:    "f",
	},
	{
		Name:         "workflow",
		Description:  "Atlantis Workflow to be used. [single-workspace|multi-workspace].",
		Required:     true,
		DefaultValue: "",
		Shorthand:    "w",
	},
	{
		Name:         "when-modified",
		Description:  "Atlantis When modified (list of strings) to run autoplan.",
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
		Required:     true,
		DefaultValue: "",
		Shorthand:    "q",
	},
	{
		Name:         "pull-num",
		Description:  "Github Pull Request Number to check diffs.",
		Required:     true,
		DefaultValue: "",
		Shorthand:    "p",
	},
	{
		Name:         "base-repo-name",
		Description:  "Github Repo Name.",
		Required:     true,
		DefaultValue: "",
		Shorthand:    "r",
	},
	{
		Name:         "base-repo-owner",
		Description:  "Github Repo Owner Name.",
		Required:     true,
		DefaultValue: "",
		Shorthand:    "o",
	},
	{
		Name:         "gh-token",
		Description:  "Github Token Value.",
		Required:     true,
		DefaultValue: "",
		Shorthand:    "t",
	},
}

// Init generates the config Parameters object and checks for missing required parameters
func Init(ccmd *cobra.Command) (err error) {
	GlobalConfig.Parameters = make(map[string]string)
	for i, param := range ParameterList {
		GlobalConfig.Parameters[param.Name] = getFlagOrEnv(ccmd, ParameterList[i].Name,
			generateEnvVarName(ParameterList[i].Name), ParameterList[i].DefaultValue)
	}
	err = CheckRequiredPamameters()
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

// CheckRequiredPamameters checks for missing required parameters
func CheckRequiredPamameters() (err error) {
	missingParams := []string{}
	for _, param := range ParameterList {
		if param.Required && GlobalConfig.Parameters[param.Name] == "" {
			missingParams = append(missingParams, param.Name)
		}
	}
	if len(missingParams) > 0 {
		err = fmt.Errorf("Missing required parameters: %s", strings.Join(missingParams, ", "))
		return err
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
