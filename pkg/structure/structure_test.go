package structure

import (
	"path/filepath"
	"testing"
)

var DirectoryIdentities = []struct {
	name string
	dir  Directory
}{
	{"EmptyDirectory",
		Directory{Name: "/tmp/dir1", Path: "dir1"},
	},
	{"DirectoryWithSubDirectory",
		Directory{
			Name: "/tmp/dir1",
			Path: "dir1",
			SubDirectories: map[string]Directory{
				"sub1": {
					Name: "sub1",
					Path: "/tmp/dir1/sub1",
				},
			},
		},
	},
	{"DirectoryWithSubDirectoryWithSubDirectory",
		Directory{
			Name: "/tmp/dir1",
			Path: "dir1",
			SubDirectories: map[string]Directory{
				"sub1": {
					Name: "sub1",
					Path: "/tmp/dir1/sub1",
					SubDirectories: map[string]Directory{
						"subsub1": {
							Name: "subsub1",
							Path: "/tmp/dir1/sub1/subsub1",
						},
					},
				},
			},
		},
	},
}

func TestDirectory_Equals_WithIdentity(t *testing.T) {
	for _, tt := range DirectoryIdentities {
		t.Run(tt.name, func(t *testing.T) {
			tt.dir.Equals(tt.dir)
		})
	}
}

var DirectoryTests = []struct {
	name string
	dir  Directory
}{
	{"EmptyDirectory", Directory{Name: "dir", Path: "/tmp"}},
	{"DirectoryWithFile", Directory{Name: "dir", Path: "/tmp", Files: map[string]File{"File1": {Name: "File1", Path: "/tmp/dir"}}}},
	{"DirectoryWithSubDirectory", Directory{Name: "dir", Path: "/tmp", SubDirectories: map[string]Directory{"subdir1": {Name: "subdir1", Path: "/tmp/dir"}}}},
}

func TestDirectory_AddFile(t *testing.T) {
	for _, tt := range DirectoryTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.dir.AddFile("/tmp/dir/file1.txt")
			if err != nil {
				t.Fatal(err)
			}
			if file, ok := tt.dir.Files["file1.txt"]; ok {
				if file.Name != "file1.txt" {
					t.Fatalf("file name was not set correctly. expected %s but was %s", "file1.txt", file.Name)
				}
				if file.Path != "/tmp/dir" {
					t.Fatalf("file path was not set correctly. expected %s but was %s", "tmp/dir", file.Path)
				}
			}
		})
	}
}

func TestDirectory_AddDirectory(t *testing.T) {
	for _, tt := range DirectoryTests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.dir.AddDirectory("/tmp/dir/subdir")
			if err != nil {
				t.Fatal(err)
			}
			if subDir, ok := tt.dir.SubDirectories["subdir"]; ok {
				if subDir.Name != "subdir" {
					t.Fatalf("file name was not set correctly. expected %s but was %s", "subdir", subDir.Name)
				}
				if subDir.Path != "/tmp/dir" {
					t.Fatalf("file path was not set correctly. expected %s but was %s", "tmp/dir", subDir.Path)
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

var FindTests = []struct {
	name           string
	dir            Directory
	fullPathToFind string
}{
	{"DirectoryWithSubDirectory",
		func() Directory {
			dir := Directory{Name: "dir1", Path: "/tmp"}
			_, _ = dir.AddDirectory("/tmp/dir1/sub1")
			_, _ = dir.AddFile("/tmp/dir1/sub1.txt")
			return dir
		}(),
		"/tmp/dir1/sub1",
	},
	{"DirectoryWithSubDirectoryWithSubDirectory",
		func() Directory {
			dir := Directory{Name: "dir1", Path: "/tmp"}
			sub1, _ := dir.AddDirectory("/tmp/dir1/sub1")
			_, _ = sub1.AddDirectory("/tmp/dir1/sub1/subsub1")
			_, _ = sub1.AddFile("/tmp/dir1/sub1/subsub1.txt")
			return dir
		}(),
		"/tmp/dir1/sub1/subsub1",
	},
	{"DirectoryWithSubDirectories",
		func() Directory {
			dir := Directory{Name: "dir1", Path: "/tmp"}
			sub1, _ := dir.AddDirectory("/tmp/dir1/sub1")
			_, _ = sub1.AddDirectory("/tmp/dir1/sub1/subsub1")
			subsub2, _ := sub1.AddDirectory("/tmp/dir1/sub1/subsub2")
			_, _ = sub1.AddDirectory("/tmp/dir1/sub1/subsub3")
			_, _ = subsub2.AddDirectory("/tmp/dir1/sub1/subsub2/subsubsub")
			_, _ = subsub2.AddFile("/tmp/dir1/sub1/subsub2/subsubsub.txt")
			return dir
		}(),
		"/tmp/dir1/sub1/subsub2/subsubsub",
	},
}

func TestDirectory_FindDirectory(t *testing.T) {
	for _, tt := range FindTests {
		t.Run(tt.name, func(t *testing.T) {
			expectedPath, expectedName := filepath.Split(tt.fullPathToFind)
			found, err := tt.dir.Find(tt.fullPathToFind)
			if err != nil {
				t.Fatal(err)
			}
			if found.Path != expectedPath {
				t.Fatalf("found path did not match expected. expected path: " +
					"'%s' actual path: '%s'", expectedPath, found.Path)
			}
			if found.Name != expectedName {
				t.Fatalf("found name did not match expected. expected name: " +
					"'%s' actual name: '%s'", expectedName, found.Name)
			}
		})
	}
}
