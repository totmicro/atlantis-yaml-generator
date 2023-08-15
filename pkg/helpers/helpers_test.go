package helpers

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestProjectFilter(t *testing.T) {
	testCases := []struct {
		name           string
		item           string
		filter         ProjectRegexFilter
		expectedResult bool
	}{
		{
			name:           "NoFilter",
			item:           "project123",
			filter:         ProjectRegexFilter{},
			expectedResult: true,
		},
		{
			name:           "IncludeOnly",
			item:           "project123",
			filter:         ProjectRegexFilter{Includes: "project\\d+"},
			expectedResult: true,
		},
		{
			name:           "ExcludeOnly",
			item:           "project123",
			filter:         ProjectRegexFilter{Excludes: "project\\d+"},
			expectedResult: false,
		},
		{
			name:           "IncludeAndExclude",
			item:           "project123",
			filter:         ProjectRegexFilter{Includes: "project\\d+", Excludes: "project123"},
			expectedResult: false,
		},
		{
			name:           "IncludeExcludeNotMatching",
			item:           "project123",
			filter:         ProjectRegexFilter{Includes: "proj\\d+", Excludes: "proj\\d+"},
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ProjectFilter(tc.item, tc.filter)
			if err != nil {
				t.Errorf("Error: %v", err)
			}
			if result != tc.expectedResult {
				t.Errorf("Expected %v, but got %v", tc.expectedResult, result)
			}
		})
	}
}

func TestCreateProjectFilter(t *testing.T) {
	testCases := []struct {
		name     string
		includes string
		excludes string
	}{
		{
			name:     "NoFilter",
			includes: "",
			excludes: "",
		},
		{
			name:     "IncludeOnly",
			includes: "project\\d+",
			excludes: "",
		},
		{
			name:     "ExcludeOnly",
			includes: "",
			excludes: "project\\d+",
		},
		{
			name:     "IncludeAndExclude",
			includes: "project\\d+",
			excludes: "project123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filter := CreateProjectFilter(tc.includes, tc.excludes)
			if filter.Includes != tc.includes || filter.Excludes != tc.excludes {
				t.Errorf("Expected includes: %v, excludes: %v, but got includes: %v, excludes: %v",
					tc.includes, tc.excludes, filter.Includes, filter.Excludes)
			}
		})
	}
}

func TestTrimFileExtension(t *testing.T) {
	testCases := []struct {
		name          string
		filename      string
		expectedValue string
	}{
		{
			name:          "NoExtension",
			filename:      "file",
			expectedValue: "file",
		},
		{
			name:          "WithExtension",
			filename:      "file.txt",
			expectedValue: "file",
		},
		{
			name:          "MultipleDots",
			filename:      "file.version.txt",
			expectedValue: "file.version",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := TrimFileExtension(tc.filename)
			if result != tc.expectedValue {
				t.Errorf("Expected %s, but got %s", tc.expectedValue, result)
			}
		})
	}
}

func TestIsStringInList(t *testing.T) {
	testCases := []struct {
		name          string
		value         string
		list          []string
		expectedValue bool
	}{
		{
			name:          "EmptyList",
			value:         "apple",
			list:          []string{},
			expectedValue: false,
		},
		{
			name:          "StringInList",
			value:         "apple",
			list:          []string{"banana", "apple", "orange"},
			expectedValue: true,
		},
		{
			name:          "StringNotInList",
			value:         "grape",
			list:          []string{"banana", "apple", "orange"},
			expectedValue: false,
		},
		{
			name:          "SubstringInList",
			value:         "app",
			list:          []string{"banana", "apple", "orange"},
			expectedValue: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsStringInList(tc.value, tc.list)
			if result != tc.expectedValue {
				t.Errorf("Expected %v, but got %v", tc.expectedValue, result)
			}
		})
	}
}

func TestLookupEnvString(t *testing.T) {
	os.Setenv("EXISTING_VAR", "existing_value")

	testCases := []struct {
		name          string
		key           string
		defaultValue  string
		expectedValue string
	}{
		{
			name:          "ExistingEnvVar",
			key:           "EXISTING_VAR",
			defaultValue:  "default_value",
			expectedValue: "existing_value",
		},
		{
			name:          "NonExistingEnvVar",
			key:           "NON_EXISTING_VAR",
			defaultValue:  "default_value",
			expectedValue: "default_value",
		},
		{
			name:          "EmptyEnvVar",
			key:           "EMPTY_VAR",
			defaultValue:  "default_value",
			expectedValue: "default_value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := LookupEnvString(tc.key, tc.defaultValue)
			if result != tc.expectedValue {
				t.Errorf("Expected %s, but got %s", tc.expectedValue, result)
			}
		})
	}
}

