package internal

var skippables = make(map[string]bool)

// ShouldSkip returns true if a file should be ignored
func ShouldSkip(fileName string) bool {
	shouldSkip, present := skippables[fileName]
	return present && shouldSkip
}

// Skip will cause the given file name to be skipped
func Skip(fileName string) {
	skippables[fileName] = true
}
