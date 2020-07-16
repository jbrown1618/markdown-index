package internal

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// MakeIndex takes a directory path and returns the contents of
// a file which indexes the markdown files in that directory
func MakeIndex(dirPath string) (string, error) {
	os.Chdir(dirPath)

	index, err := indexDirectory(dirPath, dirPath)
	if err != nil {
		return "", err
	}

	return WriteIndex(index), nil
}

func indexDirectory(rootPath string, dirPath string) (DirectoryIndex, error) {
	var index DirectoryIndex
	var subDirectories []DirectoryIndex
	var markdownFiles []FileIndex

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return index, err
	}

	for _, info := range files {
		name := info.Name()
		if ShouldSkip(name) {
			continue
		}

		if info.IsDir() {
			subDirIndex, err := indexDirectory(rootPath, path.Join(dirPath, name))
			if err != nil {
				continue
			}
			subDirectories = append(subDirectories, subDirIndex)
		} else {
			ext := filepath.Ext(name)
			if ext != ".md" {
				continue
			}
			fileIndex := indexFile(rootPath, path.Join(dirPath, name))
			markdownFiles = append(markdownFiles, fileIndex)
		}
	}

	if len(subDirectories) == 0 && len(markdownFiles) == 0 {
		return index, errors.New("Empty directory")
	}

	return DirectoryIndex{
		Name:           getDirectoryTitle(dirPath),
		SubDirectories: subDirectories,
		MarkdownFiles:  markdownFiles,
	}, nil
}

func indexFile(rootPath, filePath string) FileIndex {
	path := strings.Replace(filePath, rootPath, ".", 1)
	return FileIndex{
		Title: getFileTitle(filePath),
		Path:  path,
	}
}

func getDirectoryTitle(dirPath string) string {
	_, dirName := path.Split(dirPath)
	re := regexp.MustCompile(`[\-\_]`)
	return strings.Title(re.ReplaceAllString(dirName, " "))
}

func getFileTitle(filePath string) string {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return filePath
	}
	file, err := os.Open(absPath)
	defer file.Close()
	if err != nil {
		return filePath
	}

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
