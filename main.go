package main

import (
	"iryndin/golt/datafeed"
	"iryndin/golt/dataviz"
	"iryndin/golt/datawrite"
	"iryndin/golt/scenario"

	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

var queries []string

func myRequestFunc() scenario.ResponseData {
	query := queries[rand.Intn(len(queries))]
	requestUrl := fmt.Sprintf("http://127.0.0.1:8001/search?queries=%s", query)

	startTimeUnixMs := time.Now().UnixMilli()
	fmt.Printf("  sent: %d, url: %s\n", startTimeUnixMs, requestUrl)
	resp, _ := http.Get(requestUrl)
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	return scenario.ResponseData{
		StartTimeUnixMs: startTimeUnixMs,
		EndTimeUnixMs:   time.Now().UnixMilli(),
		StatusCode:      resp.StatusCode,
		ResponseSize:    len(bodyBytes),
	}
}

func main() {
	runTest2()
}

func readTestResultsAndWriteSimulation() {
	d := datawrite.ReadResults("results.csv")
	dataviz.WriteSimulationResults(d)
}

func runTest2() {
	queries = datafeed.LinesFromTextFile("queries.txt")
	sc := scenario.NewScenario("test2")
	sc.Ramp(5, "10s", myRequestFunc)
	sc.AtOnce(20, myRequestFunc)
	sc.AtConstantRate(5.0, "30s", myRequestFunc)
	sc.Stop()

	dataviz.WriteSimulationResults(sc.GetResults())

	fmt.Println("====================")
	fmt.Println("DONE")
}

func runTest1() {
	queries = datafeed.LinesFromTextFile("queries.txt")
	sc := scenario.NewScenario("test1")
	sc.Ramp(5, "10s", myRequestFunc)
	sc.AtOnce(10, myRequestFunc)
	sc.AtConstantRate(3.0, "20s", myRequestFunc)
	sc.Stop()

	dataviz.WriteSimulationResults(sc.GetResults())

	//datawrite.WriteResults("results.csv", data)
	fmt.Println("====================")
	fmt.Println("DONE")
}
