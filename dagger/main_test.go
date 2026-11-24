package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TEST FUNCTIONS

// Creates a temporary test directory (which is deleted when test finsihes)
// Creates one real dummy file (keep.png) inside that temp folder and leaves missing.png out
// Calls to check that the missing img is ignored and only keep.png is returned
func TestLoadExistingImages(t *testing.T) {
	dir := t.TempDir()
	existing := filepath.Join(dir, "keep.png")
	_ = os.WriteFile(existing, []byte("img"), 0o600)

	got := loadExistingImages(dir, []string{"keep.png", "missing.png"})
	if len(got) != 1 || got[0] != existing {
		t.Fatalf("got %v, want [%s]", got, existing)
	}
}

// Defines test cases for different APP_ENV values using a table-driven test
// Uses t.Setenv for each subtest so environment changes stay isolated and clean up automatically
// Checks that APP_ENV matches the expected mode or defaults to "development"
func TestGetEnvironmentMode(t *testing.T) {
	testCases := []struct {
		name         string
		envValue     string
		expectedMode string
	}{
		{
			name:         "Production mode flag set",
			envValue:     "production",
			expectedMode: "production",
		},
		{
			name:         "Empty flag defaults to development",
			envValue:     "",
			expectedMode: "development",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("APP_ENV", tc.envValue)

			mode := getAppEnvironment()

			if mode != tc.expectedMode {
				t.Errorf("expected mode %s, got %s", tc.expectedMode, mode)
			}

		})
	}
}

// Defines test cases for valid, empty and space containing image names
// Loops through each case as subtest and calls isValidImageName
// Checks that the returned boolean matches wether the image name should be valid or not
func TestIsValidImageName(t *testing.T) {
	testCases := []struct {
		name     string
		image    string
		expected bool
	}{
		{
			name:     "Valid image name",
			image:    "my-image",
			expected: true,
		},
		{
			name:     "Empty name is invalid",
			image:    "",
			expected: false,
		},
		{
			name:     "Name with spaces is invalid",
			image:    "my image",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := isValidImageName(tc.image)
			if got != tc.expected {
				t.Errorf("isValidImageName(%q) = %v, want %v", tc.image, got, tc.expected)
			}
		})
	}
}

// Creates a temporary file and uses the temp directory path to test dif path types
// Defines a table of test cases covering a real file, a direcory and a non-existent file
// Checks for each case that only existing regular file returns true
func TestCheckFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	tempFile := filepath.Join(tmpDir, "sample.txt")
	_ = os.WriteFile(tempFile, []byte("hello"), 0o600)

	testCases := []struct {
		name        string
		path        string
		expectExist bool
	}{
		{
			name:        "Existing regular file should return true",
			path:        tempFile,
			expectExist: true,
		},
		{
			name:        "Directory should return false",
			path:        tmpDir,
			expectExist: false,
		},
		{
			name:        "Missing file should return false",
			path:        filepath.Join(tmpDir, "non_existent_file.txt"),
			expectExist: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			exists := checkFileExists(tc.path)

			if exists != tc.expectExist {
				t.Errorf("expected file existence to be %v, got %v", tc.expectExist, exists)
			}
		})
	}
}
