package cmd

import (
	"fmt"
	"github.com/jbrown1618/markdown-index/internal"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var indexFileName = "index.md"

var rootCmd = &cobra.Command{
	Use: "markdown-index",
	Run: func(cmd *cobra.Command, args []string) {
		rootDir, err := cmd.Flags().GetString("root")
		if err != nil {
			log.Fatal(err)
			return
		}
		openBrowser, err := cmd.Flags().GetBool("browser")
		if err != nil {
			log.Fatal(err)
			return
		}

		os.Chdir(rootDir)
		indexFile, err := os.Create(indexFileName)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer indexFile.Close()

		indexFile.WriteString(internal.MakeIndex(rootDir))
		if openBrowser {
			open.RunWith(indexFileName, "Google Chrome")
		}
	},
}

func init() {
	rootCmd.PersistentFlags().Bool("browser", false, "Open the index file in a web browser")
	rootCmd.PersistentFlags().String("root", "./", "Specify the root directory for which to create an index")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
