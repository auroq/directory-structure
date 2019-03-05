# Directory Structure

Directory Structure is a go library for reading and transversing your local filesystem.
Currently it consists of one package: structure

__NOTE: This library does not do any modification to the actual local filesystem.__


## Parts of [Structure][Structure]

Structure provides three structs: Directory, File, and Descendants


### [Directory][Directory]

Directory is the core of structure.  It is a tree where each node contains the following:
- [Name()][Directory.Name]: the name of the directory 
- [Path()][Directory.Path]: the path to the directory _excluding_ Name()
- [FullPath()][Directory.FullPath]: the full path to the directory _including_ Name()
- [Files()][Directory.Files]: pointers to the files contained in the directory
- [SubDirectories()][Directory.SubDirectories]: pointers to Directories in this Directory


### [File][File]

File represents a single file. It consists of the following:
- [Name()][File.Name]: the name of the directory 
- [Path()][File.Path]: the path to the directory _excluding_ Name()
- [FullPath()][File.FullPath]: the full path to the directory _including_ Name()


### [Descendants][Descendants]

Descendants is used for convenience when getting all the Files and Directories that are descendants of a Directory. It consists of the following:
- Directories: the Directories that are Descendants of a Directory
- Files: the Files that are Descendants of a Directory


## Functionality of [Structure][Structure]

### Directory Tree Creation

There are two main ways to create a Directory tree
1. Call [NewDirectory()][Structure.NewDirectory]:
This create a new directory with no files or subdirectories.
2. Call [GetDirectoryStructure()][Structure.GetDirectoryStructure]:
This walks your local filesystem at the path provided and generates a full Directory tree that matches the given directory.


### Adding Items to a Directory Tree

Directories and Files can be added to a directory tree by calling either [directory.AddDirectory()][Directory.AddDirectory] or [directory.AddFile()][Directory.AddDirectory]
on any directory in the tree that is a ancestor of the new item.
Either method takes the full path to the new item, and will create directories as needed between the ancestor Directory and the new item.


### Searching a Directory or Directory Tree

#### Direct Child

A Directory or File that is a direct child of the current directory can be found by calling [directory.Directory()][Directory.Directory] or [directory.File()][Directory.File] on the current directory and passing it the name of the Directory of File.


#### Descendant By Path

Descendants of a Directory can be found using the Descendant's full path by calling either [directory.GetDirectory()][Directory.GetDirectory] or [directory.GetFile()][Directory.GetFile]


#### Descendant By Name

A __depth first search__ by name can be done by calling [directory.FindDirectoryDepth()][Directory.FindDirectoryDepth] or [directory.FindFileDepth()][Directory.FindFileDepth]

A __breadth first search__ by name can be done by calling [directory.FindDirectoryBreadth()][Directory.FindDirectoryBreadth] or [directory.FindFileBreadth()][Directory.FindDirectoryBreadth]


[Structure]: https://godoc.org/github.com/auroq/directory-structure/pkg/structure
[Structure.NewDirectory]: https://godoc.org/github.com/auroq/directory-structure/pkg/structure#NewDirectory
[Structure.GetDirectoryStructure]: https://godoc.org/github.com/auroq/directory-structure/pkg/structure#GetDirectoryStructurey

[Directory]: https://godoc.org/github.com/auroq/directory-structure/pkg/structure#Directory
[Directory.Name]: https://godoc.org/github.com/auroq/directory-structure/pkg/structure#Directory.Name
[Directory.Path]: https://godoc.org/github.com/auroq/directory-structure/pkg/structure#Directory.Path
[Directory.FullPath]: https://godoc.org/github.com/auroq/directory-structure/pkg/structure#Directory.FullPath
[Directory.Files]: https://godoc.org/github.com/auroq/directory-structure/pkg/structure#Directory.Files
[Directory.SubDirectories]: https://godoc.org/github.com/auroq/directory-structure/pkg/structure#Directory.SubDirectories
[Directory.AddDirectory]: https://godoc.org/github.com/auroq/directory-structure/pkg/structure#Directory.AddDirectory
[Directory.AddFile]: https://godoc.org/github.com/auroq/directory-structure/pkg/structure#Directory.AddFile
[Directory.Directory]: https://godoc.org/github.com/auroq/directory-structure/pkg/structure#Directory.Directory
[Directory.File]: https://godoc.org/github.com/auroq/directory-structure/pkg/structure#Directory.File
[Directory.GetDirectory]: https://godoc.org/github.com/auroq/directory-structure/pkg/structure#Directory.GetDirectory
[Directory.GetFile]: https://godoc.org/github.com/auroq/directory-structure/pkg/structure#Directory.GetFile
[Directory.FindDirectoryDepth]: https://godoc.org/github.com/auroq/directory-structure/pkg/structure#Directory.FindDirectoryDepth
[Directory.FindFileDepth]: https://godoc.org/github.com/auroq/directory-structure/pkg/structure#Directory.FindFileDepth
[Directory.FindDirectoryBreadth]: https://godoc.org/github.com/auroq/directory-structure/pkg/structure#Directory.FindDirectoryBreadth
[Directory.FindFileBreadth]: https://godoc.org/github.com/auroq/directory-structure/pkg/structure#Directory.FindFileBreadth

[File]: https://godoc.org/github.com/auroq/directory-structure/pkg/structure#File
[File.Name]: https://godoc.org/github.com/auroq/directory-structure/pkg/structure#File.Name
[File.Path]: https://godoc.org/github.com/auroq/directory-structure/pkg/structure#File.Path
[File.FullPath]: https://godoc.org/github.com/auroq/directory-structure/pkg/structure#File.FullPath

[Descendants]: https://godoc.org/github.com/auroq/directory-structure/pkg/structure#Descendants