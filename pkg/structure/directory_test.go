package structure

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestDirectory_Equals_WithIdentity(t *testing.T) {
	for _, tt := range DirectoryIdentities {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.dir.Equals(&tt.dir) {
				t.Fatal("directories were equal but were not found to be")
			}
		})
	}
}

func TestDirectory_Equals_WhenNotEqual(t *testing.T) {
	for _, tt := range DirectoryIdentities {
		for _, ott := range DirectoryIdentities {
			if tt.name == ott.name {
				continue
			}
			t.Run(fmt.Sprintf("%s_And_%s", tt.name, ott.name), func(t *testing.T) {
				if tt.dir.Equals(&ott.dir) {
					tt.dir.Equals(&ott.dir)
					t.Fatal("directories were found to be equal but were not")
				}
			})
		}
	}
}

func TestDirectory_AddDirectory(t *testing.T) {
	for _, tt := range DirectoryIdentities {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := tt.dir.FindDirectory(tt.dirFullPath)
			if err != nil {
				t.Fatal(err)
			}
			_, err = dir.AddDirectory(filepath.Join(tt.dirFullPath, "newdir"))
			if err != nil {
				t.Fatal(err)
			}
			if subDir, ok := dir.SubDirectories["newdir"]; ok {
				if subDir.Name != "newdir" {
					t.Fatalf("file name was not set correctly. expected %s but was %s", "newdir", subDir.Name)
				}
				if subDir.Path != tt.dirFullPath {
					t.Fatalf("file path was not set correctly. expected %s but was %s", tt.dirFullPath, subDir.Path)
				}
			} else {
				t.Fatalf("subdirectory was not found in the subdirectories")
			}
		})
	}
}

func TestDirectory_AddDirectory_ReturnsError(t *testing.T) {
	dir := Directory{Name: "dir", Path: "/tmp"}
	_, err := dir.AddDirectory("/tmp/dir/subdir1/subdir2/subdir3")
	if err == nil {
		t.Fatal("error should have been returned but was nil")
	}
}

func TestDirectory_AddFile_ReturnsError(t *testing.T) {
	dir := Directory{Name: "dir", Path: "/tmp"}
	_, err := dir.AddDirectory("/tmp/dir/subdir1/subdir2/subdir3/file.txt")
	if err == nil {
		t.Fatal("error should have been returned but was nil")
	}
}

func TestDirectory_FindDirectory(t *testing.T) {
	for _, tt := range FindTests {
		t.Run(tt.name, func(t *testing.T) {
			expectedPath, expectedName := filepath.Split(tt.fullDirPathToFind)
			expectedPath = filepath.Clean(expectedPath)
			found, err := tt.dir.FindDirectory(tt.fullDirPathToFind)
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