func TestCheckEnvVars(t *testing.T) {
	os.Setenv("EXISTING_VAR", "existing_value")
	testCases := []struct {
		name          string
		key           string
		envVar        string
		defaultValue  string
		expectedValue string
	}{
		{
			name:          "KeyProvided",
			key:           "provided_key",
			envVar:        "EXISTING_VAR",
			defaultValue:  "default_value",
			expectedValue: "provided_key",
		},
		{
			name:          "KeyNotProvided",
			key:           "",
			envVar:        "EXISTING_VAR",
			defaultValue:  "default_value",
			expectedValue: "existing_value",
		},
		{
			name:          "KeyNotProvidedEmptyEnvVar",
			key:           "",
			envVar:        "EMPTY_VAR",
			defaultValue:  "default_value",
			expectedValue: "default_value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CheckEnvVars(tc.key, tc.envVar, tc.defaultValue)
			if result != tc.expectedValue {
				t.Errorf("Expected %s, but got %s", tc.expectedValue, result)
			}
		})
	}
}

func TestGetFlagOrEnv(t *testing.T) {
	os.Setenv("EXISTING_VAR", "existing_value")

	rootCmd := &cobra.Command{Use: "root"}
	rootCmd.PersistentFlags().String("testflag", "", "Test flag")

	testCases := []struct {
		name          string
		flagValue     string
		envVar        string
		defaultValue  string
		expectedValue string
	}{
		{
			name:          "FlagProvided",
			flagValue:     "flag_value",
			envVar:        "EXISTING_VAR",
			defaultValue:  "default_value",
			expectedValue: "flag_value",
		},
		{
			name:          "FlagNotProvided",
			flagValue:     "",
			envVar:        "EXISTING_VAR",
			defaultValue:  "default_value",
			expectedValue: "existing_value",
		},
		{
			name:          "FlagNotProvidedEmptyEnvVar",
			flagValue:     "",
			envVar:        "EMPTY_VAR",
			defaultValue:  "default_value",
			expectedValue: "default_value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rootCmd.SetArgs([]string{"--testflag=" + tc.flagValue})
			_ = rootCmd.Execute()

			result := GetFlagOrEnv(rootCmd, "testflag", tc.envVar, tc.defaultValue)
			if result != tc.expectedValue {
				t.Errorf("Expected %s, but got %s", tc.expectedValue, result)
			}
		})
	}
}

func TestWriteFile(t *testing.T) {
	testContent := "Hello, world!"
	testFilePath := "testfile.txt"

	// Cleanup the test file after the test
	defer func() {
		_ = os.Remove(testFilePath)
	}()

	// Test writing to a valid file path
	err := WriteFile(testContent, testFilePath)
	assert.NoError(t, err, "Expected no error")

	// Read the content of the written file
	contentBytes, err := os.ReadFile(testFilePath)
	assert.NoError(t, err, "Expected no error")

	assert.Equal(t, testContent, string(contentBytes), "File content does not match")

	// Test writing to an invalid file path
	invalidFilePath := "/nonexistentfolder/testfile.txt"
	err = WriteFile(testContent, invalidFilePath)
	assert.Error(t, err, "Expected an error")

	// Test writing with a permission-denied scenario (simulate by creating a read-only file)
	readOnlyFilePath := "readonlyfile.txt"
	_ = os.WriteFile(readOnlyFilePath, []byte("initial content"), 0444) // Create a read-only file
	defer func() {
		_ = os.Remove(readOnlyFilePath)
	}()

	err = WriteFile(testContent, readOnlyFilePath)
	fmt.Println(err.Error())
	assert.Error(t, err, "Expected an error")
}

func TestCheckRequiredArgs(t *testing.T) {
	type TestArgs struct {
		Field1 string
		Field2 string
		Field3 string
	}

	testCases := []struct {
		name     string
		args     interface{}
		expected error
	}{
		{
			name: "All fields provided",
			args: TestArgs{
				Field1: "value1",
				Field2: "value2",
				Field3: "value3",
			},
			expected: nil,
		},
		{
			name: "One missing field",
			args: TestArgs{
				Field1: "value1",
				Field2: "",
				Field3: "value3",
			},
			expected: fmt.Errorf("Missing required parameters: Field2"),
		},
		{
			name: "Multiple missing fields",
			args: TestArgs{
				Field1: "",
				Field2: "",
				Field3: "",
			},
			expected: fmt.Errorf("Missing required parameters: Field1, Field2, Field3"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := CheckRequiredArgs(tc.args)
			assert.Equal(t, tc.expected, err)
		})
	}
}
