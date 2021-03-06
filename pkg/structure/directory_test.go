package structure

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDirectory_Path_CleansName(t *testing.T) {
	for _, tt := range cleanTests {
		t.Run(tt.name, func(t *testing.T) {
			dir := NewDirectory("dir1", tt.path)
			if dir.Path() != tt.cleanPath {
				t.Fatalf("path was not properly cleaned expected: '%s' actual: '%s'", tt.cleanPath, dir.Path())
			}
		})
	}
}

func TestDirectory_FullPath(t *testing.T) {
	for _, tt := range fullPathTests {
		t.Run(tt.testName, func(t *testing.T) {
			dir := NewDirectory(tt.name, tt.path)
			if dir.FullPath() != tt.fullPath {
				t.Fatalf("full path did not match expected: '%s' actual: '%s'", tt.fullPath, dir.FullPath())
			}
		})
	}
}

func TestDirectory_SubDirectory_ReturnsSubDir(t *testing.T) {
	dir := NewDirectory("dir1", filepath.Join(osRoot(), "tmp"))
	sub1, err := dir.AddDirectory(filepath.Join(osRoot(), "tmp", "dir1", "sub1"))
	if err != nil {
		t.Fatal(err)
	}
	if subDir := dir.SubDirectory("sub1"); !subDir.Equals(sub1) {
		t.Fatal("directory returned did not match")
	}
}

func TestDirectory_SubDirectory_NilWhenSubDirectoriesNil(t *testing.T) {
	dir := NewDirectory("dir1", filepath.Join(osRoot(), "tmp"))
	if subDir := dir.SubDirectory("nonexistent"); subDir != nil {
		t.Fatal("directory did not exist and should have been nil")
	}
}

func TestDirectory_SubDirectory_NilWhenNotFound(t *testing.T) {
	dir := NewDirectory("dir1", filepath.Join(osRoot(), "tmp"))
	_, err := dir.AddDirectory(filepath.Join(osRoot(), "tmp", "dir1", "sub1"))
	if err != nil {
		t.Fatal(err)
	}
	if subDir := dir.SubDirectory("nonexistent"); subDir != nil {
		t.Fatal("directory did not exist and should have been nil")
	}
}

func TestDirectory_Equals_TrueWithIdentity(t *testing.T) {
	for _, tt := range DirectoryIdentities {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.dir.Equals(&tt.dir) {
				t.Fatal("directories were equal but were not found to be")
			}
		})
	}
}

func TestDirectory_Equals_TrueWithDifferentInstances(t *testing.T) {
	directory1 := NewDirectory("dir1", filepath.Join(osRoot(), "tmp"))
	directory2 := NewDirectory("dir1", filepath.Join(osRoot(), "tmp"))
	if !directory1.Equals(directory2) {
		t.Fatal("directories were found to be equal but were not")
	}
}

func TestDirectory_Equals_TrueWhenPathStringNotClean(t *testing.T) {
	directory1 := NewDirectory("dir1", filepath.Join(osRoot(), "tmp")+string(os.PathSeparator))
	directory2 := NewDirectory("dir1", filepath.Join(osRoot(), "tmp"))
	if !directory1.Equals(directory2) {
		t.Fatal("directories were equal but were not found to be")
	}
}

func TestDirectory_Equals_FalseWhenPathStringNotSameCase(t *testing.T) {
	directory1 := NewDirectory("dir1", filepath.Join(osRoot(), "tmp"))
	directory2 := NewDirectory("dir1", filepath.Join(osRoot(), "tMp"))
	if directory1.Equals(directory2) {
		t.Fatal("directories were found to be equal but were not")
	}
}

func TestDirectory_Equals_FalseWhenDifferentName(t *testing.T) {
	directory1 := NewDirectory("dir1", filepath.Join(osRoot(), "tmp"))
	directory2 := NewDirectory("dir2", filepath.Join(osRoot(), "tmp"))
	if directory1.Equals(directory2) {
		t.Fatal("directories were found to be equal but were not")
	}
}

