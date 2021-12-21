package expoters

import "fmt"

type ExportDestinationType uint8

const (
	TEXT ExportDestinationType = iota
	JSON
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

func MakeExporter(exportType ExportDestinationType) (Exporter, error) {
	switch exportType {

	case TEXT:
		return textExporter{}, nil

	case JSON:
		return jsonExporter{}, nil

	}
	return nil, fmt.Errorf("invalid export type %d", exportType)
}
