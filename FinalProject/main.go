package main

import (
	"encoding/json"
	"finalproject/handler"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

func main() {

	fmt.Println("START")

	mux := http.NewServeMux()

	finalHandler := http.HandlerFunc(final)
	mux.Handle("/", middlewareOne(middlewareTwo(finalHandler)))
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}

func middlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		if r.URL.Path == "/users/register" || r.URL.Path == "/users/login" {
			next.ServeHTTP(w, r)
		} else {

			//Get input Token
			token := r.Header.Get("Authorization")
			splitToken := strings.Split(token, "Bearer ")
			reqToken := splitToken[1]
			fmt.Println("Token : ", reqToken)

			var jsonFile struct {
				Token string
			}
			jsonFile.Token = reqToken
			//write token to config
			file, _ := json.MarshalIndent(jsonFile, "", " ")
			_ = ioutil.WriteFile("config.json", file, 0644)
			next.ServeHTTP(w, r)
		}
		//log.Print("Executing middlewareOne again")
	})
}

func middlewareTwo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//log.Print("Executing middlewareTwo")

		next.ServeHTTP(w, r)
		//log.Print("Executing middlewareTwo again")
	})
}

func final(w http.ResponseWriter, r *http.Request) {
	log.Print("Executing finalHandler")
	handler.MainHandler(w, r)
}
