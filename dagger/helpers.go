package main

import (
	"os"
	"path/filepath"
)

// HELPER FUNCTIONS

// Takes an image folder path (imgDir) and a list of expected filenames
// Builds the full path for each filename using filepaht.Join
// Checks with os.Stat wheter that file exists on disk
// Keeps only the path of file that exist and skips the missing ones
// Returns the existing file paths
func loadExistingImages(imgDir string, filenames []string) []string {
	var found []string
	for _, name := range filenames {
		p := filepath.Join(imgDir, name)
		if _, err := os.Stat(p); err == nil {
			found = append(found, p)
		}
	}
	return found
}

// Checks that an image name is non-empty and onctains no spaces
// Returns false for empty strings or names with spaces and true otherwise
func isValidImageName(name string) bool {
	if name == "" {
		return false
	}
	for _, char := range name {
		if char == ' ' {
			return false
		}
	}
	return true
}

// Checks wether filePath points to an existing file and not a missing path or directory
// Returns true only when the file exists and is not a directory
func checkFileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir() //returns true if it exists and is a file
}

// Reads the APP_ENV environment variable and returns it
// Falling back to "deployment" when APP_ENV is unset or empty
func getAppEnvironment() string {
	env := os.Getenv("APP_ENV")

	//if nothing was set, default to development
	if env == "" {
		return "development"
	}
	return env
}
