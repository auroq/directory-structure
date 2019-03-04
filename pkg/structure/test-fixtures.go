package structure

import (
	"os"
	"path/filepath"
	"runtime"
)

func osRoot() string {
	switch runtime.GOOS {
	case "windows":
		return "C:"
	case "linux":
		return "/"
	default:
		return "/"
	}
}

var DirectoryIdentities = []struct {
	name        string
	dir         Directory
	dirFullPath string
}{
	{
		"EmptyDirectory",
		Directory{name: "dir1", path: osRoot() + "tmp"},
		filepath.Join(osRoot(), "tmp", "dir1"),
	},
	{"DirectoryWithSubDirectory",
		Directory{
			name: "dir1",
			path: osRoot() + "tmp",
			subDirectories: map[string]*Directory{
				"sub1": {
					name: "sub1",
					path: filepath.Join(osRoot(), "tmp", "dir1"),
				},
			},
		},
		filepath.Join(osRoot(), "tmp", "dir1"),
	},
	{"DirectoryWithSubDirectories",
		Directory{
			name: "dir1",
			path: osRoot() + "tmp",
			subDirectories: map[string]*Directory{
				"sub1": {
					name: "sub1",
					path: filepath.Join(osRoot(), "tmp", "dir1"),
				},
				"sub2": {
					name: "sub2",
					path: filepath.Join(osRoot(), "tmp", "dir1"),
				},
				"sub3": {
					name: "sub3",
					path: filepath.Join(osRoot(), "tmp", "dir1"),
				},
			},
		},
		filepath.Join(osRoot(), "tmp", "dir1", "sub2"),
	},
	{"DirectoryWithSubDirectories2",
		Directory{
			name: "dir1",
			path: osRoot() + "tmp",
			subDirectories: map[string]*Directory{
				"sub-1": {
					name: "sub-1",
					path: filepath.Join(osRoot(), "tmp", "dir1"),
				},
				"sub-2": {
					name: "sub-2",
					path: filepath.Join(osRoot(), "tmp", "dir1"),
				},
				"sub-3": {
					name: "sub-3",
					path: filepath.Join(osRoot(), "tmp", "dir1"),
				},
			},
		},
		filepath.Join(osRoot(), "tmp", "dir1", "sub-2"),
	},
	{"DirectoryWithSubDirectoryWithSubDirectory",
		Directory{
			name: "dir1",
			path: osRoot() + "tmp",
			subDirectories: map[string]*Directory{
				"sub1": {
					name: "sub1",
					path: filepath.Join(osRoot(), "tmp", "dir1"),
				},
				"sub2": {
					name: "sub2",
					path: filepath.Join(osRoot(), "tmp", "dir1"),
					subDirectories: map[string]*Directory{
						"subsub1": {
							name: "subsub1",
							path: filepath.Join(osRoot(), "tmp", "dir1", "sub2"),
						},
					},
				},
				"sub3": {
					name: "sub3",
					path: filepath.Join(osRoot(), "tmp", "dir1"),
				},
			},
		},
		filepath.Join(osRoot(), "tmp", "dir1", "sub2", "subsub1"),
	},
	{"DirectoryWithSubDirectoryWithSubDirectories",
		Directory{
			name: "dir1",
			path: osRoot() + "tmp",
			subDirectories: map[string]*Directory{
				"sub1": {
					name: "sub1",
					path: filepath.Join(osRoot(), "tmp", "dir1"),
				},
				"sub2": {
					name: "sub2",
					path: filepath.Join(osRoot(), "tmp", "dir1"),
					subDirectories: map[string]*Directory{
						"subsub1": {
							name: "subsub1",
							path: filepath.Join(osRoot(), "tmp", "dir1", "sub2"),
						},
						"subsub2": {
							name: "subsub2",
							path: filepath.Join(osRoot(), "tmp", "dir1", "sub2"),
						},
						"subsub3": {
							name: "subsub3",
							path: filepath.Join(osRoot(), "tmp", "dir1", "sub2"),
						},
					},
				},
				"sub3": {
					name: "sub3",
					path: filepath.Join(osRoot(), "tmp", "dir1"),
				},
			},
		},
		filepath.Join(osRoot(), "tmp", "dir1", "sub2", "subsub2"),
	},
}

