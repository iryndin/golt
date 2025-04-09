package datawrite

import (
	"bufio"
	"fmt"
	"iryndin/golt/scenario"
	"os"
)

func WriteResults(filepath string, data []scenario.ResponseData) (string, error) {
	file, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString("StartTimeMs,EndTimeMs,ElapsedTimeMs,ResponseStatus Code,ResponseSize\n")
	if err != nil {
		return "", err
	}
	for _, d := range data {
		line := fmt.Sprintf("%d,%d,%d,%d,%d\n", d.StartTimeUnixMs, d.EndTimeUnixMs,
			d.EndTimeUnixMs-d.StartTimeUnixMs, d.StatusCode, d.ResponseSize)
		_, err := writer.WriteString(line)
		if err != nil {
			return "", err
		}
	}
	return filepath, writer.Flush()
}
