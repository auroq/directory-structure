package structure

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	name string
	path string
}

// Name returns the name of the File
func (file File) Name() string { return file.name }

// Path returns the path to the File excluding the File itself
func (file File) Path() string { return filepath.Clean(file.path) }

// FullPath returns the full path to the File including the File itself
func (file File) FullPath() string { return filepath.Clean(filepath.Join(file.path, file.name)) }

// Equals determines if other is equivalent to the current File.
func (file File) Equals(other *File) bool {
	return filepath.Clean(file.path) == filepath.Clean(other.path) &&
		file.name == other.name
}

// NewFile creates a new File using a name and a path
// Name is the name of of the File itself.
// Path is the path to the File not including name
func NewFile(name string, path string) File {
	return File{name: name, path: path}
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
	newFile := NewFile(name, path)
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
	if parent.files == nil {
		parent.files = map[string]*File{}
	}
	parent.files[name] = &newFile
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

	if file := fileDir.File(name); file != nil {
		return file, nil
	}
	return nil, errors.New(fmt.Sprintf("file could not be found in directory '%s'", dir.Path()))
}

// FindFileDepth searches the directory tree for a File using depth first search.
// When it finds a File with name fileName, it returns it.
// If the File is not found, nil is returned
func (dir Directory) FindFileDepth(fileName string) *File {
	if file := dir.File(fileName); file != nil {
		return file
	}
	for _, subDir := range dir.SubDirectories() {
		d := subDir.FindFileDepth(fileName)
		if d != nil {
			return d
		}
	}
	return nil
}

// FindFileBreadth searches the directory tree for a File using breadth first search.
// When it finds a File with name fileName, it returns it.
// If the File is not found, nil is returned
func (dir Directory) FindFileBreadth(fileName string) *File {
	queue := []*Directory{&dir}
	for len(queue) > 0 {
		pop := queue[0]
		queue = queue[1:]
		if file := pop.File(fileName); file != nil {
			return file
		}
		for _, subDir := range pop.subDirectories {
			queue = append(queue, subDir)
		}
	}
	return nil
}
