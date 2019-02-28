package structure

import (
	"github.com/auroq/symfigurator/pkg/structure"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestGetDirectoryStructure(t *testing.T) {

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	expected := structure.Directory{
		Path: tmpDir,
		SubDirectories: map[string]structure.Directory{
			"FolderA": {
				"FolderA",
				path.Join(tmpDir, "FolderA"),
				nil,
				nil,
			},
			"FolderB": {
				"FolderB",
				path.Join(tmpDir, "FolderB"),
				nil,
				nil,
			},
		},
	}
	err = os.Mkdir(path.Join(tmpDir, "FolderA"), 0700)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir(path.Join(tmpDir, "FolderB"), 0700)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := structure.GetDirectoryStructure(tmpDir)
	if !actual.Equals(expected) {
		t.Fatal(err)
	}
}
