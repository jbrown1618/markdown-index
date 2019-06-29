package internal

import (
	"bufio"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// MakeIndex takes a directory path and returns the contents of a file which indexes the markdown files in that directory
func MakeIndex(dirPath string) string {
	indexContents := make([]string, 0)
	indexContents = append(indexContents, "# Index")
	filepath.Walk(dirPath, func(currentPath string, info os.FileInfo, err error) error {
		if currentPath == dirPath {
			return nil
		}
		currentPath, err = filepath.Rel(dirPath, currentPath)
		if err != nil {
			return err
		}

		if info.IsDir() {
			indexContents = append(indexContents, indexDirectory(currentPath))
		} else if filepath.Ext(currentPath) == ".md" {
			if info.Name() == "index.md" {
				return nil
			}
			indexContents = append(indexContents, indexFile(currentPath))
		}
		return nil
	})
	return strings.Join(indexContents, "\n")
}

func indexDirectory(dirPath string) string {
	_, dirName := path.Split(dirPath)
	re := regexp.MustCompile(`[\-\_]`)
	header := re.ReplaceAllString(dirName, " ")

	numSeparators := strings.Count(dirPath, string(filepath.Separator))
	prefix := "##"
	for i := 0; i < numSeparators; i++ {
		prefix += "#"
	}
	return prefix + " " + header
}

func indexFile(filePath string) string {
	relativePath := "./" + filePath
	title := getFileTitle(filePath)
	return "- [" + title + "](" + relativePath + ")"
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
