package structure

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Directory struct {
	name           string
	path           string
	subDirectories map[string]*Directory
	files          map[string]*File
}

// Name returns the name of the Directory
func (dir Directory) Name() string { return dir.name }

// Path returns the path to the Directory excluding the Directory itself
func (dir Directory) Path() string { return filepath.Clean(dir.path) }

// FullPath returns the full path to the Directory including the Directory itself
func (dir Directory) FullPath() string { return filepath.Clean(filepath.Join(dir.path, dir.name)) }

// SubDirectories returns a map where the key is the name of each subdirectory and the value is a
// pointer to the subdirectory
func (dir Directory) SubDirectories() map[string]*Directory { return dir.subDirectories }

// Files returns a map where the key is the name of each File and the value is a
// pointer to the File
func (dir Directory) Files() map[string]*File { return dir.files }

// SubDirectory returns a s pointer to a subdirectory named name
// If returns nil if the given name is not found
func (dir Directory) SubDirectory(name string) *Directory {
	return dir.subDirectories[name]
}

// SubDirectory returns a s pointer to a File named name
// If returns nil if the given name is not found
func (dir Directory) File(name string) *File {
	return dir.files[name]
}

// Equals determines if other is equivalent to the current Directory.
// It does so using only path and name and therefore does not take
// into account the structure of either Directory's children.
func (dir Directory) Equals(other *Directory) bool {
	return filepath.Clean(dir.path) == filepath.Clean(other.path) &&
		dir.name == other.name
}

// NewDirectory creates a new Directory using a name and a path.
// Name is the name of of the Directory itself.
// Path is the path to the Directory not including name
func NewDirectory(name string, path string) Directory {
	return Directory{name: name, path: path}
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
		return nil, errors.New(fmt.Sprintf("fullPath must be a subdirectory of the directory to which it is "+
			"being added: '%s' is not a subdirectory of '%s'", fullPath, filepath.Join(dir.Path(), dir.Name())))
	}

	var parent *Directory
	newDirectory := NewDirectory(name, path)
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
	if parent.subDirectories == nil {
		parent.subDirectories = map[string]*Directory{}
	}
	parent.subDirectories[name] = &newDirectory
	return &newDirectory, nil
}

// GetDirectory transverses the current Directory to find a directory whose
// path is fullPath. It returns the Directory and an error if fullPath is
// not a descendant of the current Directory.
func (dir Directory) GetDirectory(fullPath string) (*Directory, error) {
	path, name := filepath.Split(fullPath)
	path = filepath.Clean(path)
	if path == dir.Path() && name == dir.Name() {
		return &dir, nil
	}
	currentDir := filepath.Join(dir.Path(), dir.Name())
	if len(path) < len(currentDir) || path[:len(currentDir)] != currentDir {
		return nil, errors.New(fmt.Sprintf("item '%s' is not found in directory '%s'", fullPath, dir.Path()))
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
	if dir.Name() == dirName {
		return &dir
	}
	for _, subDir := range dir.SubDirectories() {
		d := subDir.FindDirectoryDepth(dirName)
		if d != nil {
			return d
		}
	}
	return nil
}

// FindDirectoryBreadth searches the directory tree for a Directory using breadth first search.
// When it finds a Directory with name dirName, it returns it.
// If the Directory is not found, nil is returned
func (dir Directory) FindDirectoryBreadth(dirName string) *Directory {
	queue := []*Directory{&dir}
	for len(queue) > 0 {
		pop := queue[0]
		queue = queue[1:]
		if pop.name == dirName {
			return pop
		}
		for _, subDir := range pop.subDirectories {
			queue = append(queue, subDir)
		}
	}
	return nil
}