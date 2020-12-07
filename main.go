package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
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

	buf := &bytes.Buffer{}

	if err := writeToBuffer(buf, targetDir); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	files, err := ioutil.ReadDir(targetDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	maxFileNum := len(files) - 1
	for i, f := range files {
		fileName := f.Name()

		row := ""
		if isAll {
			row = rowWithEdge(i, maxFileNum, fileName)
			if err := writeToBuffer(buf, row); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			if !strings.HasPrefix(fileName, ".") {
				row = rowWithEdge(i, maxFileNum, fileName)
				if err := writeToBuffer(buf, row); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		}
	}

	io.Copy(os.Stdout, buf)
}

func writeToBuffer(buf *bytes.Buffer, row string) error {
	_, err := buf.WriteString(row + "\n")
	return err
}

func rowWithEdge(i, maxFileNum int, fileName string) string {
	const (
		edge1 = "├──" + " "
		edge2 = "└──" + " "
	)

	if i == maxFileNum {
		return edge2 + fileName
	} else {
		return edge1 + fileName
	}
}
