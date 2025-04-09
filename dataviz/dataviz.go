package dataviz

import (
	"bufio"
	_ "embed"
	"fmt"
	"iryndin/golt/datawrite"
	"iryndin/golt/scenario"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const Line1 = "const LoadTestData = ["
const FormatRecordLine = "  { startMs: %d, endMs: %d, elapsedMs: %d, size: %d },\n"
const LineEnd = "];"

//go:embed static/js/response_time_distribution.js
var responseTimeDistributionJsContent string

//go:embed static/css/styles.css
var stylesCssContent string

//go:embed static/index.html
var indexHtmlContent string

func WriteSimulationResults(data []scenario.ResponseData) (string, error) {
	// 1. Create simulation folder
	simulationFolderName := createFolder()

	// 2. Write results.csv file
	resultsFilePath := filepath.Join(simulationFolderName, "results.csv")
	datawrite.WriteResults(resultsFilePath, data)

	// 3. Write data.js file
	jsFolder := filepath.Join(simulationFolderName, "js")
	os.MkdirAll(jsFolder, 0755)
	dataJsFilepath := filepath.Join(jsFolder, "data.js")
	WriteDataJs(dataJsFilepath, data)

	// 4. Write response_time_distribution.js
	responseTimeDistributionJsFilepath := filepath.Join(jsFolder, "response_time_distribution.js")
	writeFileStringContent(responseTimeDistributionJsFilepath, responseTimeDistributionJsContent)

	// 5. Write styles.css
	cssFolder := filepath.Join(simulationFolderName, "css")
	os.MkdirAll(cssFolder, 0755)
	stylesCssFilepath := filepath.Join(cssFolder, "styles.css")
	writeFileStringContent(stylesCssFilepath, stylesCssContent)

	// 6. Write index.html into simulation folder root
	indexHtmlFilepath := filepath.Join(simulationFolderName, "index.html")
	writeFileStringContent(indexHtmlFilepath, indexHtmlContent)

	return indexHtmlFilepath, nil
}

func createFolder() string {
	dtSuffix := time.Now().Format("20060102T150405")
	folderName := fmt.Sprintf("simulation-" + dtSuffix)
	err := os.MkdirAll(folderName, 0755)
	if err != nil {
		log.Fatal(err)
	}
	return folderName
}

func WriteDataJs(filepath string, data []scenario.ResponseData) {
	file, _ := os.Create(filepath)
	defer file.Close()

	writer := bufio.NewWriter(file)

	elps := getElapsedOnly(data)
	astats := CalculateAggregateStats(elps)

	simulationStartTimeUnixMs, _ := findMinMaxForI64(data, func(x scenario.ResponseData) int64 {
		return x.StartTimeUnixMs
	})
	_, simulationEndTimeUnixMs := findMinMaxForI64(data, func(x scenario.ResponseData) int64 {
		return x.EndTimeUnixMs
	})

	astats.simulationStartTime = time.UnixMilli(simulationStartTimeUnixMs).Format("2006-01-02 15:04:05")
	astats.simulationEndTime = time.UnixMilli(simulationEndTimeUnixMs).Format("2006-01-02 15:04:05")
	astats.simulationDurationSeconds = (simulationEndTimeUnixMs - simulationStartTimeUnixMs) / 1000

	writeData(writer, data)
	writer.WriteString("\n\n")
	writeStats(writer, astats)
}

func getElapsedOnly(data []scenario.ResponseData) []int {
	res := make([]int, 0, len(data))
	for _, val := range data {
		res = append(res, int(val.EndTimeUnixMs-val.StartTimeUnixMs))
	}
	return res
}

func writeData(writer *bufio.Writer, data []scenario.ResponseData) {
	writer.WriteString(Line1)

	for _, d := range data {
		line := fmt.Sprintf(FormatRecordLine, d.StartTimeUnixMs, d.EndTimeUnixMs,
			d.EndTimeUnixMs-d.StartTimeUnixMs, d.ResponseSize)
		writer.WriteString(line)
	}
	writer.WriteString(LineEnd)
	writer.Flush()
}

func writeStats(writer *bufio.Writer, astats AggregateStats) {
	writer.WriteString("const AggregateStats = {\n")
	writer.WriteString("  startTimestamp: \"" + astats.simulationStartTime + "\",\n")
	writer.WriteString("  endTimestamp: \"" + astats.simulationEndTime + "\",\n")
	writer.WriteString("  durationSeconds: " + strconv.FormatInt(astats.simulationDurationSeconds, 10) + ",\n")
	writer.WriteString("  total: " + strconv.Itoa(astats.total) + ",\n")
	writer.WriteString("  min: " + strconv.Itoa(astats.min) + ",\n")
	writer.WriteString("  max: " + strconv.Itoa(astats.max) + ",\n")
	writer.WriteString("  mean: " + strconv.Itoa(astats.mean) + ",\n")
	writer.WriteString("  stddev: " + strconv.Itoa(astats.stdDev) + ",\n")
	writer.WriteString("  p50: " + strconv.Itoa(astats.p50) + ",\n")
	writer.WriteString("  p75: " + strconv.Itoa(astats.p75) + ",\n")
	writer.WriteString("  p90: " + strconv.Itoa(astats.p90) + ",\n")
	writer.WriteString("  p95: " + strconv.Itoa(astats.p95) + ",\n")
	writer.WriteString("  p99: " + strconv.Itoa(astats.p99) + ",\n")
	writer.WriteString("};\n")
	writer.Flush()
}

func writeFileStringContent(filepath string, content string) {
	file, _ := os.Create(filepath)
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(content)
	writer.Flush()
}
