package main

import (
	"math/rand/v2"
	"time"

	"github.com/HappyFatPug/go-course-metrics/internal/agent"
)

var PollCount int64 = 0
var RandomValue float64

func main() {
	app := agent.ConfigureTestApp(2*time.Second, 5*time.Second, true)

	pollTicker := time.NewTicker(app.PollInterval)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(app.ReportInterval)
	defer reportTicker.Stop()

	go func() {
		for range pollTicker.C {
			PollCount += 1
			RandomValue = rand.Float64()
			app.Service.UpdateAll(PollCount, RandomValue)
		}
	}()

	// go func() {
	// 	for range reportTicker.C {
	// 		fmt.Println("REPORT!")
	// 	}
	// }()

	select {}
}
