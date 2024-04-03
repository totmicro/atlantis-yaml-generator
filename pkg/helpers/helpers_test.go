package helpers

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestLookupEnvString(t *testing.T) {
	// Set up environment variables for testing
	key := "TEST_KEY"
	value := "test_value"
	err := os.Setenv(key, value)
	assert.NoError(t, err, "Failed to set up environment variable for testing")
	defer os.Unsetenv(key)
	// Test when the environment variable exists
	result := LookupEnvString(key)
	assert.Equal(t, value, result, "Returned value doesn't match expected")
	// Test when the environment variable doesn't exist
	missingKey := "MISSING_KEY"
	missingResult := LookupEnvString(missingKey)
	assert.Equal(t, "", missingResult, "Expected an empty string for missing key")
}

func TestReadFile(t *testing.T) {
	// Create a temporary test file with content
	tempFile, err := os.CreateTemp("", "testfile*.txt")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
	// Write content to the test file
	content := "This is a test file."
	_, err = tempFile.WriteString(content)
	assert.NoError(t, err)
	// Call the function and test the logic
	readContent, err := ReadFile(tempFile.Name())
	assert.NoError(t, err)                // Check that there are no errors
	assert.Equal(t, content, readContent) // Check that the read content matches the original content
}

func TestReadFile_OpenError(t *testing.T) {
	// Call the function with a non-existent file
	_, err := ReadFile("nonexistent.txt")
	assert.Error(t, err) // Check that an error is returned
}

func TestReadFile_StatError(t *testing.T) {
	// Create a temporary test file with content
	tempFile, err := os.CreateTemp("", "testfile*.txt")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()
	// Remove read permissions from the file to simulate a Stat error
	err = os.Chmod(tempFile.Name(), 0222)
	assert.NoError(t, err)
	// Call the function and test the logic
	_, err = ReadFile(tempFile.Name())
	assert.Error(t, err) // Check that an error is returned
}

func TestReadFile_ReadError(t *testing.T) {
	// Create a temporary test file with content
	tempFile, err := os.CreateTemp("", "testfile*.txt")
	assert.NoError(t, err)
	tempFile.Close()
	os.Remove(tempFile.Name())
	// Call the function and test the logic
	_, err = ReadFile(tempFile.Name())
	assert.Error(t, err) // Check that an error is returned
}

func TestMatchesPattern(t *testing.T) {
	tests := []struct {
		pattern string
		str     string
		want    bool
	}{
		{"workspace_vars", "workspace_vars", true},
		{"main.tf", "main.tf", true},
		{".*\\.tf", "main.tf", true},
		{".*\\.tf", "vpc.tf", true},
		{"^a.*z$", "alphabet", false}, // Negative test case
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			if got := MatchesPattern(tt.pattern, tt.str); got != tt.want {
				t.Errorf("MatchesPattern(%q, %q) = %v, want %v", tt.pattern, tt.str, got, tt.want)
			}
		})
	}
}

// HashableString is a simple Hashable type for testing.
type HashableString string

func (h HashableString) Hash() string {
	return string(h)
}

func TestNewSet(t *testing.T) {
	s := NewSet()
	if s == nil || len(s.Elements) != 0 {
		t.Errorf("NewSet() = %v, want a new Set instance with empty Elements", s)
	}
}

func TestAddAndContains(t *testing.T) {
	s := NewSet()
	element := HashableString("test")
	s.Add(element)
	if !s.Contains(element) {
		t.Errorf("Set does not contain element %v after Add", element)
	}
}

func TestRemove(t *testing.T) {
	s := NewSet()
	element := HashableString("test")
	s.Add(element)
	s.Remove(element)
	if s.Contains(element) {
		t.Errorf("Set contains element %v after Remove", element)
	}
}

func TestSize(t *testing.T) {
	s := NewSet()
	s.Add(HashableString("one"))
	s.Add(HashableString("two"))
	if s.Size() != 2 {
		t.Errorf("Size() = %d, want 2", s.Size())
	}
}

func TestList(t *testing.T) {
	s := NewSet()
	elements := []Hashable{HashableString("one"), HashableString("two")}
	for _, e := range elements {
		s.Add(e)
	}
	list := s.List()
	if !reflect.DeepEqual(list, elements) && len(list) == len(elements) {
		t.Errorf("List() = %v, want %v", list, elements)
	}
}
