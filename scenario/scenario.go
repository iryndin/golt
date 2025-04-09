package scenario

import "sync"

const TickMillis = 10
const ChanLen = 10000

type ResponseData struct {
	StartTimeUnixMs int64
	EndTimeUnixMs   int64
	StatusCode      int
	ResponseSize    int
}

type RequestFunc func() ResponseData

type Scenario interface {
	GetName() string
	Ramp(totalRequests int, durationString string, requestFunc RequestFunc)
	AtConstantRate(rps float64, durationString string, requestFunc RequestFunc)
	AtOnce(totalRequests int, requestFunc RequestFunc)
	Wait(durationString string)
	Stop()
	GetResults() []ResponseData
}

func NewScenario(name string) Scenario {
	scenario := new(simpleScenario)

	scenario.name = name
	scenario.tickMillis = TickMillis
	scenario.wg = &sync.WaitGroup{}
	scenario.responseDataChan = make(chan ResponseData, ChanLen)

	return scenario
}
