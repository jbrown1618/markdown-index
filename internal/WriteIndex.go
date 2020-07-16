package internal

import "strings"

// WriteIndex produces the string contents of a markdown index given an index model
func WriteIndex(index DirectoryIndex) string {
	lines := make([]string, 0)
	lines = writeDirectoryIndex(lines, index, 0)
	return strings.Join(lines, "\n")
}

func writeDirectoryIndex(lines []string, dir DirectoryIndex, indentationLevel int) []string {
	prefix := "#"
	for i := 0; i < indentationLevel; i++ {
		prefix += "#"
	}
	lines = append(lines, prefix+" "+dir.Name)

	for _, file := range dir.MarkdownFiles {
		lines = writeFileIndex(lines, file)
	}

	for _, subDir := range dir.SubDirectories {
		lines = writeDirectoryIndex(lines, subDir, indentationLevel+1)
	}

	return lines
}

func writeFileIndex(lines []string, file FileIndex) []string {
	return append(lines, "- ["+file.Title+"]("+file.Path+")")
}
