package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/cors"
	"math"
	"math/rand"
	"net/http"
	"time"
)

type WaitTimes struct {
	Id        uuid.UUID `json:"id"`
	DateTime  time.Time `json:"datetime"`
	Liftname  string    `json:"liftname"`
	Waittimes float64   `json:"waittime"`
}

func generateData(t time.Time) float64 {
	minutes := float64(t.Sub(time.Date(t.Year(), t.Month(), t.Day(), 8, 0, 0, 0, t.Location())).Minutes())

	maxWaitTime := rand.Float64()*30 + 10
	peakTime := 240.0
	width := 240.0
	var value float64
	if minutes < peakTime {
		value = maxWaitTime * (minutes / width)
	} else {
		value = maxWaitTime * ((480 - minutes) / width)
	}

	noise := rand.Float64()*5 - 2.5
	result := value + noise
	return math.Max(0, result)
}

func main() {
	liftnames := []string{
		"Alpine Ascender",
		"Frosty Flyer",
		"Summit Seeker",
		"Blizzard Breeze",
		"Glacier Glide",
		"Peak Pioneer",
		"Snowdrift Soarer",
		"Everest Express",
		"Winter Whirl",
		"Crystal Cruiser",
	}

	var waitTimes []WaitTimes

	for _, l := range liftnames {
		start := time.Date(2023, 12, 24, 8, 0, 0, 0, time.Local)
		end := time.Date(2023, 12, 24, 17, 0, 0, 0, time.Local)

		for t := start; t.Before(end) || t.Equal(end); t = t.Add(5 * time.Minute) {
			w := WaitTimes{Waittimes: generateData(t), Id: uuid.New(), DateTime: t, Liftname: l}
			waitTimes = append(waitTimes, w)
		}
	}

	fmt.Println("Server started on Port 4242")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /waittimes/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		wjson, err := json.Marshal(waitTimes)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write(wjson)
	})
	handler := cors.Default().Handler(mux)
	http.ListenAndServe("localhost:4242", handler)
}
