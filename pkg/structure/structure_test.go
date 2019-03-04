package structure

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestDirectory_StructureEquals_WithIdentity(t *testing.T) {
	for _, tt := range DirectoryIdentities {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.dir.StructureEquals(&tt.dir) {
				t.Fatal("directory structures were found to be unequal but were equal")
			}
		})
	}
}

func TestDirectory_StructureEquals_WithIdentityAndFile(t *testing.T) {
	for _, tt := range FindTests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.dir.StructureEquals(&tt.dir) {
				t.Fatal("directory structures were found to be unequal but were equal")
			}
		})
	}
}

func TestDirectory_StructureEquals_WhenNotEqual(t *testing.T) {
	for _, tt := range DirectoryIdentities {
		for _, ott := range DirectoryIdentities {
			if tt.name == ott.name {
				continue
			}
			t.Run(fmt.Sprintf("%s_And_%s", tt.name, ott.name), func(t *testing.T) {
				if tt.dir.StructureEquals(&ott.dir) {
					t.Fatal("directory structures were found to be equal but were not")
				}
			})
		}
	}
}

func TestDirectory_StructureEquals_WhenNotEqualWithFile(t *testing.T) {
	for _, tt := range DirectoryIdentities {
		for _, ott := range DirectoryIdentities {
			t.Run(fmt.Sprintf("%s_And_%s", tt.name, ott.name), func(t *testing.T) {
				if tt.name == ott.name {
					_, err := tt.dir.AddFile(filepath.Join(tt.dir.Path(), tt.dir.Name(), "randomfile.txt"))
					if err != nil {
						t.Fatal(err)
					}
				}
				if tt.dir.StructureEquals(&ott.dir) {
					t.Fatal("directory structures were found to be equal but were not")
				}
			})
		}
	}
}

func TestDirectory_IsSubPath(t *testing.T) {
	for _, tt := range DirectoryIdentities {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.dir.IsSubPath(tt.dirFullPath) {
				t.Fatalf("'%s' is a subdirectory but was not found to be", tt.dirFullPath)
			}
		})
	}
}

func TestDirectory_IsSubPath_WhenPathIsParent(t *testing.T) {
	dir := Directory{name: "dir1", path: filepath.Join(osRoot(), "tmp")}
	if dir.IsSubPath(filepath.Join(osRoot(), "tmp")) {
		t.Fatalf("'/tmp' is not a subdirectory of '/tmp/dir1' but was found to be")
	}
}

func TestDirectory_IsSubPath_WhenPathIsSibling(t *testing.T) {
	parent := Directory{name: "dir1", path: filepath.Join(osRoot(), "tmp")}
	dir, err := parent.AddDirectory(filepath.Join(osRoot(), "tmp", "dir1", "subdir1"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = parent.AddDirectory(filepath.Join(osRoot(), "tmp", "dir1", "subdir2"))
	if err != nil {
		t.Fatal(err)
	}
	if dir.IsSubPath(filepath.Join(osRoot(), "tmp", "dir1", "subdir2")) {
		t.Fatalf("'/tmp/dir1/subdir2' is not a subdirectory of '/tmp/dir1/subdir1' but was found to be")
	}
}

func TestDirectory_IsSubPath_WhenPathIsUnrelated(t *testing.T) {
	dir := Directory{name: "dir1", path: filepath.Join(osRoot(), "tmp", "dir1")}
	if dir.IsSubPath(filepath.Join(osRoot(), "other")) {
		t.Fatalf("'/tmp' is not a subdirectory of '/tmp/dir1' but was found to be")
	}
}

func TestDescendants_Contains_TrueWhenContainsDirectory(t *testing.T) {
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
	if !desc.Contains(dir2) {
		t.Fatal("directory was not found in descendants")
	}
}

func TestDescendants_Contains_FalseWhenNotContainsDirectory(t *testing.T) {
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
	if desc.Contains(dir2) {
		t.Fatal("directory was found in descendants but should not have been")
	}
}

func TestDescendants_Contains_TrueWhenContainsFile(t *testing.T) {
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
	if !desc.Contains(file2) {
		t.Fatal("directory was not found in descendants")
	}
}

func TestDescendants_Contains_FalseWhenNotContainsFile(t *testing.T) {
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
	if desc.Contains(file2) {
		t.Fatal("directory was found in descendants but should not have been")
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

	if descDirLen := len(descendants.directories); descDirLen != len(dirList) {
		t.Fatalf("incorrect number of descendant directories found expected: %d actual: %d", len(dirList), descDirLen)
	}

	if descFileLen := len(descendants.files); descFileLen != len(fileList) {
		t.Fatalf("incorrect number of descendant files found expected: %d actual: %d", len(fileList), descFileLen)
	}

	for _, directory := range dirList {
		if !(descendants.Contains(directory)) {
			t.Fatalf("directory '%s' was not found in descendants", directory.FullPath())
		}
	}

	for _, file := range fileList {
		if !(descendants.Contains(file)) {
			t.Fatalf("directory '%s' was not found in descendants", file.FullPath())
		}
	}
}
