package structure

var DirectoryIdentities = []struct {
	name        string
	dir         Directory
	dirFullPath string
}{
	{
		"EmptyDirectory",
		Directory{name: "dir1", path: "/tmp"},
		"/tmp/dir1",
	},
	{"DirectoryWithSubDirectory",
		Directory{
			name: "dir1",
			path: "/tmp",
			subDirectories: map[string]*Directory{
				"sub1": {
					name: "sub1",
					path: "/tmp/dir1",
				},
			},
		},
		"/tmp/dir1",
	},
	{"DirectoryWithSubDirectories",
		Directory{
			name: "dir1",
			path: "/tmp",
			subDirectories: map[string]*Directory{
				"sub1": {
					name: "sub1",
					path: "/tmp/dir1",
				},
				"sub2": {
					name: "sub2",
					path: "/tmp/dir1",
				},
				"sub3": {
					name: "sub3",
					path: "/tmp/dir1",
				},
			},
		},
		"/tmp/dir1/sub2",
	},
	{"DirectoryWithSubDirectories2",
		Directory{
			name: "dir1",
			path: "/tmp",
			subDirectories: map[string]*Directory{
				"sub-1": {
					name: "sub-1",
					path: "/tmp/dir1",
				},
				"sub-2": {
					name: "sub-2",
					path: "/tmp/dir1",
				},
				"sub-3": {
					name: "sub-3",
					path: "/tmp/dir1",
				},
			},
		},
		"/tmp/dir1/sub-2",
	},
	{"DirectoryWithSubDirectoryWithSubDirectory",
		Directory{
			name: "dir1",
			path: "/tmp",
			subDirectories: map[string]*Directory{
				"sub1": {
					name: "sub1",
					path: "/tmp/dir1",
				},
				"sub2": {
					name: "sub2",
					path: "/tmp/dir1",
					subDirectories: map[string]*Directory{
						"subsub1": {
							name: "subsub1",
							path: "/tmp/dir1/sub2",
						},
					},
				},
				"sub3": {
					name: "sub3",
					path: "/tmp/dir1",
				},
			},
		},
		"/tmp/dir1/sub2/subsub1",
	},
	{"DirectoryWithSubDirectoryWithSubDirectories",
		Directory{
			name: "dir1",
			path: "/tmp",
			subDirectories: map[string]*Directory{
				"sub1": {
					name: "sub1",
					path: "/tmp/dir1",
				},
				"sub2": {
					name: "sub2",
					path: "/tmp/dir1",
					subDirectories: map[string]*Directory{
						"subsub1": {
							name: "subsub1",
							path: "/tmp/dir1/sub2",
						},
						"subsub2": {
							name: "subsub2",
							path: "/tmp/dir1/sub2",
						},
						"subsub3": {
							name: "subsub3",
							path: "/tmp/dir1/sub2",
						},
					},
				},
				"sub3": {
					name: "sub3",
					path: "/tmp/dir1",
				},
			},
		},
		"/tmp/dir1/sub2/subsub2",
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
			dir := NewDirectory("dir1", "/tmp")
			_, _ = dir.AddDirectory("/tmp/dir1/sub1")
			_, _ = dir.AddFile("/tmp/dir1/sub1.txt")
			return dir
		}(),
		"/tmp/dir1/sub1",
		"/tmp/dir1/sub1.txt",
	},
	{"DirectoryWithSubDirectoryWithSubDirectory",
		func() Directory {
			dir := NewDirectory("dir1", "/tmp")
			sub1, _ := dir.AddDirectory("/tmp/dir1/sub1")
			_, _ = sub1.AddDirectory("/tmp/dir1/sub1/subsub1")
			_, _ = sub1.AddFile("/tmp/dir1/sub1/subsub1.txt")
			return dir
		}(),
		"/tmp/dir1/sub1/subsub1",
		"/tmp/dir1/sub1/subsub1.txt",
	},
	{"DirectoryWithSubDirectories",
		func() Directory {
			dir := NewDirectory("dir1", "/tmp")
			sub1, _ := dir.AddDirectory("/tmp/dir1/sub1")
			_, _ = sub1.AddDirectory("/tmp/dir1/sub1/subsub1")
			subsub2, _ := sub1.AddDirectory("/tmp/dir1/sub1/subsub2")
			_, _ = sub1.AddDirectory("/tmp/dir1/sub1/subsub3")
			_, _ = subsub2.AddDirectory("/tmp/dir1/sub1/subsub2/subsubsub")
			_, _ = subsub2.AddFile("/tmp/dir1/sub1/subsub2/subsubsub.txt")
			return dir
		}(),
		"/tmp/dir1/sub1/subsub2/subsubsub",
		"/tmp/dir1/sub1/subsub2/subsubsub.txt",
	},
}

var cleanTests = []struct {
	name string
	path string
	cleanPath string
} {
	{
		"trailingSlash",
		"/tmp/",
		"/tmp",
	},
	{
		"extraCenterSlash",
		"/tmp//dir",
		"/tmp/dir",
	},
	{
		"extraLeadingSlash",
		"//tmp/dir",
		"/tmp/dir",
	},
}

var fullPathTests = []struct {
	testName string
	name string
	path string
	fullPath string
} {
	{
		"pathIsRoot",
		"item",
		"/",
		"/item",
	},
	{
		"fullPathIsRoot",
		"/",
		"",
		"/",
	},
	{
		"trailingSlashInPath",
		"item",
		"/tmp/",
		"/tmp/item",
	},
	{
		"extraCenterSlashInPath",
		"item",
		"/tmp//dir",
		"/tmp/dir/item",
	},
	{
		"extraLeadingSlashInPath",
		"item",
		"//tmp/dir",
		"/tmp/dir/item",
	},
}
