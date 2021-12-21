package expoters

import "fmt"

type ExportDestinationType uint8

const (
	TEXT ExportDestinationType = iota
)

// DirectoryAnalysisData is a struct that holds the information for the entire directory
type DirectoryAnalysisData struct {
	FileCounter       uint32
	JavaFileCount     uint32
	CodeLinesCount    uint32
	CommentLinesCount uint32
	AllLinesCount     uint32
}

// FileAnalysisData is a struct to hold information about a single file
type FileAnalysisData struct {
	CodeLinesCount    uint32
	CommentLinesCount uint32
	AllLinesCount     uint32
}

type Exporter interface {
	Name() string
	DoExport(data DirectoryAnalysisData) error
}

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

func MakeExporter(exportType ExportDestinationType) (Exporter, error) {
	switch exportType {

	case TEXT:
		return textExporter{}, nil

	}
	return nil, fmt.Errorf("invalid export type %d", exportType)
}
