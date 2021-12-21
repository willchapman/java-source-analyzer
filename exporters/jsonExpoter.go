package exporters

import (
	"encoding/json"
	"fmt"
)

type jsonExporter struct{}

func (je jsonExporter) Name() string {
	return "JSON"
}

func (je jsonExporter) DoExport(data DirectoryAnalysisData) error {
	var jsonMap = make(map[string]uint32)

	jsonMap["totalFiles"] = data.FileCounter
	jsonMap["javaFiles"] = data.JavaFileCount
	jsonMap["codeLines"] = data.CodeLinesCount
	jsonMap["commentLines"] = data.CommentLinesCount
	jsonMap["allLines"] = data.AllLinesCount

	dataBytes, err := json.Marshal(jsonMap)
	if err != nil {
		return fmt.Errorf("Failed to encode JSON %#v", err)
	}

	fmt.Println(string(dataBytes))
	return nil

}
