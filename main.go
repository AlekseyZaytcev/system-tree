package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

// LastItemSign char for last file/folder in tree
const LastItemSign = "└───"

// ItemSign char for file/folder if it is at middle of tree
const ItemSign = "├───"

// IndentationSize default indentation for tree level(2 spaces)
const IndentationSize = "\t"

// IndentationWithDelimeter additional indentation for draw tree level
const IndentationWithDelimeter = "│\t"

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"

	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	drawLevel(out, sortedLevel(path, printFiles), path, "", printFiles)
	return nil
}

func drawLevel(output io.Writer, files []os.FileInfo, path string, indentation string, printFiles bool) {
	for id, file := range files {
		if file.IsDir() {
			print(
				folderPresentation(isLastFile(id, files), file, indentation),
				output,
			)

			subPath := filepath.Join(path, file.Name())
			drawLevel(
				output,
				sortedLevel(subPath, printFiles),
				subPath,
				subIndentation(id, files, indentation),
				printFiles,
			)
			continue

		} else if isLastFile(id, files) {
			print(itemPresentation(file, indentation, true), output)
			break

		} else {
			print(itemPresentation(file, indentation, false), output)
		}
	}
}

func folderPresentation(isLast bool, file os.FileInfo, indentation string) string {
	var str string
	if isLast {
		str = itemPresentation(file, indentation, true)
	} else {
		str = itemPresentation(file, indentation, false)
	}

	return str
}

func sortedLevel(path string, printFiles bool) []os.FileInfo {
	levelFiles, _ := ioutil.ReadDir(path)
	if !printFiles {
		levelFiles = filterLevel(levelFiles)
	}

	sort.Slice(levelFiles[:], func(i, j int) bool {
		return levelFiles[i].Name() < levelFiles[j].Name()
	})

	return levelFiles
}

func filterLevel(files []os.FileInfo) []os.FileInfo {
	result := make([]os.FileInfo, 0)
	for _, file := range files {
		if file.IsDir() {
			result = append(result, file)
		}
	}
	return result
}

func isLastFile(id int, files []os.FileInfo) bool {
	return id == countFiles(files)
}

func itemPresentation(file os.FileInfo, indentation string, isLastFile bool) string {
	var formattedSize string
	var formatter string

	if file.Size() == 0 && !file.IsDir() {
		formattedSize = "(empty)"
	} else if !file.IsDir() {
		formattedSize = fmt.Sprintf("(%db)", file.Size())
	}

	if isLastFile {
		formatter = LastItemSign
	} else {
		formatter = ItemSign
	}

	return fmt.Sprintf("%s%s%s %s", indentation, formatter, file.Name(), formattedSize)
}

func subIndentation(id int, files []os.FileInfo, indentation string) string {
	str := ""
	isIdentationEmpty := len(indentation) == 0

	if isIdentationEmpty && !isLastFile(id, files) {
		str = IndentationWithDelimeter

	} else if !isLastFile(id, files) {
		str = indentation + IndentationWithDelimeter

	} else {
		str = indentation + IndentationSize
	}

	return str
}

func countFiles(files []os.FileInfo) int {
	return len(files) - 1
}

func print(str string, output io.Writer) {
	fmt.Fprintln(output, str)
}
