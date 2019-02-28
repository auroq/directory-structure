package structure

import (
	"path/filepath"
	"testing"
)

func TestDirectory_AddFile(t *testing.T) {
	for _, tt := range DirectoryIdentities {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := tt.dir.FindDirectory(tt.dirFullPath)
			if err != nil {
				t.Fatal(err)
			}
			_, err = dir.AddFile(filepath.Join(tt.dirFullPath, "newfile"))
			if err != nil {
				t.Fatal(err)
			}
			if file, ok := tt.dir.Files["file1.txt"]; ok {
				if file.Name != "file1.txt" {
					t.Fatalf("file name was not set correctly. expected %s but was %s", "file1.txt", file.Name)
				}
				if file.Path != tt.dirFullPath {
					t.Fatalf("file path was not set correctly. expected %s but was %s", tt.dirFullPath, file.Path)
				}
			}
		})
	}
}

func TestDirectory_FindFile(t *testing.T) {
	for _, tt := range FindTests {
		t.Run(tt.name, func(t *testing.T) {
			expectedPath, expectedName := filepath.Split(tt.fullFilePathToFind)
			expectedPath = filepath.Clean(expectedPath)
			found, err := tt.dir.FindFile(tt.fullFilePathToFind)
			if err != nil {
				t.Fatal(err)
			}
			if found.Path != expectedPath {
				t.Fatalf("found path did not match expected. expected path: "+
					"'%s' actual path: '%s'", expectedPath, found.Path)
			}
			if found.Name != expectedName {
				t.Fatalf("found name did not match expected. expected name: "+
					"'%s' actual name: '%s'", expectedName, found.Name)
			}
		})
	}
}
