package structure

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Name string
	Path string
}

func (file File) Equals(other *File) bool {
	return file.Path == other.Path && file.Name == other.Name
}

func (dir *Directory) AddFile(fullPath string) (*File, error) {
	path, name := filepath.Split(fullPath)
	path = filepath.Clean(path)
	if !dir.IsSubPath(fullPath) {
		return nil, errors.New("fullPath must be an immediate child of the directory to which it is being added")
	}

	var parent *Directory
	newFile := File{name, path}
	if relativePath := dir.relativePath(path); relativePath == "" {
		parent = dir
	} else {
		pathSlice := strings.Split(relativePath, string(os.PathSeparator))
		newParent, err := dir.createPath(pathSlice)
		if err != nil {
			return nil, err
		}
		parent = newParent
	}
	if parent.Files == nil {
		parent.Files = map[string]*File{}
	}
	parent.Files[name] = &newFile
	return &newFile, nil
}

func (dir Directory) GetFile(fullPath string) (*File, error) {
	path, name := filepath.Split(fullPath)
	path = filepath.Clean(path)
	fileDir, err := dir.GetDirectory(path)
	if err != nil {
		return nil, err
	}

	if file, ok := fileDir.Files[name]; ok {
		return file, nil
	}
	return nil, errors.New(fmt.Sprintf("file could not be found in directory '%s'", dir.Path))
}

// FindFileDepth searches the directory tree for a File using depth first search.
// When it finds a File with name filenName, it returns it.
// If the File is not found, nil is returned
func (dir Directory) FindFileDepth(fileName string) *File {
	if file, ok := dir.Files[fileName]; ok {
		return file
	}
	for _, subDir := range dir.SubDirectories {
		d := subDir.FindFileDepth(fileName)
		if d != nil {
			return d
		}
	}
	return nil
}
