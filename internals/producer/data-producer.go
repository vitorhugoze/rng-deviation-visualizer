package producer

import (
	"log"
	ws "main/internals/websocket"
	"main/pkg/rng"
	"time"
)

type DeviationData struct {
	Index     int     `json:"index"`
	Deviation float64 `json:"deviation"`
}

/*
Generates a list of random numbers, then calculates the standard deviation
whithin that group and add that deviation to the WebSocket queue to be transmitted
*/
func ProduceDeviationData(interval time.Duration, totalRandoms int, shards int) {

	index := 0

	for {
		time.Sleep(interval)

		generator, err := rng.NewRandomGenerator(totalRandoms, shards)
		if err != nil {
			log.Fatal(err)
		}

		group := generator.GenerateRandomGroup()
		std := generator.CalculatesStandardDeviation(group)

		deviation := DeviationData{
			Index:     index,
			Deviation: std,
		}

		ws.AddMessageToQueue(deviation)
		index++
	}

}
