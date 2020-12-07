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

	deepCnt := 0
	deep(targetDir, isAll, deepLevel, deepCnt, buf)

	if outFileName == "" {
		io.Copy(os.Stdout, buf)
	} else {
		f, err := os.Create(outFileName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer f.Close()
		io.Copy(f, buf)
	}
}

func deep(targetDir string, isAll bool, deepLevel, deepCnt int, buf *bytes.Buffer) {
	if deepCnt >= deepLevel {
		return
	}

	files, err := ioutil.ReadDir(targetDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	deepCnt++
	maxFileNum := len(files) - 1
	for i, f := range files {
		fileName := f.Name()

		row := ""
		if isAll {
			row = rowWithEdge(i, maxFileNum, deepCnt, fileName)
			if err := writeToBuffer(buf, row); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if f.IsDir() {
				deep(targetDir+"/"+fileName, isAll, deepLevel, deepCnt, buf)
			}
		} else {
			if !strings.HasPrefix(fileName, ".") {
				row = rowWithEdge(i, maxFileNum, deepCnt, fileName)
				if err := writeToBuffer(buf, row); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				if f.IsDir() {
					deep(targetDir+"/"+fileName, isAll, deepLevel, deepCnt, buf)
				}
			}
		}
	}
}

func writeToBuffer(buf *bytes.Buffer, row string) error {
	_, err := buf.WriteString(row + "\n")
	return err
}

func rowWithEdge(i, maxFileNum, deepCnt int, fileName string) string {
	const (
		space = "   "
		edge0 = "│" + space
		edge1 = "├──" + " "
		edge2 = "└──" + " "
	)

	row := ""
	for i := 1; i < deepCnt; i++ {
		row += edge0
	}

	if i == maxFileNum {
		row += edge2 + fileName
	} else {
		row += edge1 + fileName
	}

	return row
}
