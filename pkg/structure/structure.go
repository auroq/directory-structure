package structure

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Descendants struct {
	directories []*Directory
	files       []*File
}

func (desc Descendants) Contains(item interface{}) bool {
	switch t := item.(type) {
	case *File:
		for _, file := range desc.files {
			if file.Equals(t) {
				return true
			}
		}
	case File:
		for _, file := range desc.files {
			if file.Equals(&t) {
				return true
			}
		}
	case *Directory:
		for _, dir := range desc.directories {
			if dir.Equals(t) {
				return true
			}
		}
	case Directory:
		for _, dir := range desc.directories {
			if dir.Equals(&t) {
				return true
			}
		}
	}
	return false
}

func (dir Directory) GetAllDescendants() Descendants {
	descDirs, descFiles := dir.getDescendants()
	return Descendants{directories: descDirs, files: descFiles}
}

func (dir Directory) getDescendants() (dirs []*Directory, files []*File) {
	for _, file := range dir.files {
		files = append(files, file)
	}
	for _, subdir := range dir.subDirectories {
		dirs = append(dirs, subdir)
		descDirs, descFiles := subdir.getDescendants()
		dirs = append(dirs, descDirs...)
		files = append(files, descFiles...)
	}
	return
}

// GetDirectoryStructure walks through a directory on disk and its descendants
// and builds a Directory tree containing that matches the filesystem on disk
// It returns the root Directory whose path is fullPath and an error if one occurs
func GetDirectoryStructure(fullPath string) (*Directory, error) {
	d, err := os.Stat(fullPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("fullPath '%s' does not exist", fullPath))
	}
	if !d.IsDir() {
		return nil, errors.New(fmt.Sprintf("fullPath '%s' is not a directory", fullPath))
	}
	rootPath, rootName := filepath.Split(fullPath)
	rootPath = filepath.Clean(rootPath)
	root := Directory{name: rootName, path: rootPath}
	err = filepath.Walk(fullPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			newFullPath := filepath.Join(path, info.Name())
			if info.IsDir() {
				_, err = root.AddDirectory(newFullPath)
				if err != nil {
					return err
				}
			} else {
				_, err = root.AddFile(newFullPath)
				if err != nil {
					return err
				}
			}
			return nil
		})
	return &root, err
}

// StructureEquals determines if other and its descendants are identical
// to the current Directory and its descendants. It takes into account
// the full structure of both Directories including SubDirectories
// and Files.
func (dir Directory) StructureEquals(other *Directory) bool {
	if !(dir.path == other.path &&
		dir.name == other.name &&
		len(dir.subDirectories) == len(other.subDirectories) &&
		len(dir.files) == len(other.files)) {
		return false
	}
	for fileName, file := range dir.files {
		if otherFile, ok := other.files[fileName]; !ok || !otherFile.Equals(file) {
			return false
		}
	}

	for subDirectoryName, subDirectory := range dir.subDirectories {
		if otherSubDir, ok := other.subDirectories[subDirectoryName]; !ok || !otherSubDir.StructureEquals(subDirectory) {
			return false
		}
	}
	return true
}

// IsSubPath determines if fullPath is a descendant of the current Directory.
// fullPath does not need to actually exist in the Directory. It just has
// to be a descendant. It returns true or false accordingly.
func (dir *Directory) IsSubPath(fullPath string) bool {
	fullPath = filepath.Clean(fullPath)
	relPath := strings.TrimPrefix(fullPath, dir.FullPath())
	return relPath != fullPath
}

func (dir *Directory) relativePath(fullPath string) string {
	path := strings.TrimPrefix(fullPath, dir.FullPath())
	path = strings.TrimPrefix(path, string(os.PathSeparator))
	return path
}

func (dir *Directory) createPath(pathSlice []string) (*Directory, error) {
	if len(pathSlice) <= 0 {
		return dir, nil
	}
	if dir.subDirectories == nil {
		dir.subDirectories = map[string]*Directory{}
	}
	name := pathSlice[0]
	path := filepath.Join(dir.Path(), dir.Name())
	newDirectory := Directory{name: pathSlice[0], path: path}
	dir.subDirectories[name] = &newDirectory
	return newDirectory.createPath(pathSlice[1:])
}

func (dir Directory) findPath(relativePath []string) (*Directory, error) {
	if subDir := dir.SubDirectory(relativePath[0]); subDir != nil {
		if len(relativePath) == 1 {
			return subDir, nil
		}
		return subDir.findPath(relativePath[1:])
	}
	return nil, errors.New(fmt.Sprintf("directory could not be found. "+
		"Current dir: %s Looking for: %s", dir.Path(), strings.Join(relativePath, string(os.PathSeparator))))
}
