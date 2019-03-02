package structure

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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
	root := Directory{Name: rootName, Path: rootPath}
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
	if !(dir.Path == other.Path &&
		dir.Name == other.Name &&
		len(dir.SubDirectories) == len(other.SubDirectories) &&
		len(dir.Files) == len(other.Files)) {
		return false
	}
	for fileName, file := range dir.Files {
		if otherFile, ok := other.Files[fileName]; !ok || !otherFile.Equals(file) {
			return false
		}
	}

	for subDirectoryName, subDirectory := range dir.SubDirectories {
		if otherSubDir, ok := other.SubDirectories[subDirectoryName]; !ok || !otherSubDir.StructureEquals(subDirectory) {
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
	currentDir := filepath.Join(dir.Path, dir.Name)
	relPath := strings.TrimPrefix(fullPath, currentDir)
	return relPath != fullPath
}

func (dir *Directory) relativePath(fullPath string) string {
	currentDir := filepath.Join(dir.Path, dir.Name)
	path := strings.TrimPrefix(fullPath, currentDir)
	path = strings.TrimPrefix(path, "/")
	return path
}

func (dir *Directory) createPath(pathSlice []string) (*Directory, error) {
	if len(pathSlice) <= 0 {
		return dir, nil
	}
	if dir.SubDirectories == nil {
		dir.SubDirectories = map[string]*Directory{}
	}
	name := pathSlice[0]
	path := filepath.Join(dir.Path, dir.Name)
	newDirectory := Directory{Name: pathSlice[0], Path: path}
	dir.SubDirectories[name] = &newDirectory
	return newDirectory.createPath(pathSlice[1:])
}

func (dir Directory) findPath(relativePath []string) (*Directory, error) {
	if subDir, ok := dir.SubDirectories[relativePath[0]]; ok {
		if len(relativePath) == 1 {
			return subDir, nil
		}
		return subDir.findPath(relativePath[1:])
	}
	return nil, errors.New(fmt.Sprintf("directory could not be found. "+
		"Current dir: %s Looking for: %s", dir.Path, strings.Join(relativePath, string(os.PathSeparator))))
}
