package helpers

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

type ProjectRegexFilter struct {
	Excludes string
	Includes string
}

func LookupEnvString(key string, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultValue
}

func CheckEnvVars(key string, envVar string, defaultValue string) string {
	if key != "" {
		return key
	} else {
		return LookupEnvString(envVar, defaultValue)
	}
}

func GetFlagOrEnv(ccmd *cobra.Command, flagName, envVar string, defaultValue string) string {
	val, _ := ccmd.Flags().GetString(flagName)
	return CheckEnvVars(val, envVar, defaultValue)
}

func WriteFile(content, filePath string) error {
	// Open the file for writing, creating it if it doesn't exist
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Convert the string to a byte slice
	contentBytes := []byte(content)

	// Write the content to the file
	_, err = file.Write(contentBytes)
	if err != nil {
		return err
	}
	return err
}

func GetRelativePath(path, basePath string) (string, error) {
	return filepath.Rel(basePath, filepath.Dir(path))
}

func TrimFileExtension(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func IsStringInList(value string, list []string) bool {
	for _, str := range list {
		if strings.Contains(value, str) {
			return true
		}
	}
	return false
}

func ProjectFilter(item string, filter ProjectRegexFilter) (result bool, err error) {
	// If the regexp is not defined, we don't filter the project
	if filter.Includes == "" && filter.Excludes == "" {
		return true, nil
	}

	// Compile the regular expressions
	var patternInclude, patternExclude *regexp.Regexp
	if filter.Includes != "" {
		patternInclude, err = regexp.Compile(filter.Includes)
	}
	if filter.Excludes != "" {
		patternExclude, err = regexp.Compile(filter.Excludes)
	}

	// Check if the item matches the include and exclude patterns
	if patternInclude != nil && !patternInclude.MatchString(item) {
		return false, nil
	}
	if patternExclude != nil && patternExclude.MatchString(item) {
		return false, nil
	}

	return true, nil
}

func CreateProjectFilter(includes string, excludes string) ProjectRegexFilter {
	return ProjectRegexFilter{
		Excludes: excludes,
		Includes: includes,
	}
}
