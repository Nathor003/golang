package main

import (
	"UserService/entity"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var allData = make(map[int]entity.User)

func main() {

	fmt.Print("TEST")

	router := mux.NewRouter()

	router.HandleFunc("/users", getData).Methods(http.MethodGet)
	router.HandleFunc("/users", register).Methods(http.MethodPost)
	router.HandleFunc("/users/{id}", getDataById).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", deleteData).Methods(http.MethodDelete)
	router.HandleFunc("/users/{id}", putData).Methods(http.MethodPut)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getData(w http.ResponseWriter, r *http.Request) {

	u, _ := json.Marshal(&allData)
	w.Header().Add("Content-Type", "application/json")
	w.Write(u)
	fmt.Fprint(w)
}

func getDataById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := vars["id"]
	idData, _ := strconv.Atoi(id)

	fmt.Println("id data", idData)

	var data entity.User
	for in, _ := range allData {
		if allData[in].Id == idData {
			data = allData[in]
		}

	}

	fmt.Println(data)

	u, _ := json.Marshal(data)
	w.Header().Add("Content-Type", "application/json")
	w.Write(u)
	fmt.Fprint(w)
}

func register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user entity.User
	if err := decoder.Decode(&user); err != nil {
		w.Write([]byte("error decoding json body"))
		return
	}
	idx := len(allData)
	allData[idx] = user
	// u, _ := json.Marshal(&user)
	// w.Header().Add("Content-Type", "application/json")
	// w.Write(u)
	// fmt.Fprint(w)
}

func putData(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var user entity.User
	if err := decoder.Decode(&user); err != nil {
		w.Write([]byte("error decoding json body"))
		return
	}

	vars := mux.Vars(r)
	id, _ := vars["id"]
	idData, _ := strconv.Atoi(id)
	idx := len(allData)
	allData[idx] = user

	for in, _ := range allData {
		if allData[in].Id == idData {
			allData[in] = user
		}

	}
	// u, _ := json.Marshal(&user)
	// w.Header().Add("Content-Type", "application/json")
	// w.Write(u)
	// fmt.Fprint(w)
}

func deleteData(w http.ResponseWriter, r *http.Request) {

	fmt.Println("DELETE DATA")

	vars := mux.Vars(r)
	id, _ := vars["id"]
	idData, _ := strconv.Atoi(id)

	for idx, _ := range allData {

		if allData[idx].Id == idData {
			delete(allData, idx)
		}

	}
}
