package structure

import (
	"github.com/auroq/directory-structure/pkg/structure"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGetDirectoryStructure(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	path, name := filepath.Split(tmpDir)
	expected := structure.Directory{ Path: path, Name: name, }
	_, err = expected.AddDirectory(filepath.Join(tmpDir, "dir1"))
	if err != nil {
		t.Fatal(err)
	}
	dir2, err := expected.AddDirectory(filepath.Join(tmpDir, "dir2"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = expected.AddDirectory(filepath.Join(tmpDir, "dir3"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = dir2.AddDirectory(filepath.Join(tmpDir, "dir2", "sub1"))
	if err != nil {
		t.Fatal(err)
	}
	sub2, err := dir2.AddDirectory(filepath.Join(tmpDir, "dir2", "sub2"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = dir2.AddDirectory(filepath.Join(tmpDir, "dir2", "sub3"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = sub2.AddFile(filepath.Join(tmpDir, "dir2", "sub2", "file"))
	if err != nil {
		t.Fatal(err)
	}

	err = os.Mkdir(filepath.Join(tmpDir, "dir1"), 0700)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir(filepath.Join(tmpDir, "dir2"), 0700)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir(filepath.Join(tmpDir, "dir3"), 0700)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir(filepath.Join(tmpDir, "dir2", "sub1"), 0700)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir(filepath.Join(tmpDir, "dir2", "sub2"), 0700)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir(filepath.Join(tmpDir, "dir2", "sub3"), 0700)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.OpenFile(filepath.Join(tmpDir, "dir2", "sub2", "file"), os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := structure.GetDirectoryStructure(tmpDir)
	if !actual.Equals(&expected) {
		t.Fatal(err)
	}
}
