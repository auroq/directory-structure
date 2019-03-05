package structure

import (
	"path/filepath"
	"testing"
)

func TestDescendants_ContainsDirectory_TrueWhenContainsDirectory(t *testing.T) {
	dir1 := NewDirectory("dir1", osRoot())
	dir2 := NewDirectory("dir2", osRoot())
	subdir1 := NewDirectory("subdir1", filepath.Join(osRoot(), "dir1"))
	file1 := NewFile("file1", osRoot())
	file2 := NewFile("file2", osRoot())
	file3 := NewFile("file1", filepath.Join(osRoot(), "dir1"))
	desc := Descendants{
		[]*Directory{&dir1, &dir2, &subdir1},
		[]*File{&file1, &file2, &file3},
	}
	if !desc.ContainsDirectory(&dir2) {
		t.Fatal("directory was not found in descendants")
	}
}

func TestDescendants_ContainsDirectory_FalseWhenNotContainsDirectory(t *testing.T) {
	dir1 := NewDirectory("dir1", osRoot())
	dir2 := NewDirectory("dir2", osRoot())
	subdir1 := NewDirectory("subdir1", filepath.Join(osRoot(), "dir1"))
	file1 := NewFile("file1", osRoot())
	file2 := NewFile("file2", osRoot())
	file3 := NewFile("file1", filepath.Join(osRoot(), "dir1"))
	desc := Descendants{
		[]*Directory{&dir1, &subdir1},
		[]*File{&file1, &file2, &file3},
	}
	if desc.ContainsDirectory(&dir2) {
		t.Fatal("directory was found in descendants but should not have been")
	}
}

func TestDescendants_ContainsFile_TrueWhenContainsFile(t *testing.T) {
	dir1 := NewDirectory("dir1", osRoot())
	dir2 := NewDirectory("dir2", osRoot())
	subdir1 := NewDirectory("subdir1", filepath.Join(osRoot(), "dir1"))
	file1 := NewFile("file1", osRoot())
	file2 := NewFile("file2", osRoot())
	file3 := NewFile("file1", filepath.Join(osRoot(), "dir1"))
	desc := Descendants{
		[]*Directory{&dir1, &dir2, &subdir1},
		[]*File{&file1, &file2, &file3},
	}
	if !desc.ContainsFile(&file2) {
		t.Fatal("file was not found in descendants")
	}
}

func TestDescendants_ContainsFile_FalseWhenNotContainsFile(t *testing.T) {
	dir1 := NewDirectory("dir1", osRoot())
	dir2 := NewDirectory("dir2", osRoot())
	subdir1 := NewDirectory("subdir1", filepath.Join(osRoot(), "dir1"))
	file1 := NewFile("file1", osRoot())
	file2 := NewFile("file2", osRoot())
	file3 := NewFile("file1", filepath.Join(osRoot(), "dir1"))
	desc := Descendants{
		[]*Directory{&dir1, &dir2, &subdir1},
		[]*File{&file1, &file3},
	}
	if desc.ContainsFile(&file2) {
		t.Fatal("file was found in descendants but should not have been")
	}
}

func TestDirectory_GetAllDescendants(t *testing.T) {
	dir := NewDirectory("dir", filepath.Join(osRoot(), "tmp"))
	subdir1, err := dir.AddDirectory(filepath.Join(osRoot(), "tmp", "dir", "subdir1"))
	if err != nil {
		t.Fatal(err)
	}
	subdir2, err := dir.AddDirectory(filepath.Join(osRoot(), "tmp", "dir", "subdir1", "subdir2"))
	if err != nil {
		t.Fatal(err)
	}
	subdir3, err := dir.AddDirectory(filepath.Join(osRoot(), "tmp", "dir", "subdir1", "subdir2", "subdir3"))
	if err != nil {
		t.Fatal(err)
	}
	file1, err := dir.AddFile(filepath.Join(osRoot(), "tmp", "dir", "file1.txt"))
	if err != nil {
		t.Fatal(err)
	}
	file2, err := dir.AddFile(filepath.Join(osRoot(), "tmp", "dir", "subdir1", "file2.txt"))
	if err != nil {
		t.Fatal(err)
	}
	file3, err := dir.AddFile(filepath.Join(osRoot(), "tmp", "dir", "subdir1", "subdir2", "file3.txt"))
	if err != nil {
		t.Fatal(err)
	}
	file4, err := dir.AddFile(filepath.Join(osRoot(), "tmp", "dir", "subdir1", "subdir2", "subdir3", "file4.txt"))
	if err != nil {
		t.Fatal(err)
	}

	dirList := []*Directory{subdir1, subdir2, subdir3}
	fileList := []*File{file1, file2, file3, file4}

	descendants := dir.GetAllDescendants()

	if descDirLen := len(descendants.Directories); descDirLen != len(dirList) {
		t.Fatalf("incorrect number of descendant directories found expected: %d actual: %d", len(dirList), descDirLen)
	}

	if descFileLen := len(descendants.Files); descFileLen != len(fileList) {
		t.Fatalf("incorrect number of descendant files found expected: %d actual: %d", len(fileList), descFileLen)
	}

	for _, directory := range dirList {
		if !(descendants.ContainsDirectory(directory)) {
			t.Fatalf("directory '%s' was not found in descendants", directory.FullPath())
		}
	}

	for _, file := range fileList {
		if !(descendants.ContainsFile(file)) {
			t.Fatalf("directory '%s' was not found in descendants", file.FullPath())
		}
	}
}
