package structure

type Descendants struct {
	Directories []*Directory
	Files       []*File
}

// ContainsDirectory determines whether Descendants contains a Directory
// It returns true or false accordingly
func (desc Descendants) ContainsDirectory(dir *Directory) bool {
	for _, d := range desc.Directories {
		if d.Equals(dir) {
			return true
		}
	}
	return false
}

// ContainsFile determines whether Descendants contains a File
// It returns true or false accordingly
func (desc Descendants) ContainsFile(file *File) bool {
	for _, f := range desc.Files {
		if f.Equals(file) {
			return true
		}
	}
	return false
}

// GetAllDescendants walks through the given directory builds a structure of its descendants
// It returns a Descendants which has two lists: one for all the Directory descendants and
// one for all the File descendants
func (dir Directory) GetAllDescendants() Descendants {
	descDirs, descFiles := dir.getDescendants()
	return Descendants{Directories: descDirs, Files: descFiles}
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
