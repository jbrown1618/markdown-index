package main

import (
	"github.com/jbrown1618/markdown-index/internal"
	"github.com/skratchdot/open-golang/open"
	"log"
	"os"
)

func main() {
	index := internal.MakeIndex("/Users/joabro/Documents/Notes")

	indexFile, err := os.Create("index.md")
	defer indexFile.Close()
	if err != nil {
		log.Fatal(err)
		return
	}

	indexFile.WriteString(index)
	open.RunWith("index.md", "Google Chrome")
}
