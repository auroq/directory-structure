package structure

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type File struct {
	Name string
	Path string
}

type Directory struct {
	Name           string
	Path           string
	SubDirectories map[string]Directory
	Files          map[string]File
}

func (dir Directory) Equals(other Directory) bool {
	if dir.Path != other.Path {
		return false
	}
	for subDirectoryName, subDirectory := range dir.SubDirectories {
		if otherDir, ok := other.SubDirectories[subDirectoryName]; ok {
			otherDir.Equals(subDirectory)
		} else {
			return false
		}
	}

	return true
}

func (dir *Directory) AddDirectory(fullPath string) (newDir Directory, err error) {
	path, name := filepath.Split(fullPath)
	path = filepath.Clean(path)
	if currentDir := filepath.Join(dir.Path, dir.Name); path != currentDir {
		return newDir, errors.New(fmt.Sprintf("fullPath must be an immediate subdirectory of the directory to which it "+
			"is being added. currentPath: '%s' fullpath: '%s'", currentDir, fullPath))
	}
	if dir.SubDirectories == nil {
		dir.SubDirectories = map[string]Directory{}
	}
	newDirectory := Directory{Name: name, Path: path}
	dir.SubDirectories[name] = newDirectory
	return
}

func (dir *Directory) AddFile(fullPath string) (newFile File, err error) {
	path, name := filepath.Split(fullPath)
	path = filepath.Clean(path)
	if currentDir := filepath.Join(dir.Path, dir.Name); path != currentDir {
		return newFile, errors.New(fmt.Sprintf("fullPath must be an immediate child of the directory to which it "+
			"is being added. currentPath: '%s' fullpath: '%s'", currentDir, fullPath))
	}
	if dir.Files == nil {
		dir.Files = map[string]File{}
	}
	newFile = File{name, path}
	dir.Files[name] = newFile
	return
}

func (dir Directory) Find(fullPath string) (foundDir *Directory, err error) {
	path, name := filepath.Split(fullPath)
	path = filepath.Clean(path)
	if currentDir := filepath.Join(dir.Path, dir.Name); len(path) < len(currentDir) || path[:len(currentDir)] != currentDir {
		return foundDir, errors.New(fmt.Sprintf("item '%s' is not found in directory '%s'", fullPath, dir.Path))
	}
	path = path[len(dir.Path):]
	pathSlice := strings.Split(path, strconv.QuoteRune(os.PathSeparator))
	return dir.find(pathSlice, name)
}

func (dir Directory) find(relativePath []string, name string) (foundDir *Directory, err error) {
	if subDir, ok := dir.SubDirectories[relativePath[0]]; ok {
		return subDir.find(relativePath[1:], name)
	}
	return foundDir, errors.New("")
}

func GetDirectoryStructure(fullPath string) (root Directory, err error) {
	path, name := filepath.Split(fullPath)
	root = Directory{Name: name, Path: path}
	err = filepath.Walk(fullPath,
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
	return
}
