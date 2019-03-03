package structure

import (
	"path/filepath"
	"testing"
)

func TestFile_Path_CleansName(t *testing.T) {
	for _, tt := range cleanTests {
		t.Run(tt.name, func(t *testing.T) {
			file := NewFile("file1", tt.path)
			if file.Path() != tt.cleanPath {
				t.Fatalf("path was not properly cleaned expected: '%s' actual: '%s'", tt.cleanPath, file.Path())
			}
		})
	}
}

func TestFile_FullPath(t *testing.T) {
	for _, tt := range fullPathTests {
		t.Run(tt.testName, func(t *testing.T) {
			file := NewFile(tt.name, tt.path)
			if file.FullPath() != tt.fullPath {
				t.Fatalf("full path did not match expected: '%s' actual: '%s'", tt.fullPath, file.FullPath())
			}
		})
	}
}

func TestFile_File_ReturnsFile(t *testing.T) {
	dir := NewDirectory("dir1", "/tmp")
	sub1, err := dir.AddFile("/tmp/dir1/sub1")
	if err != nil {
		t.Fatal(err)
	}
	if subFile := dir.File("sub1"); !subFile.Equals(sub1) {
		t.Fatal("file returned did not match")
	}
}

func TestFile_File_NilWhenFilesNil(t *testing.T) {
	dir := NewDirectory("dir1", "/tmp")
	if subFile := dir.File("nonexistent"); subFile != nil {
		t.Fatal("file did not exist and should have been nil")
	}
}

func TestFile_File_NilWhenNotFound(t *testing.T) {
	dir := NewDirectory("dir1", "/tmp")
	_, err := dir.AddFile("/tmp/dir1/sub1")
	if err != nil {
		t.Fatal(err)
	}
	if subFile := dir.File("nonexistent"); subFile != nil {
		t.Fatal("file did not exist and should have been nil")
	}
}

func TestFile_Equals_TrueWithIdentity(t *testing.T) {
	file := NewFile("file1", "/tmp")
	if !file.Equals(&file) {
		t.Fatal("files were equal but were not found to be")
	}
}

func TestFile_Equals_TrueWithDifferentInstances(t *testing.T) {
	file1 := NewFile("file1", "/tmp")
	file2 := NewFile("file1", "/tmp")
	if !file1.Equals(&file2) {
		t.Fatal("files were equal but were not found to be")
	}
}

func TestFile_Equals_TrueWhenPathStringNotClean(t *testing.T) {
	file1 := NewFile("file1", "/tmp/")
	file2 := NewFile("file1", "/tmp")
	if !file1.Equals(&file2) {
		t.Fatal("files were equal but were not found to be")
	}
}

func TestFile_Equals_FalseWhenPathStringNotSameCase(t *testing.T) {
	file1 := NewFile("file1", "/tmp")
	file2 := NewFile("file1", "/tMp")
	if file1.Equals(&file2) {
		t.Fatal("files were found to be equal but were not")
	}
}

func TestFile_Equals_FalseWhenDifferentNname(t *testing.T) {
	file1 := NewFile("file1", "/tmp")
	file2 := NewFile("file2", "/tMp")
	if file1.Equals(&file2) {
		t.Fatal("files were found to be equal but were not")
	}
}

func TestFile_Equals_FalseWhenDifferentPath(t *testing.T) {
	file1 := NewFile("file1", "/tmp")
	file2 := NewFile("file2", "/tmp/dir")
	if file1.Equals(&file2) {
		t.Fatal("files were found to be equal but were not")
	}
}

func TestDirectory_AddFile(t *testing.T) {
	for _, tt := range DirectoryIdentities {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := tt.dir.GetDirectory(tt.dirFullPath)
			if err != nil {
				t.Fatal(err)
			}
			_, err = dir.AddFile(filepath.Join(tt.dirFullPath, "newfile"))
			if err != nil {
				t.Fatal(err)
			}
			if file, ok := tt.dir.Files()["file1.txt"]; ok {
				if file.Name() != "file1.txt" {
					t.Fatalf("file name was not set correctly. expected %s but was %s", "file1.txt", file.Name())
				}
				if file.Path() != tt.dirFullPath {
					t.Fatalf("file path was not set correctly. expected %s but was %s", tt.dirFullPath, file.path)
				}
			}
		})
	}
}

func TestDirectory_AddFile_CreatesSubdirectoriesAsNecessary(t *testing.T) {
	dir := NewDirectory("dir", "/tmp")
	_, err := dir.AddFile("/tmp/dir/subdir1/subdir2/file")
	if err != nil {
		t.Fatal(err)
	}
	if subdir1, ok := dir.SubDirectories()["subdir1"]; !ok {
		t.Fatal("subdir1 was not created")
		expected := NewDirectory("subdir1", "/tmp/dir")
		if !subdir1.Equals(&expected) {
			t.Fatal("subdir1 structure was incorrect")
		}
	}
	if subdir2, ok := dir.SubDirectories()["subdir1"].SubDirectories()["subdir2"]; !ok {
		t.Fatal("subdir2 was not created")
		expected := NewDirectory("subdir2", "/tmp/dir/subdir1")
		if !subdir2.Equals(&expected) {
			t.Fatal("subdir2 structure was incorrect")
		}
	}
	if file, ok := dir.SubDirectories()["subdir1"].SubDirectories()["subdir2"].Files()["file"]; !ok {
		t.Fatal("file was not created")
		expected := NewFile("file", "/tmp/dir/subdir1/subdir2")
		if !file.Equals(&expected) {
			t.Fatal("file structure was incorrect")
		}
	}
}

func TestDirectory_AddFile_ReturnsErrorIfNotSubdirectory(t *testing.T) {
	dir := NewDirectory("dir", "/tmp")
	_, err := dir.AddDirectory("/tmp/other/subdir1/subdir2/subdir3/file.txt")
	if err == nil {
		t.Fatal("error should have been returned but was nil")
	}
}

func TestDirectory_AddFile_ReturnsErrorIfDifferentParent(t *testing.T) {
	dir := NewDirectory("dir", "/tmp")
	_, err := dir.AddDirectory("/other/dir/subdir1/subdir2/subdir3/file.txt")
	if err == nil {
		t.Fatal("error should have been returned but was nil")
	}
}

func TestDirectory_FindFile(t *testing.T) {
	for _, tt := range FindTests {
		t.Run(tt.name, func(t *testing.T) {
			expectedPath, expectedName := filepath.Split(tt.fullFilePathToFind)
			expectedPath = filepath.Clean(expectedPath)
			found, err := tt.dir.GetFile(tt.fullFilePathToFind)
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

func TestDirectory_FindFileDepth(t *testing.T) {
	for _, tt := range FindTests {
		t.Run(tt.name, func(t *testing.T) {
			expectedPath, expectedName := filepath.Split(tt.fullFilePathToFind)
			expectedPath = filepath.Clean(expectedPath)
			found := tt.dir.FindFileDepth(expectedName)
			if found == nil {
				t.Fatal("nil was returned but actual file was expected")
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
