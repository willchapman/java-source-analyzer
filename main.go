package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// handleError is a generic error handler
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

// main is the program entry point
func main() {

	start := time.Now().UnixMilli()

	var srcDirectory string
	if len(os.Args) > 1 {
		srcDirectory = os.Args[1]
	} else {
		srcDirectory = "."
	}

	fullPath, _ := filepath.Abs(srcDirectory)
	fmt.Printf("Analyzing %s...\n", fullPath)

	fileCounter, dirCounter, javaFileCounter := 0, 0, 0
	var lines uint32
	var comments uint32
	var allLines uint32
	filepath.Walk("elasticsearch-master", func(path string, info fs.FileInfo, err error) error {
		handleError(err)
		if info.IsDir() {
			dirCounter++
		} else {
			fileCounter++
			if strings.HasSuffix(path, ".java") {
				javaFileCounter++
				thisFileLines, thisCommentLines, fileLines := processJavaFile(path)
				lines += thisFileLines
				comments += thisCommentLines
				allLines += fileLines
			}
		}
		return nil
	})

	fmt.Println("Files\n\t", fileCounter)
	fmt.Println("Java Source Files\n\t", javaFileCounter)
	fmt.Println("Lines\n\t", lines, " ", fmt.Sprintf("( %.2f%% )", (float32(lines)/float32(allLines))*100))
	fmt.Println("Comments\n\t", comments, " ", fmt.Sprintf("( %.2f%% )", (float32(comments)/float32(allLines))*100))
	fmt.Println("All Lines\n\t", allLines)

	duration := float64(time.Now().UnixMilli()-start) / 1000
	fmt.Printf("Scan took %.1f seconds", duration)

}

// processJavaFile is a method that analyzes a source file and counts and returns the
// number of java source code lines, number of comments and the total number of non-empty
// lines.
func processJavaFile(path string) (uint32, uint32, uint32) {

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	var lines uint32
	var commentLines uint32
	var inComment bool
	var allLines uint32

	inComment = false
	for scanner.Scan() {
		txt := scanner.Text()

		if len(txt) > 0 {
			allLines++
		}

		//
		// Don't count empty lines, or the package/import statements... although it is "code" right?
		if len(txt) == 0 {
			continue
		}

		if strings.HasPrefix(txt, "/*") {
			inComment = true
		} else if strings.HasSuffix(txt, "*/") {
			inComment = false
		}

		if strings.HasPrefix(txt, "//") {
			commentLines++
			continue
		}

		if !inComment {
			lines++
		} else {
			commentLines++
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return lines, commentLines, allLines
}