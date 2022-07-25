package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Status struct {
	Disaster DisasterStatus `json:"status"`
}

type DisasterStatus struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/updateStatus", disasterInfo).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8080", router))
	// ticker := time.NewTicker(14 * time.Second)
	// quit := make(chan struct{})
	// go func() {
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			disasterInfo()
	// 		case <-quit:
	// 			ticker.Stop()
	// 			return
	// 		}
	// 	}
	// }()
	// http.Handle("/", http.FileServer(http.Dir("./html")))
	// http.ListenAndServe(":8080", nil)

}

func disasterInfo(w http.ResponseWriter, r *http.Request) { //w http.ResponseWriter, r *http.Request
	var disasterInfoStat Status

	data, _ := ioutil.ReadFile("data.json")
	err := json.Unmarshal(data, &disasterInfoStat)

	min := 1
	max := 20

	fmt.Println(disasterInfoStat)

	if err != nil {
		panic(err)
	}
	rand.Seed(time.Now().UnixNano())

	windStatus := rand.Intn(max-min) + min
	waterStatus := rand.Intn(max-min) + min

	fmt.Println(windStatus)
	fmt.Println(waterStatus)

	disasterInfoStat.Disaster.Water = waterStatus
	disasterInfoStat.Disaster.Wind = windStatus
	jsonFile, _ := json.Marshal(disasterInfoStat)
	_ = ioutil.WriteFile("data.json", jsonFile, 0644)

	type templateInfo struct {
		Wind        int
		StatusWind  string
		Water       int
		StatusWater string
	}

	var info templateInfo

	info.Water = disasterInfoStat.Disaster.Water
	info.Wind = disasterInfoStat.Disaster.Wind
	info.StatusWater = waterCondition(info.Water)
	info.StatusWind = windCondition(info.Wind)
	tpl, err := template.ParseFiles("index.html")

	if err != nil {
		panic(err)
	}
	err = tpl.Execute(w, info)

}

func waterCondition(speed int) string {
	if speed <= 5 {
		return "aman"
	} else if 6 <= speed && speed <= 8 {
		return "siaga"
	} else {
		return "bahaya"
	}

}

func windCondition(speed int) string {
	if speed <= 6 {
		return "aman"
	} else if 7 <= speed && speed <= 15 {
		return "siaga"
	} else {
		return "bahaya"
	}

}
