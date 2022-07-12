package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Subscription struct {
	Plan          string `json:"plan"`
	Status        string `json:"status"`
	PaymentMethod string `json:"payment_method"`
	Term          string `json:"term"`
}

type CreditCard struct {
	Cc_number string `json:"cc_number"`
}

type Coordinate struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Address struct {
	City          string     `json:"city"`
	StreetName    string     `json:"street_name"`
	StreetAddress string     `json:"street_address"`
	Zipcode       string     `json:"zip_code"`
	State         string     `json:"state"`
	Countery      string     `json:"country"`
	Coordinate    Coordinate `json:"coordinates"`
}

type Employment struct {
	Title    string `json:"title"`
	KeySkill string `json:"key_skill"`
}

type Users struct {
	Id                    int          `json:"id"`
	Uid                   string       `json:"uid"`
	Password              string       `json:"password"`
	FirstName             string       `json:"first_name"`
	LastName              string       `json:"last_name"`
	Username              string       `json:"username"`
	Email                 string       `json:"email"`
	Avatar                string       `json:"avatar"`
	Gender                string       `json:"gender"`
	PhoneNumber           string       `json:"phone_number"`
	SocialInsuranceNumber string       `json:"social_insurance_number"`
	DateOfBirth           string       `json:"date_of_birth"`
	Employment            Employment   `json:"employment"`
	Address               Address      `json:"address"`
	CreditCard            CreditCard   `json:"credit_card"`
	Subscription          Subscription `json:"subscription"`
}

type User struct {
	Id         int     `json:"id"`
	Uid        string  `json:"uid"`
	First_name string  `json:"first_name"`
	Last_name  string  `json:"last_name"`
	Username   string  `json:"username"`
	Address    Address `json:"address"`
}

func main() {
	mux := http.NewServeMux()

	endpoint := http.HandlerFunc(getDataUser)

	mux.Handle("/users", checkAuthentication(goToEndpoint(endpoint)))

	err := http.ListenAndServe(":8080", mux)

	log.Fatal(err)
}

func getDataUser(w http.ResponseWriter, r *http.Request) {
	var allData []User
	response, err := http.Get("https://random-data-api.com/api/users/random_user?size=10")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result []Users
	json.Unmarshal([]byte(responseData), &result)

	fmt.Println(result[0].Address.Coordinate)
	fmt.Println(result[0].Subscription)

	for idx, _ := range result {
		var temp User
		temp.Id = result[idx].Id
		temp.Uid = result[idx].Uid
		temp.First_name = result[idx].FirstName
		temp.Last_name = result[idx].LastName
		temp.Username = result[idx].Username
		temp.Address = result[idx].Address

		allData = append(allData, temp)
	}

	dataJSON, err := json.Marshal(allData)
	if err != nil {
		panic(err)
	}

	w.Write(dataJSON)
}

func checkAuthentication(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		username, password, ok := r.BasicAuth()
		if !ok {
			w.Write([]byte(`Something wrong reading ok`))
			return
		}

		isValid := (username == "username") && (password == "PASSWORD")
		if !isValid {
			w.Write([]byte(`Username or Password incorrect`))
			return
		} else {
			next.ServeHTTP(w, r)

		}
	})
}

func goToEndpoint(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