func TestDirectory_Equals_FalseWhenDifferentPath(t *testing.T) {
	directory1 := NewDirectory("dir1", filepath.Join(osRoot(), "tmp"))
	directory2 := NewDirectory("dir1", filepath.Join(osRoot(), "tmp", "dir"))
	if directory1.Equals(directory2) {
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
			if subDir, ok := dir.SubDirectories()["newdir"]; ok {
				if subDir.Name() != "newdir" {
					t.Fatalf("file name was not set correctly. expected %s but was %s", "newdir", subDir.Name())
				}
				if subDir.Path() != tt.dirFullPath {
					t.Fatalf("file path was not set correctly. expected %s but was %s", tt.dirFullPath, subDir.Path())
				}
			} else {
				t.Fatalf("subdirectory was not found in the subdirectories")
			}
		})
	}
}

func TestDirectory_AddDirectory_CreatesSubdirectoriesAsNecessary(t *testing.T) {
	dir := NewDirectory("dir", filepath.Join(osRoot(), "tmp"))
	_, err := dir.AddDirectory(filepath.Join(osRoot(), "tmp", "dir", "subdir1", "subdir2", "subdir3"))
	if err != nil {
		t.Fatal(err)
	}
	if subdir1, ok := dir.SubDirectories()["subdir1"]; !ok {
		t.Fatal("subdir1 was not created")
		expected := Directory{name: "subdir1", path: filepath.Join(osRoot(), "tmp", "dir")}
		if !subdir1.Equals(&expected) {
			t.Fatal("subdir1 structure was incorrect")
		}
	}
	if subdir2, ok := dir.SubDirectories()["subdir1"].SubDirectories()["subdir2"]; !ok {
		t.Fatal("subdir2 was not created")
		expected := Directory{name: "subdir2", path: filepath.Join(osRoot(), "tmp", "dir", "subdir1")}
		if !subdir2.Equals(&expected) {
			t.Fatal("subdir2 structure was incorrect")
		}
	}
	if subdir3, ok := dir.SubDirectories()["subdir1"].SubDirectories()["subdir2"].SubDirectories()["subdir3"]; !ok {
		t.Fatal("subdir1 was not created")
		expected := Directory{name: "subdir3", path: filepath.Join(osRoot(), "tmp", "dir", "subdir1", "subdir2")}
		if !subdir3.Equals(&expected) {
			t.Fatal("subdir1 structure was incorrect")
		}
	}
}

func TestDirectory_AddDirectory_AddingSiblingAsParentOfNewDirectoryDoesntDeleteCurrentDirectories(t *testing.T) {
	dir := NewDirectory("dir", filepath.Join(osRoot(), "tmp"))
	_, err := dir.AddDirectory(filepath.Join(osRoot(), "tmp", "dir", "subdir1", "subsub1"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = dir.AddDirectory(filepath.Join(osRoot(), "tmp", "dir", "subdir1", "subsub2", "subsubsub1"))
	if err != nil {
		t.Fatal(err)
	}
	subsub1, err := dir.GetDirectory(filepath.Join(osRoot(), "tmp", "dir", "subdir1", "subsub1"))
	if err != nil {
		t.Fatal(err)
	}
	if subsub1 == nil {
		t.Fatal("subsub1 was deleted but should not have been")
	}
}

func TestDirectory_AddDirectory_LeafDoesNotContainChildren(t *testing.T) {
	dir := NewDirectory("dir", filepath.Join(osRoot(), "tmp"))
	_, err := dir.AddDirectory(filepath.Join(osRoot(), "tmp", "dir", "subdir1", "subdir2", "subdir3"))
	if err != nil {
		t.Fatal(err)
	}
	if subdir3, ok := dir.SubDirectories()["subdir1"].SubDirectories()["subdir2"].SubDirectories()["subdir3"]; ok {
		if subdir3.SubDirectories() != nil {
			t.Fatal("subdir3 contained unexpected children")
		}
		if subdir3.Files() != nil {
			t.Fatal("subdir3 contained unexpected children")
		}
	} else {
		t.Fatal("subdir3 did not exist")
	}
}

func TestDirectory_AddDirectory_ReturnsErrorIfNotSubdirectory(t *testing.T) {
	dir := NewDirectory("dir", filepath.Join(osRoot(), "tmp"))
	_, err := dir.AddDirectory(filepath.Join(osRoot(), "tmp", "other", "subdir1", "subdir2", "subdir3"))
	if err == nil {
		t.Fatal("error should have been returned but was nil")
	}
}

func TestDirectory_AddDirectory_ReturnsErrorIfDifferentParent(t *testing.T) {
	dir := NewDirectory("dir", filepath.Join(osRoot(), "tmp"))
	_, err := dir.AddDirectory(filepath.Join(osRoot(), "other", "dir", "subdir1", "subdir2", "subdir3"))
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
			if found.Path() != expectedPath {
				t.Fatalf("found path did not match expected. expected path: "+
					"'%s' actual path: '%s'", expectedPath, found.Path())
			}
			if found.Name() != expectedName {
				t.Fatalf("found name did not match expected. expected name: "+
					"'%s' actual name: '%s'", expectedName, found.Name())
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
			if found.Path() != expectedPath {
				t.Fatalf("found path did not match expected. expected path: "+
					"'%s' actual path: '%s'", expectedPath, found.Path())
			}
			if found.Name() != expectedName {
				t.Fatalf("found name did not match expected. expected name: "+
					"'%s' actual name: '%s'", expectedName, found.Name())
			}
		})
	}
}

