package structure

import (
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

func TestDirectory_Equals_WithDifferentInstances(t *testing.T) {
	directory1 := Directory{Name: "dir1", Path: "/tmp"}
	directory2 := Directory{Name: "dir1", Path: "/tmp"}
	if !directory1.Equals(&directory2) {
		t.Fatal("directories were found to be equal but were not")
	}
}

func TestDirectory_Equals_WhenDifferentName(t *testing.T) {
	directory1 := Directory{Name: "dir1", Path: "/tmp"}
	directory2 := Directory{Name: "dir2", Path: "/tmp"}
	if directory1.Equals(&directory2) {
		t.Fatal("directories were found to be equal but were not")
	}
}

func TestDirectory_Equals_WhenDifferentPath(t *testing.T) {
	directory1 := Directory{Name: "dir1", Path: "/tmp"}
	directory2 := Directory{Name: "dir1", Path: "/tmp/dir"}
	if directory1.Equals(&directory2) {
		t.Fatal("directories were found to be equal but were not")
	}
}

func TestDirectory_AddDirectory(t *testing.T) {
	for _, tt := range DirectoryIdentities {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := tt.dir.GetDirectory(tt.dirFullPath)
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

func TestDirectory_AddDirectory_CreatesSubdirectoriesAsNecessary(t *testing.T) {
	dir := Directory{Name: "dir", Path: "/tmp"}
	_, err := dir.AddDirectory("/tmp/dir/subdir1/subdir2/subdir3")
	if err != nil {
		t.Fatal(err)
	}
	if subdir1, ok := dir.SubDirectories["subdir1"]; !ok {
		t.Fatal("subdir1 was not created")
		expected := Directory{Name: "subdir1", Path: "/tmp/dir"}
		if !subdir1.Equals(&expected) {
			t.Fatal("subdir1 structure was incorrect")
		}
	}
	if subdir2, ok := dir.SubDirectories["subdir1"].SubDirectories["subdir2"]; !ok {
		t.Fatal("subdir2 was not created")
		expected := Directory{Name: "subdir2", Path: "/tmp/dir/subdir1"}
		if !subdir2.Equals(&expected) {
			t.Fatal("subdir2 structure was incorrect")
		}
	}
	if subdir3, ok := dir.SubDirectories["subdir1"].SubDirectories["subdir2"].SubDirectories["subdir3"]; !ok {
		t.Fatal("subdir1 was not created")
		expected := Directory{Name: "subdir3", Path: "/tmp/dir/subdir1/subdir2"}
		if !subdir3.Equals(&expected) {
			t.Fatal("subdir1 structure was incorrect")
		}
	}
}

func TestDirectory_AddDirectory_ReturnsErrorIfNotSubdirectory(t *testing.T) {
	dir := Directory{Name: "dir", Path: "/tmp"}
	_, err := dir.AddDirectory("/tmp/other/subdir1/subdir2/subdir3")
	if err == nil {
		t.Fatal("error should have been returned but was nil")
	}
}

func TestDirectory_AddDirectory_ReturnsErrorIfDifferentParent(t *testing.T) {
	dir := Directory{Name: "dir", Path: "/tmp"}
	_, err := dir.AddDirectory("/other/dir/subdir1/subdir2/subdir3")
	if err == nil {
		t.Fatal("error should have been returned but was nil")
	}
}

func TestDirectory_GetDirectory(t *testing.T) {
	for _, tt := range FindTests {
		t.Run(tt.name, func(t *testing.T) {
			expectedPath, expectedName := filepath.Split(tt.fullDirPathToFind)
			expectedPath = filepath.Clean(expectedPath)
			found, err := tt.dir.GetDirectory(tt.fullDirPathToFind)
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

func TestDirectory_FindDirectoryDepth(t *testing.T) {
	for _, tt := range FindTests {
		t.Run(tt.name, func(t *testing.T) {
			expectedPath, expectedName := filepath.Split(tt.fullDirPathToFind)
			expectedPath = filepath.Clean(expectedPath)
			found := tt.dir.FindDirectoryDepth(expectedName)
			if found == nil {
				t.Fatal("nil was returned but actual directory was expected")
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
