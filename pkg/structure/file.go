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

// Equals determines if other is equivalent to the current File.
func (file File) Equals(other *File) bool {
	return filepath.Clean(file.Path) == filepath.Clean(other.Path) &&
		file.Name == other.Name
}

// AddFile creates a new File and adds it to the current Directory tree
// The new File will contain a name and a path specified by fullPath.
// AddDirectory will return the new File and an error if fullPath is not a
// descendant of the current Directory
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

// GetFile transverses the current Directory to find a File whose
// path is fullPath. It returns the File and an error if fullPath is
// not a descendant of the current Directory.
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
// When it finds a File with name fileName, it returns it.
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