var FindTests = []struct {
	name               string
	dir                Directory
	fullDirPathToFind  string
	fullFilePathToFind string
}{
	{"DirectoryWithSubDirectory",
		func() Directory {
			dir := NewDirectory("dir1", filepath.Join(osRoot(), "tmp"))
			_, _ = dir.AddDirectory(filepath.Join(osRoot(), "tmp", "dir1", "sub1"))
			_, _ = dir.AddFile(filepath.Join(osRoot(), "tmp", "dir1", "sub1.txt"))
			return dir
		}(),
		filepath.Join(osRoot(), "tmp", "dir1", "sub1"),
		filepath.Join(osRoot(), "tmp", "dir1", "sub1.txt"),
	},
	{"DirectoryWithSubDirectoryWithSubDirectory",
		func() Directory {
			dir := NewDirectory("dir1", filepath.Join(osRoot(), "tmp"))
			sub1, _ := dir.AddDirectory(filepath.Join(osRoot(), "tmp", "dir1", "sub1"))
			_, _ = sub1.AddDirectory(filepath.Join(osRoot(), "tmp", "dir1", "sub1", "subsub1"))
			_, _ = sub1.AddFile(filepath.Join(osRoot(), "tmp", "dir1", "sub1", "subsub1.txt"))
			return dir
		}(),
		filepath.Join(osRoot(), "tmp", "dir1", "sub1", "subsub1"),
		filepath.Join(osRoot(), "tmp", "dir1", "sub1", "subsub1.txt"),
	},
	{"DirectoryWithSubDirectories",
		func() Directory {
			dir := NewDirectory("dir1", filepath.Join(osRoot(), "tmp"))
			sub1, _ := dir.AddDirectory(filepath.Join(osRoot(), "tmp", "dir1", "sub1"))
			_, _ = sub1.AddDirectory(filepath.Join(osRoot(), "tmp", "dir1", "sub1", "subsub1"))
			subsub2, _ := sub1.AddDirectory(filepath.Join(osRoot(), "tmp", "dir1", "sub1", "subsub2"))
			_, _ = sub1.AddDirectory(filepath.Join(osRoot(), "tmp", "dir1", "sub1", "subsub3"))
			_, _ = subsub2.AddDirectory(filepath.Join(osRoot(), "tmp", "dir1", "sub1", "subsub2", "subsubsub"))
			_, _ = subsub2.AddFile(filepath.Join(osRoot(), "tmp", "dir1", "sub1", "subsub2", "subsubsub.txt"))
			return dir
		}(),
		filepath.Join(osRoot(), "tmp", "dir1", "sub1", "subsub2", "subsubsub"),
		filepath.Join(osRoot(), "tmp", "dir1", "sub1", "subsub2", "subsubsub.txt"),
	},
}

var cleanTests = []struct {
	name      string
	path      string
	cleanPath string
}{
	{
		"trailingSlash",
		osRoot() + "tmp" + string(os.PathSeparator),
		osRoot() + "tmp",
	},
	{
		"extraCenterSlash",
		osRoot() + "tmp" + string(os.PathSeparator) + string(os.PathSeparator) + "dir",
		filepath.Join(osRoot(), "tmp", "dir"),
	},
	{
		"extraLeadingSlash",
		"//tmp/dir",
		func() string {
			if runtime.GOOS == "windows" {
				return "\\\\tmp\\dir"
			}
			return "/tmp/dir"
		}(),
	},
}

var fullPathTests = []struct {
	testName string
	name     string
	path     string
	fullPath string
}{
	{
		"pathIsRoot",
		"item",
		osRoot(),
		filepath.Join(osRoot(), "item"),
	},
	{
		"fullPathIsRoot",
		osRoot(),
		"",
		func() string {
			if runtime.GOOS == "windows" {
				return osRoot() + "."
			}
			return osRoot()
		}(),
	},
	{
		"trailingSlashInPath",
		"item",
		osRoot() + "tmp" + string(os.PathSeparator),
		filepath.Join(osRoot(), "tmp", "item"),
	},
	{
		"extraCenterSlashInPath",
		"item",
		osRoot() + "tmp" + string(os.PathSeparator) + string(os.PathSeparator) + "dir",
		filepath.Join(osRoot(), "tmp", "dir", "item"),
	},
	{
		"extraLeadingSlashInPath",
		"item",
		"//tmp/dir",
		func() string {
			if runtime.GOOS == "windows" {
				return "\\\\tmp\\dir\\item"
			}
			return "/tmp/dir"
		}(),
	},
}
