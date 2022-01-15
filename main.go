package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"text/template"
	"time"
)

type StatusData struct {
	Status struct {
		Water int `json:"water"`
		Wind  int `json:"wind"`
	} `json:"status"`
}

func main() {
	go AutoReloadJSON()
	http.HandleFunc("/", AutoReloadWeb)
	http.ListenAndServe(":8080", nil)
}

func AutoReloadJSON() {
	for {
		min := 1
		max := 100
		wind := rand.Intn(max-min) + min
		water := rand.Intn(max-min) + min

		data := StatusData{}
		data.Status.Wind = wind
		data.Status.Water = water

		jsonData, err := json.Marshal(data)

		if err != nil {
			log.Fatal(err.Error())
		}
		err = ioutil.WriteFile("data.json", jsonData, 0644)

		if err != nil {
			log.Fatal(err.Error())
		}
		time.Sleep(15 * time.Second)
	}
}

func AutoReloadWeb(w http.ResponseWriter, r *http.Request) {
	fileData, err := ioutil.ReadFile("data.json")

	if err != nil {
		log.Fatal(err.Error())
	}

	var statusData StatusData

	err = json.Unmarshal(fileData, &statusData)
	if err != nil {
		log.Fatal(err.Error())
	}

	waterVal := statusData.Status.Water
	windVal := statusData.Status.Wind

	var (
		waterStatus string
		windStatus  string
	)
	switch {
	case waterVal <= 5:
		waterStatus = "aman"
	case waterVal >= 6 && waterVal <= 8:
		waterStatus = "siaga"
	default:
		waterStatus = "bahaya"
	}

	switch {
	case windVal <= 6:
		windStatus = "aman"
	case windVal >= 7 && windVal <= 15:
		windStatus = "siaga"
	default:
		windStatus = "bahaya"
	}

	data := map[string]interface{}{
		"waterStatus": waterStatus,
		"windStatus":  windStatus,
		"waterHeight": waterVal,
		"windSpeed":   windVal,
	}

	tpl, err := template.ParseFiles("index.html")

	if err != nil {
		log.Fatal(err.Error())
	}

	tpl.Execute(w, data)

}
