package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
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
		if checkFileExists(p) {
			found = append(found, p)
		}
	}
	return found
}

// Checks that an image name is not empty
// Returns false for empty strings and true otherwise
func isValidImageName(name string) bool {
	return name != ""
}

// Checks wether filePath points to an existing file and not a missing path or directory
// Any error from os.Stat (missing, permission denied, path through a file) counts as "no"
// Returns true only when the file exists and is not a directory
func checkFileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// Reads the APP_ENV environment variable and returns it
// Falling back to "development" when APP_ENV is unset or empty
func getAppEnvironment() string {
	env := os.Getenv("APP_ENV")

	//if nothing was set, default to development
	if env == "" {
		return "development"
	}
	return env
}

const minImageSize = 64 // width and height must be at least this (the model input is 64x64)

var ErrEmptyImage = errors.New("image is empty")

// Checks that image data is safe to send on to the Python pipeline
// Rejects: empty data, not a real PNG/JPEG, too small
// Accepts everything else, because Python resizes to 64x64 anyway
// Returns nil when the image is fine, otherwise an error that says what is wrong
func ValidateImage(data []byte) error {
	// nothing to check if there is no data at all
	if len(data) == 0 {
		return ErrEmptyImage
	}

	// DecodeConfig only reads the header (cheap), so it tells us the type and size
	// without decoding the whole image. It fails if the data is not a real image
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("not a valid image: %w", err)
	}
	if cfg.Width < minImageSize || cfg.Height < minImageSize {
		return fmt.Errorf("image too small: %dx%d, min %dx%d", cfg.Width, cfg.Height, minImageSize, minImageSize)
	}
	return nil
}
