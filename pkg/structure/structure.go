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
	Files          []File
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

func (dir Directory) AddFileOrDir(info os.FileInfo) {
	fullPath := info.Name()
	path, name := filepath.Split(fullPath)
	if info.IsDir() {
		if dir.SubDirectories == nil {
			dir.SubDirectories = map[string]Directory{}
		}
		dir.SubDirectories[info.Name()] = Directory{
			name,
			path,
			nil,
			nil,
		}
	} else {
		dir.Files = append(dir.Files, File{name, path})
	}
}

func (dir Directory) Find(fullPath string) (foundDir Directory, err error) {
	path, name := filepath.Split(fullPath)
	if path[:len(dir.Path)] != dir.Path {
		return foundDir, errors.New(fmt.Sprintf("item '%s' is not found in directory '%s'", fullPath, dir.Path))
	}
	path = path[len(dir.Path):]
	pathSlice := strings.Split(path, strconv.QuoteRune(os.PathSeparator))
	return dir.find(pathSlice, name)
}

func (dir Directory) find(relativePath []string, name string) (foundDir Directory, err error) {
	if subDir, ok := dir.SubDirectories[relativePath[0]]; ok {
		return subDir.find(relativePath[1:], name)
	}
	return foundDir, errors.New("")
}

func GetDirectoryStructure(fullPath string) (root Directory, err error) {
	path, name := filepath.Split(fullPath)
	root = Directory{
		name,
		path,
		nil,
		nil,
	}
	err = filepath.Walk(fullPath,
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			root.AddFileOrDir(info)
			return nil
		})
	return
}
