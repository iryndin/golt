package datafeed

import (
	"bufio"
	"log"
	"os"
)

func LinesFromTextFile(filepath string) []string {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	result := make([]string, 0, 100)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, line)
	}

	return result
}
