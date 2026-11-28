package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// TEST FUNCTIONS

// Creates a temporary file and uses the temp directory path to test dif path types
// Defines a table of test cases covering a real file, a direcory and a non-existent file
// Checks for each case that only existing regular file returns true
func TestCheckFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	tempFile := filepath.Join(tmpDir, "sample.txt")
	if err := os.WriteFile(tempFile, []byte("hello"), 0o600); err != nil {
		t.Fatal(err)
	}

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
		{
			name:        "Path through a file (not a directory) should return false",
			path:        filepath.Join(tempFile, "child"),
			expectExist: false,
		},
		{
			name:        "Empty path should return false",
			path:        "",
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

// Creates a temporary test directory (which is deleted when test finsihes)
// Defines a table of cases: some files present, some missing, empty list, etc.
// Checks that only the existing files come back, in the same order as the input
func TestLoadExistingImages(t *testing.T) {
	dir := t.TempDir()
	keep := filepath.Join(dir, "keep.png")
	other := filepath.Join(dir, "other.png")
	for _, p := range []string{keep, other} {
		if err := os.WriteFile(p, []byte("img"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "folder.png"), 0o700); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name  string
		names []string
		want  []string
	}{
		{
			name:  "Missing file is ignored",
			names: []string{"keep.png", "missing.png"},
			want:  []string{keep},
		},
		{
			name:  "All files exist",
			names: []string{"keep.png", "other.png"},
			want:  []string{keep, other},
		},
		{
			name:  "Order is preserved",
			names: []string{"other.png", "keep.png"},
			want:  []string{other, keep},
		},
		{
			name:  "Folder with an image name is ignored",
			names: []string{"folder.png", "keep.png"},
			want:  []string{keep},
		},
		{
			name:  "All files missing",
			names: []string{"a.png", "b.png"},
			want:  []string{},
		},
		{
			name:  "Empty list",
			names: []string{},
			want:  []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := loadExistingImages(dir, tc.names)

			if len(got) != len(tc.want) {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("index %d: got %s, want %s", i, got[i], tc.want[i])
				}
			}
		})
	}
}

// Defines test cases for valid or empty containing image names
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
			name:     "Name with spaces is valid",
			image:    "my image",
			expected: true,
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

// Creates a blank PNG image of the given width and height and returns it as bytes
// t.Helper() makes failures point at the test that called this, not at the helper
func makePNG(t *testing.T, w, h int) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := png.Encode(&buf, image.NewRGBA(image.Rect(0, 0, w, h))); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func makeJPEG(t *testing.T, w, h int) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, image.NewRGBA(image.Rect(0, 0, w, h)), nil); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// Builds different kinds of input (good and bad images) and puts them in a table
// Calls ValidateImage for each one as a subtest
// Checks that only empty, fake and too small images return an error
func TestValidateImage(t *testing.T) {
	testCases := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		// valid images
		{
			name:    "Normal image is valid",
			data:    makePNG(t, 100, 100),
			wantErr: false,
		},
		{
			name:    "Original pipeline size 1600x768 is valid",
			data:    makePNG(t, 1600, 768),
			wantErr: false,
		},
		{
			name:    "Exactly the minimum size is valid",
			data:    makePNG(t, 64, 64),
			wantErr: false,
		},

		// broken input
		{
			name:    "Empty data is invalid",
			data:    []byte{},
			wantErr: true,
		},
		{
			name:    "Nil data is invalid",
			data:    nil,
			wantErr: true,
		},
		{
			name:    "Text pretending to be an image is invalid",
			data:    []byte("this is not a png"),
			wantErr: true,
		},

		// too small (one pixel under the minimum, on each side)
		{
			name:    "Width below minimum is invalid",
			data:    makePNG(t, 63, 100),
			wantErr: true,
		},
		{
			name:    "Height below minimum is invalid",
			data:    makePNG(t, 100, 63),
			wantErr: true,
		},

		{
			name:    "Normal JPEG is valid",
			data:    makeJPEG(t, 100, 100),
			wantErr: false,
		},
		{
			name:    "JPEG below minimum size is invalid",
			data:    makeJPEG(t, 63, 100),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateImage(tc.data)

			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateImage() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

// Defines test cases for APP_ENV being set and being empty
// Uses t.Setenv for each subtest so environment changes stay isolated and clean up automatically
// Checks that APP_ENV is returned as-is or defaults to "development"
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
