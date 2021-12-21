package main

import (
	"bufio"
	flag "flag"
	"fmt"
	"io/fs"
	"java-source-analyzer/expoters"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// handleError is a generic error handler
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

// main is the program entry point
func main() {
	var outputType string
	flag.StringVar(&outputType, "output", "TEXT", "Chooses the output format (JSON,TEXT)")
	flag.Parse()

	var exportType expoters.ExportDestinationType
	switch outputType {

	case "JSON":
		exportType = expoters.JSON
		break

	default:
		exportType = expoters.TEXT
		break

	}

	exporter, err := expoters.MakeExporter(exportType)
	handleError(err)

	//start := time.Now().UnixMilli()

	var srcDirectory string
	if len(flag.Args()) == 0 {
		fmt.Println("Usage: " + os.Args[0] + ": <options> source-directory")
		fmt.Println(" Options are listed below")
		flag.PrintDefaults()
		return
	} else {
		srcDirectory = flag.Arg(0)
	}

	//fullPath, _ := filepath.Abs(srcDirectory)
	//fmt.Printf("Analyzing %s...\n", fullPath)

	var dirData expoters.DirectoryAnalysisData

	err = filepath.Walk(srcDirectory, func(path string, info fs.FileInfo, err error) error {
		handleError(err)
		if !info.IsDir() {
			dirData.FileCounter++
			if strings.HasSuffix(path, ".java") {
				dirData.JavaFileCount++
				fileData := processJavaFile(path)
				dirData.CodeLinesCount += fileData.CodeLinesCount
				dirData.CommentLinesCount += fileData.CommentLinesCount
				dirData.AllLinesCount += fileData.AllLinesCount
			}
		}
		return nil
	})
	handleError(err)

	err = exporter.DoExport(dirData)
	handleError(err)

	//duration := float64(time.Now().UnixMilli()-start) / 1000
	//fmt.Printf("\nScan took %.1f seconds", duration)

}

// processJavaFile is a method that analyzes a source file and counts and returns the
// number of java source code lines, number of comments and the total number of non-empty
// lines.
func processJavaFile(path string) expoters.FileAnalysisData {

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

	return expoters.FileAnalysisData{CodeLinesCount: lines, CommentLinesCount: commentLines, AllLinesCount: allLines}
}
