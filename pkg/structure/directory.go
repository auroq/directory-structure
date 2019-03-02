package structure

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Directory struct {
	Name           string
	Path           string
	SubDirectories map[string]*Directory
	Files          map[string]*File
}

// Equals determines if other is equivalent to the current Directory.
// It does so using only path and name and therefore does not take
// into account the structure of either Directory's children.
func (dir Directory) Equals(other *Directory) bool {
	return dir.Path == other.Path && dir.Name == other.Name
}

// AddDirectory creates a new Directory and adds it to the current Directory tree
// The new Directory will contain a name and a path specified by fullPath.
// SubDirectories and Files of the new Directory will be nil
// AddDirectory will return the new Directory and an error if fullPath is not a
// descendant of the current Directory
func (dir *Directory) AddDirectory(fullPath string) (*Directory, error) {
	path, name := filepath.Split(fullPath)
	path = filepath.Clean(path)
	if !dir.IsSubPath(fullPath) {
		return nil, errors.New(fmt.Sprintf("fullPath must be a subdirectory of the directory to which it is " +
			"being added: '%s' is not a subdirectory of '%s'", fullPath, filepath.Join(dir.Path, dir.Name)))
	}

	var parent *Directory
	newDirectory := Directory{Name: name, Path: path}
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
	if parent.SubDirectories == nil {
		parent.SubDirectories = map[string]*Directory{}
	}
	parent.SubDirectories[name] = &newDirectory
	return &newDirectory, nil
}

// GetDirectory transverses the current Directory to find a directory whose
// path is fullPath. It returns the Directory and an error if fullPath is
// not a descendant of the current Directory.
func (dir Directory) GetDirectory(fullPath string) (*Directory, error) {
	path, name := filepath.Split(fullPath)
	path = filepath.Clean(path)
	if path == dir.Path && name == dir.Name {
		return &dir, nil
	}
	currentDir := filepath.Join(dir.Path, dir.Name)
	if len(path) < len(currentDir) || path[:len(currentDir)] != currentDir {
		return nil, errors.New(fmt.Sprintf("item '%s' is not found in directory '%s'", fullPath, dir.Path))
	}
	pathSlice := strings.Split(strings.TrimPrefix(fullPath, currentDir), string(os.PathSeparator))
	if len(pathSlice) >= 1 && pathSlice[0] == "" {
		pathSlice = pathSlice[1:]
	}
	return dir.findPath(pathSlice)
}

// FindDirectoryDepth searches the directory tree for a Directory using depth first search.
// When it finds a Directory with name dirName, it returns it.
// If the Directory is not found, nil is returned
func (dir Directory) FindDirectoryDepth(dirName string) *Directory {
	if dir.Name == dirName {
		return &dir
	}
	for _, subDir := range dir.SubDirectories {
		d := subDir.FindDirectoryDepth(dirName)
		if d != nil {
			return d
		}
	}
	return nil
}