func TestDirectory_FindDirectoryBreadth(t *testing.T) {
	for _, tt := range FindTests {
		t.Run(tt.name, func(t *testing.T) {
			expectedPath, expectedName := filepath.Split(tt.fullDirPathToFind)
			expectedPath = filepath.Clean(expectedPath)
			found := tt.dir.FindDirectoryBreadth(expectedName)
			if found == nil {
				t.Fatal("nil was returned but actual directory was expected")
			}
			if found.Path() != expectedPath {
				t.Fatalf("found path did not match expected. expected path: "+
					"'%s' actual path: '%s'", expectedPath, found.Path())
			}
			if found.Name() != expectedName {
				t.Fatalf("found name did not match expected. expected name: "+
					"'%s' actual name: '%s'", expectedName, found.Name())
			}
		})
	}
}

func TestDirectory_MapFnDepth(t *testing.T) {
	for _, tt := range directoryMapTests() {
		fn := func(dir *Directory) error {
			dir.name = dir.name + "-new"
			for _, file := range dir.files {
				ext := filepath.Ext(file.name)
				file.name = strings.TrimSuffix(file.name, ext) + "-new" + ext
			}

			return nil
		}

		t.Run(tt.name, func(t *testing.T) {
			actual := &tt.dir
			err := tt.dir.MapFnDepth(fn)
			if err != nil {
				t.Fatal(err)
			}
			if expected := &tt.mappedDir; !actual.StructureEquals(expected)  {
				actualPrinted, _ := actual.Print()
				expectedPrinted, _ := expected.Print()
				t.Fatalf("structure of directories did not match\nexpected: %s\nactual: %s", expectedPrinted, actualPrinted)
			}
		})
	}
}

func TestDirectory_MapFnBreadth(t *testing.T) {
	for _, tt := range directoryMapTests() {
		fn := func(dir *Directory) error {
			dir.name = dir.name + "-new"
			for _, file := range dir.files {
				ext := filepath.Ext(file.name)
				file.name = strings.TrimSuffix(file.name, ext) + "-new" + ext
			}

			return nil
		}

		t.Run(tt.name, func(t *testing.T) {
			actual := &tt.dir
			err := tt.dir.MapFnBreadth(fn)
			if err != nil {
				t.Fatal(err)
			}
			if expected := &tt.mappedDir; !actual.StructureEquals(expected)  {
				actualPrinted, _ := actual.Print()
				expectedPrinted, _ := expected.Print()
				t.Fatalf("structure of directories did not match\nexpected: %s\nactual: %s", expectedPrinted, actualPrinted)
			}
		})
	}
}

func TestDirectory_Print(t *testing.T) {
	for _, tt := range DirectoryPrint {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := tt.dir.Print()
			if err != nil {
				t.Fatal(err)
			}
			expected := strings.ReplaceAll(tt.expected, "~", "/")
			expected = strings.ReplaceAll(expected, "//", "/")
			if actual != expected {
				t.Fatalf("printed output did not match.\nexpected: \n%s\nactual: \n%s", expected, actual)
			}
		})
	}
}
