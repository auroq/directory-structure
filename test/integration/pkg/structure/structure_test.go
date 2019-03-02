package structure

import (
	"fmt"
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

	_, err = os.OpenFile(filepath.Join(tmpDir, "dir2", "sub2", "file"), os.O_RDONLY|os.O_CREATE, 0700)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := structure.GetDirectoryStructure(tmpDir)
	if !actual.Equals(&expected) {
		t.Fatal(err)
	}
}

func TestGetDirectoryStructure_WhenFullPathIsNotADirectory(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.OpenFile(filepath.Join(tmpDir, "file"), os.O_RDONLY|os.O_CREATE, 0700)
	if err != nil {
		t.Fatal(err)
	}

	_, err = structure.GetDirectoryStructure(file.Name())
	if err == nil {
		t.Fatal("an error was expected but err was nil")
	}
	if expected := fmt.Sprintf("fullPath '%s' is not a directory", file.Name()); err.Error() != expected {
		t.Fatalf("error message was incorrect. expected: '%s' actual: '%s'", expected, err.Error())
	}
}

func TestGetDirectoryStructure_WhenFullPathDoesNotExist(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(tmpDir, "missingDir")
	_, err = structure.GetDirectoryStructure(path)
	if err == nil {
		t.Fatal("an error was expected but err was nil")
	}
	if expected := fmt.Sprintf("fullPath '%s' does not exist", path); err.Error() != expected {
		t.Fatalf("error message was incorrect. expected: '%s' actual: '%s'", expected, err.Error())
	}
}
