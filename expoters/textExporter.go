package expoters

import "fmt"

type textExporter struct{}

func (je textExporter) Name() string {
	return "TEXT"
}

func (je textExporter) DoExport(data DirectoryAnalysisData) error {
	fmt.Printf("Files %d\n", data.FileCounter)
	fmt.Printf("Java Files %d\n", data.JavaFileCount)
	fmt.Printf("Code Lines %d\n", data.CodeLinesCount)
	fmt.Printf("Comment Lines %d\n", data.CommentLinesCount)
	fmt.Printf("All Lines %d\n", data.AllLinesCount)
	return nil
}
