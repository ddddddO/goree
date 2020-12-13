package tree

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// Run is ...
func Run(targetDir, outFileName string, isAll bool, deepLevel int) {
	buf := &bytes.Buffer{}

	if err := writeToBuffer(buf, targetDir); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	deepCnt := 0
	deep(targetDir, isAll, deepLevel, deepCnt, 0, buf)

	if err := output(buf, outFileName); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func deep(targetDir string, isAll bool, deepLevel, deepCnt, parentFileNum int, buf *bytes.Buffer) {
	if deepCnt >= deepLevel {
		return
	}

	files, hiddenFileNum, err := getFilesAndHiddenCnt(targetDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	deepCnt++
	currentFileNum := len(files)
	for i, f := range files {
		fileName := f.Name()

		row := ""
		if isAll {
			row = rowWithEdge(i, currentFileNum, deepCnt, deepLevel, parentFileNum, fileName)
			if err := writeToBuffer(buf, row); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if f.IsDir() {
				deep(targetDir+"/"+fileName, isAll, deepLevel, deepCnt, currentFileNum, buf)
			}
		} else {
			if strings.HasPrefix(fileName, ".") {
				continue
			}

			row = rowWithEdge(i, currentFileNum-hiddenFileNum, deepCnt, deepLevel, parentFileNum, fileName)
			if err := writeToBuffer(buf, row); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if f.IsDir() {
				deep(targetDir+"/"+fileName, isAll, deepLevel, deepCnt, currentFileNum-hiddenFileNum, buf)
			}
		}
	}
}

// ...
func getFilesAndHiddenCnt(targetDir string) ([]os.FileInfo, int, error) {
	files, err := ioutil.ReadDir(targetDir)
	if err != nil {
		return nil, 0, err
	}

	var hiddens, nonHiddens, all []os.FileInfo
	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".") {
			hiddens = append(hiddens, f)
		} else {
			nonHiddens = append(nonHiddens, f)
		}
	}
	all = append(nonHiddens, hiddens...)

	return all, len(hiddens), nil
}

func writeToBuffer(buf *bytes.Buffer, row string) error {
	_, err := buf.WriteString(row + "\n")
	return err
}

func rowWithEdge(i, targetFileNum, deepCnt, deepLevel, parentFileNum int, fileName string) string {
	const (
		four  = "    "
		three = "   "
		edge0 = "│" + three
		edge1 = "├──" + " "
		edge2 = "└──" + " "
	)

	row := ""
	if targetFileNum == 1 && parentFileNum > 1 {
		for i := 1; i < deepCnt; i++ {
			row += edge0
		}
	} else if targetFileNum == 1 || parentFileNum == 1 {
		for i := 1; i < deepCnt; i++ {
			row += four
		}
	} else {
		for i := 1; i < deepCnt; i++ {
			row += edge0
		}
	}

	// ...
	if (i == targetFileNum-1 && deepCnt != deepLevel) || (i == targetFileNum && deepCnt == deepLevel) || (i == targetFileNum-1 && deepCnt == deepLevel) {
		row += edge2 + fileName
	} else {
		row += edge1 + fileName
	}

	return row
}

func output(buf *bytes.Buffer, outFileName string) error {
	if outFileName == "" {
		io.Copy(os.Stdout, buf)
		return nil
	} else {
		f, err := os.Create(outFileName)
		if err != nil {
			return err
		}
		defer f.Close()
		io.Copy(f, buf)
	}

	return nil
}
