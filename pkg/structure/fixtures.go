package structure

var DirectoryIdentities = []struct {
	name         string
	dir          Directory
	dirFullPath  string
}{
	{
		"EmptyDirectory",
		Directory{Name: "dir1", Path: "/tmp"},
		"/tmp/dir1",
	},
	{"DirectoryWithSubDirectory",
		Directory{
			Name: "dir1",
			Path: "/tmp",
			SubDirectories: map[string]*Directory{
				"sub1": {
					Name: "sub1",
					Path: "/tmp/dir1",
				},
			},
		},
		"/tmp/dir1",
	},
	{"DirectoryWithSubDirectories",
		Directory{
			Name: "dir1",
			Path: "/tmp",
			SubDirectories: map[string]*Directory{
				"sub1": {
					Name: "sub1",
					Path: "/tmp/dir1",
				},
				"sub2": {
					Name: "sub2",
					Path: "/tmp/dir1",
				},
				"sub3": {
					Name: "sub3",
					Path: "/tmp/dir1",
				},
			},
		},
		"/tmp/dir1/sub2",
	},
	{"DirectoryWithSubDirectoryWithSubDirectory",
		Directory{
			Name: "dir1",
			Path: "/tmp",
			SubDirectories: map[string]*Directory{
				"sub1": {
					Name: "sub1",
					Path: "/tmp/dir1",
				},
				"sub2": {
					Name: "sub2",
					Path: "/tmp/dir1",
					SubDirectories: map[string]*Directory{
						"subsub1": {
							Name: "subsub1",
							Path: "/tmp/dir1/sub2",
						},
					},
				},
				"sub3": {
					Name: "sub3",
					Path: "/tmp/dir1",
				},
			},
		},
		"/tmp/dir1/sub2/subsub1",
	},
	{"DirectoryWithSubDirectoryWithSubDirectories",
		Directory{
			Name: "dir1",
			Path: "/tmp",
			SubDirectories: map[string]*Directory{
				"sub1": {
					Name: "sub1",
					Path: "/tmp/dir1",
				},
				"sub2": {
					Name: "sub2",
					Path: "/tmp/dir1",
					SubDirectories: map[string]*Directory{
						"subsub1": {
							Name: "subsub1",
							Path: "/tmp/dir1/sub2",
						},
						"subsub2": {
							Name: "subsub2",
							Path: "/tmp/dir1/sub2",
						},
						"subsub3": {
							Name: "subsub3",
							Path: "/tmp/dir1/sub2",
						},
					},
				},
				"sub3": {
					Name: "sub3",
					Path: "/tmp/dir1",
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
			dir := Directory{Name: "dir1", Path: "/tmp"}
			_, _ = dir.AddDirectory("/tmp/dir1/sub1")
			_, _ = dir.AddFile("/tmp/dir1/sub1.txt")
			return dir
		}(),
		"/tmp/dir1/sub1",
		"/tmp/dir1/sub1.txt",
	},
	{"DirectoryWithSubDirectoryWithSubDirectory",
		func() Directory {
			dir := Directory{Name: "dir1", Path: "/tmp"}
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
			dir := Directory{Name: "dir1", Path: "/tmp"}
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
