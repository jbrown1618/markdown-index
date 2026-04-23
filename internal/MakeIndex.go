package internal

import (
	"bufio"
	"errors"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// MakeIndex takes a directory path and returns the contents of
// a file which indexes the markdown files in that directory
func MakeIndex(dirPath string) (string, error) {
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return "", err
	}

	index, err := indexDirectory(absPath, absPath)
	if err != nil {
		return "", err
	}

	return WriteIndex(index), nil
}

func indexDirectory(rootPath string, dirPath string) (DirectoryIndex, error) {
	var index DirectoryIndex

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return index, err
	}

	// Separate entries into dirs and markdown files, preserving order
	type dirEntry struct {
		idx  int
		path string
	}
	type fileEntry struct {
		idx  int
		path string
	}

	var dirEntries []dirEntry
	var fileEntries []fileEntry

	for _, entry := range entries {
		name := entry.Name()
		if ShouldSkip(name) {
			continue
		}
		if entry.IsDir() {
			dirEntries = append(dirEntries, dirEntry{len(dirEntries), path.Join(dirPath, name)})
		} else if filepath.Ext(name) == ".md" {
			fileEntries = append(fileEntries, fileEntry{len(fileEntries), path.Join(dirPath, name)})
		}
	}

	// Index subdirectories concurrently
	subDirResults := make([]DirectoryIndex, len(dirEntries))
	subDirErrors := make([]error, len(dirEntries))
	var wg sync.WaitGroup

	for _, d := range dirEntries {
		wg.Add(1)
		go func(i int, p string) {
			defer wg.Done()
			subDirResults[i], subDirErrors[i] = indexDirectory(rootPath, p)
		}(d.idx, d.path)
	}

	// Index markdown files concurrently
	fileResults := make([]FileIndex, len(fileEntries))
	for _, f := range fileEntries {
		wg.Add(1)
		go func(i int, p string) {
			defer wg.Done()
			fileResults[i] = indexFile(rootPath, p)
		}(f.idx, f.path)
	}

	wg.Wait()

	// Collect results in original order, skipping failed subdirectories
	var subDirectories []DirectoryIndex
	for i, result := range subDirResults {
		if subDirErrors[i] == nil {
			subDirectories = append(subDirectories, result)
		}
	}

	if len(subDirectories) == 0 && len(fileResults) == 0 {
		return index, errors.New("empty directory")
	}

	return DirectoryIndex{
		Name:           getDirectoryTitle(dirPath),
		SubDirectories: subDirectories,
		MarkdownFiles:  fileResults,
	}, nil
}

func indexFile(rootPath, filePath string) FileIndex {
	p := strings.Replace(filePath, rootPath, ".", 1)
	return FileIndex{
		Title: getFileTitle(filePath),
		Path:  p,
	}
}

func getDirectoryTitle(dirPath string) string {
	_, dirName := path.Split(dirPath)
	re := regexp.MustCompile(`[\-\_]`)
	caser := cases.Title(language.English, cases.NoLower)
	return caser.String(re.ReplaceAllString(dirName, " "))
}

func getFileTitle(filePath string) string {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return filePath
	}
	file, err := os.Open(absPath)
	if err != nil {
		return filePath
	}
	defer file.Close()

	titlePattern := regexp.MustCompile(`^\s?#\s`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if titlePattern.MatchString(line) {
			title := titlePattern.ReplaceAllString(line, "")
			return title
		}
	}

	return filePath
}
