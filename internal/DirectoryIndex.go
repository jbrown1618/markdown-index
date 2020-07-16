package internal

// DirectoryIndex models the information needed to construct a markdown index of a directory
type DirectoryIndex struct {
	Name           string
	MarkdownFiles  []FileIndex
	SubDirectories []DirectoryIndex
}

// FileIndex models the information needed to add an index line for a markdown file
type FileIndex struct {
	Path  string
	Title string
}
