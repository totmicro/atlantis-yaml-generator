package version

import "testing"

func TestGetVersion(t *testing.T) {
	expectedVersion := "0.0.1"
	actualVersion := GetVersion()
	if actualVersion != expectedVersion {
		t.Errorf("Expected version %s, but got %s", expectedVersion, actualVersion)
	}
}
