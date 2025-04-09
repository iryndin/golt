package datawrite

import (
	"bufio"
	"iryndin/golt/scenario"
	"strconv"
	"strings"

	//"fmt"
	"log"
	"os"
)

func ReadResults(filepath string) []scenario.ResponseData {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	result := make([]scenario.ResponseData, 0, 100)
	scanner := bufio.NewScanner(file)
	lineCounter := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineCounter++
		if lineCounter == 1 {
			// we skip first line - it is a header
			continue
		}
		rd := buildResponseData(line)
		result = append(result, rd)
	}
	return result
}

func buildResponseData(line string) scenario.ResponseData {
	a := strings.Split(line, ",")

	startMs, _ := strconv.Atoi(a[0])
	endMs, _ := strconv.Atoi(a[1])
	status, _ := strconv.Atoi(a[3])
	respSize, _ := strconv.Atoi(a[4])

	return scenario.ResponseData{
		StartTimeUnixMs: int64(startMs),
		EndTimeUnixMs:   int64(endMs),
		StatusCode:      status,
		ResponseSize:    respSize,
	}
}
