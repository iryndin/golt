package scenario

import (
	"log"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type simpleScenario struct {
	name             string
	tickMillis       int
	wg               *sync.WaitGroup
	responseDataChan chan ResponseData
	results          []ResponseData
}

func (s *simpleScenario) GetName() string {
	return s.name
}

func (s *simpleScenario) Ramp(totalRequests int, durationString string, requestFunc RequestFunc) {
	if totalRequests < 0 {
		panic("totalRequests must be positive: " + strconv.Itoa(totalRequests))
	}

	duration := parseDurationString(durationString)

	if totalRequests == 0 {
		time.Sleep(duration)
		return
	}

	var requestsSent atomic.Uint64
	requestRatePerTick := calcRequestsPerTickForRamp(s.tickMillis, totalRequests, duration)

	startTime := time.Now()
	log.Printf("Start ramp (totalRequests: %d, duration: %s), unixTimeMillis: %d\n",
		totalRequests, durationString, startTime.UnixMilli())

	for time.Since(startTime) < duration {
		time.Sleep(time.Duration(s.tickMillis) * time.Millisecond)
		ticksPassed := time.Since(startTime).Milliseconds() / int64(s.tickMillis)
		currentRatePerTick := float64(requestsSent.Load()) / float64(ticksPassed)
		if currentRatePerTick < requestRatePerTick && requestsSent.Load() < uint64(totalRequests) {
			go func() {
				s.wg.Add(1)
				defer s.wg.Done()
				requestsSent.Add(1)
				s.responseDataChan <- requestFunc()
			}()
		}
	}

	for requestsSent.Load() < uint64(totalRequests) {
		go func() {
			s.wg.Add(1)
			defer s.wg.Done()
			requestsSent.Add(1)
			s.responseDataChan <- requestFunc()
		}()
	}
}

func (s *simpleScenario) AtOnce(totalRequests int, requestFunc RequestFunc) {
	if totalRequests < 0 {
		panic("totalRequests must be positive: " + strconv.Itoa(totalRequests))
	}

	if totalRequests == 0 {
		return
	}

	requestsSent := 0

	startTime := time.Now()
	log.Printf("Start at_once (totalRequests: %d), unixTimeMillis: %d\n",
		totalRequests, startTime.UnixMilli())

	for requestsSent < totalRequests {
		go func() {
			s.wg.Add(1)
			defer s.wg.Done()
			s.responseDataChan <- requestFunc()
		}()
		requestsSent++
	}
}

func (s *simpleScenario) AtConstantRate(rps float64, durationString string, requestFunc RequestFunc) {
	if rps < 0 {
		panic("rps must be positive")
	}

	duration := parseDurationString(durationString)

	if rps == 0 {
		time.Sleep(duration)
		return
	}

	var requestsSent atomic.Uint64
	requestRatePerTick := calcRequestsPerTickForConstantRps(s.tickMillis, rps)

	startTime := time.Now()
	log.Printf("Start constant_rate (rps: %f, duration: %s), unixTimeMillis: %d\n",
		rps, durationString, startTime.UnixMilli())

	for time.Since(startTime) < duration {
		time.Sleep(time.Duration(s.tickMillis) * time.Millisecond)
		ticksPassed := time.Since(startTime).Milliseconds() / int64(s.tickMillis)
		currentRatePerTick := float64(requestsSent.Load()) / float64(ticksPassed)
		if currentRatePerTick < requestRatePerTick {
			go func() {
				s.wg.Add(1)
				defer s.wg.Done()
				requestsSent.Add(1)
				s.responseDataChan <- requestFunc()
			}()
		}
	}
}

func (s *simpleScenario) Wait(durationString string) {
	duration := parseDurationString(durationString)
	time.Sleep(duration)
}

func (s *simpleScenario) Stop() {
	startTime := time.Now()
	log.Printf("Entered stop, unixTimeMillis: %d\n", startTime.UnixMilli())

	s.wg.Wait()
	close(s.responseDataChan)

	s.results = copyFromChanToSlice(s.responseDataChan)
}

func (s *simpleScenario) GetResults() []ResponseData {
	return s.results
}

func copyFromChanToSlice(responseDataChan chan ResponseData) []ResponseData {
	res := make([]ResponseData, 0, ChanLen)
	for d := range responseDataChan {
		res = append(res, d)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].StartTimeUnixMs < res[j].StartTimeUnixMs
	})
	return res
}

func parseDurationString(durationString string) time.Duration {
	d, err := time.ParseDuration(durationString)
	if err != nil {
		panic("Incorrect duration string: " + durationString)
	}
	return d
}

func calcRequestsPerTickForRamp(tickMillis int, totalRequests int, duration time.Duration) float64 {
	rps := float64(totalRequests) / duration.Seconds()
	return calcRequestsPerTickForConstantRps(tickMillis, rps)
}

func calcRequestsPerTickForConstantRps(tickMillis int, rps float64) float64 {
	return (float64(tickMillis) / 1000.0) * rps
}
