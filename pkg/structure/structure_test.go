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
	dir := Directory{name: "dir1", path: "/tmp"}
	if dir.IsSubPath("/tmp") {
		t.Fatalf("'/tmp' is not a subdirectory of '/tmp/dir1' but was found to be")
	}
}

func TestDirectory_IsSubPath_WhenPathIsSibling(t *testing.T) {
	parent := Directory{name: "dir1", path: "/tmp"}
	dir, err := parent.AddDirectory("/tmp/dir1/subdir1")
	if err != nil {
		t.Fatal(err)
	}
	_, err = parent.AddDirectory("/tmp/dir1/subdir2")
	if err != nil {
		t.Fatal(err)
	}
	if dir.IsSubPath("/tmp/dir1/subdir2") {
		t.Fatalf("'/tmp/dir1/subdir2' is not a subdirectory of '/tmp/dir1/subdir1' but was found to be")
	}
}

func TestDirectory_IsSubPath_WhenPathIsUnrelated(t *testing.T) {
	dir := Directory{name: "dir1", path: "/tmp/dir1"}
	if dir.IsSubPath("/other") {
		t.Fatalf("'/tmp' is not a subdirectory of '/tmp/dir1' but was found to be")
	}
}
