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

func (dir Directory) Equals(other *Directory) bool {
	return dir.Path == other.Path && dir.Name == other.Name
}

func (dir Directory) StructureEquals(other *Directory) bool {
	if !(dir.Path == other.Path &&
		dir.Name == other.Name &&
		len(dir.SubDirectories) == len(other.SubDirectories) &&
		len(dir.Files) == len(other.Files)) {
		return false
	}
	for subDirectoryName, subDirectory := range dir.SubDirectories {
		if otherSubDir, ok := other.SubDirectories[subDirectoryName]; ok {
			if !otherSubDir.StructureEquals(subDirectory) {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func (dir *Directory) IsSubPath(fullPath string) bool {
	fullPath = filepath.Clean(fullPath)
	currentDir := filepath.Join(dir.Path, dir.Name)
	relPath := strings.TrimPrefix(fullPath, currentDir)
	return relPath != fullPath
}

func (dir *Directory) AddDirectory(fullPath string) (*Directory, error) {
	path, name := filepath.Split(fullPath)
	path = filepath.Clean(path)
	if !dir.IsSubPath(fullPath) {
		return nil, errors.New("fullPath must be a subdirectory of the directory to which it is being added")
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

func (dir Directory) FindDirectory(fullPath string) (*Directory, error) {
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
	return dir.find(pathSlice)
}

func (dir Directory) find(relativePath []string) (*Directory, error) {
	if subDir, ok := dir.SubDirectories[relativePath[0]]; ok {
		if len(relativePath) == 1 {
			return subDir, nil
		}
		return subDir.find(relativePath[1:])
	}
	return nil, errors.New(fmt.Sprintf("directory could not be found. "+
		"Current dir: %s Looking for: %s", dir.Path, strings.Join(relativePath, string(os.PathSeparator))))
}

func GetDirectoryStructure(fullPath string) (*Directory, error) {
	path, name := filepath.Split(fullPath)
	root := Directory{Name: name, Path: path}
	err := filepath.Walk(fullPath,
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				_, err = root.AddDirectory(info.Name())
				if err != nil {
					return err
				}
			} else {
				_, err = root.AddFile(info.Name())
				if err != nil {
					return err
				}
			}
			return nil
		})
	return &root, err
}
