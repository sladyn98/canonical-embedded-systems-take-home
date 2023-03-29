package main

import (
	"crypto/rand"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// Test helper function to create a file with random content
func createTempFileWithRandomContent(t *testing.T, dir string, size int64) string {
	file, err := ioutil.TempFile(dir, "testfile-")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	data := make([]byte, size)
	_, err = rand.Read(data)
	if err != nil {
		t.Fatal(err)
	}

	_, err = file.Write(data)
	if err != nil {
		t.Fatal(err)
	}

	return file.Name()
}

func TestShred(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "shred-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	testCases := []struct {
		name        string
		setup       func() string
		expectError bool
	}{
		{
			name: "Small text file",
			setup: func() string {
				return createTempFileWithRandomContent(t, tempDir, 1024)
			},
			expectError: false,
		},
		{
			name: "Large text file",
			setup: func() string {
				return createTempFileWithRandomContent(t, tempDir, 1024*1024)
			},
			expectError: false,
		},
		{
			name: "Small binary file",
			setup: func() string {
				return createTempFileWithRandomContent(t, tempDir, 1024)
			},
			expectError: false,
		},
		{
			name: "Large binary file",
			setup: func() string {
				return createTempFileWithRandomContent(t, tempDir, 1024*1024)
			},
			expectError: false,
		},
		{
			name: "Empty file",
			setup: func() string {
				file, err := ioutil.TempFile(tempDir, "testfile-empty-")
				if err != nil {
					t.Fatal(err)
				}
				file.Close()
				return file.Name()
			},
			expectError: false,
		},
		{
			name:        "Non-existent file",
			setup:       func() string { return filepath.Join(tempDir, "nonexistent-file") },
			expectError: true,
		},
		{
			name: "Read-only file",
			setup: func() string {
				filePath := createTempFileWithRandomContent(t, tempDir, 1024)
				err := os.Chmod(filePath, 0400)
				if err != nil {
					t.Fatal(err)
				}
				return filePath
			},
			expectError: true,
		},
		{
			name: "Directory",
			setup: func() string {
				dir, err := ioutil.TempDir(tempDir, "testdir-")
				if err != nil {
					t.Fatal(err)
				}
				return dir
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := tc.setup()
			err := Shred(path)
			if tc.expectError {
				if err == nil {
					t.Errorf("expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, but got: %v", err)
				} else {
					if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
						t.Errorf("expected the file to be deleted, but it still exists")
					}
				}
			}
		})
	}
}

			
