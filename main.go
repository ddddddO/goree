package main

import (
	"flag"

	"github.com/ddddddO/goree/tree"
)

func main() {
	var (
		isAll       bool   // tree -a
		isDirOnly   bool   // tree -d
		isFullPath  bool   // tree -f
		outFileName string // tree -o <filename>
		deepLevel   int    // tree -L <integer>
	)
	flag.BoolVar(&isAll, "a", false, "All files are listed.")
	flag.BoolVar(&isDirOnly, "d", false, "List directories only.")
	flag.BoolVar(&isFullPath, "f", false, "Print the full path prefix for each file.")
	flag.StringVar(&outFileName, "o", "", "Output to file instead of stdout.")
	flag.IntVar(&deepLevel, "L", 1, "Descend only level directories deep.")
	flag.Parse()

	// tree表示対象のディレクトリ
	args := flag.Args()
	targetDir := "."
	if len(args) > 0 {
		targetDir = args[0]
	}

	tree.Run(targetDir, outFileName, isAll, deepLevel)
}
