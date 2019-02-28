package structure

import (
	"errors"
	"fmt"
	"path/filepath"
)

type File struct {
	Name string
	Path string
}

func (dir *Directory) AddFile(fullPath string) (*File, error) {
	path, name := filepath.Split(fullPath)
	path = filepath.Clean(path)
	if currentDir := filepath.Join(dir.Path, dir.Name); path != currentDir {
		return nil, errors.New(fmt.Sprintf("fullPath must be an immediate child of the directory to which it "+
			"is being added. currentPath: '%s' fullpath: '%s'", currentDir, fullPath))
	}
	if dir.Files == nil {
		dir.Files = map[string]*File{}
	}
	newFile := File{name, path}
	dir.Files[name] = &newFile
	return &newFile, nil
}

func (dir Directory) FindFile(fullPath string) (*File, error) {
	path, name := filepath.Split(fullPath)
	path = filepath.Clean(path)
	fileDir, err := dir.FindDirectory(path)
	if err != nil {
		return nil, err
	}

	if file, ok := fileDir.Files[name]; ok {
		return file, nil
	}
	return nil, errors.New(fmt.Sprintf("file could not be found in directory '%s'", dir.Path))
}
