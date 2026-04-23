package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jbrown1618/markdown-index/internal"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "markdown-index",
	Run: func(cmd *cobra.Command, args []string) {
		rootDir, err := cmd.Flags().GetString("root")
		if err != nil {
			log.Fatal(err)
			return
		}
		indexFileName, err := cmd.Flags().GetString("out")
		if err != nil {
			log.Fatal(err)
			return
		}
		openBrowser, err := cmd.Flags().GetBool("browser")
		if err != nil {
			log.Fatal(err)
			return
		}

		absRoot, err := filepath.Abs(rootDir)
		if err != nil {
			log.Fatal(err)
			return
		}

		internal.Skip(indexFileName)
		internal.Skip(".git")

		contents, err := internal.MakeIndex(absRoot)
		if err != nil {
			log.Fatal(err)
			return
		}

		outPath := filepath.Join(absRoot, indexFileName)
		indexFile, err := os.Create(outPath)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer indexFile.Close()

		indexFile.WriteString(contents)

		if openBrowser {
			open.Run(outPath)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().Bool("browser", false, "Open the index file in a web browser")
	rootCmd.PersistentFlags().String("root", "./", "Specify the root directory for which to create an index")
	rootCmd.PersistentFlags().String("out", "index.md", "Specify the name of the index file to create")
}

// Execute is the entry point for the root command.
// It will parse the command line arguments and create the index.md file.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
