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
		return nil, os.ErrNotExist
	}
	if !d.IsDir() {
		return nil, errors.New(fmt.Sprintf("fullPath '%s' is not a directory", fullPath))
	}
	rootPath, rootName := filepath.Split(fullPath)
	rootPath = filepath.Clean(rootPath)
	root := NewDirectory(rootName, rootPath)
	err = filepath.Walk(fullPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path == root.FullPath() {
				return nil
			}
			if info.IsDir() {
				_, err = root.AddDirectory(path)
				if err != nil {
					return err
				}
			} else {
				_, err = root.AddFile(path)
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
	var directory *Directory
	name := pathSlice[0]
	if existingDir := dir.SubDirectory(name); existingDir != nil {
		directory = existingDir
	} else {
		path := filepath.Join(dir.Path(), dir.Name())
		newDirectory := NewDirectory(name, path)
		directory = &newDirectory
		dir.subDirectories[name] = directory
	}
	return directory.createPath(pathSlice[1:])
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
