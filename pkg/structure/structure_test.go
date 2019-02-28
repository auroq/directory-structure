package structure

import (
	"testing"
)

var DirectoryIdentities = []struct {
	name string
	dir  Directory
}{
	{"EmptyDirectory",
		Directory{
			"/tmp/dir1",
			"dir1",
			nil,
			nil,
		},
	},
	{"DirectoryWithSubDirectory",
		Directory{
			"/tmp/dir1",
			"dir1",
			map[string]Directory{
				"sub1": {
					"sub1",
					"/tmp/dir1/sub1",
					nil,
					nil,
				},
			},
			nil,
		},
	},
	{"DirectoryWithSubDirectoryWithSubDirectory",
		Directory{
			"/tmp/dir1",
			"dir1",
			map[string]Directory{
				"sub1": {
					"sub1",
					"/tmp/dir1/sub1",
					map[string]Directory{
						"subsub1": {
							"subsub1",
							"/tmp/dir1/sub1/subsub1",
							nil,
							nil,
						},
					},
					nil,
				},
			},
			nil,
		},
	},
}

func TestDirectoryEqualsWithIdentity(t *testing.T) {
	for _, tt := range DirectoryIdentities {
		t.Run(tt.name, func(t *testing.T) {
			tt.dir.Equals(tt.dir)
		})
	}
}
