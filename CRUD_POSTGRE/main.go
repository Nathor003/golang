package main

import (
	"UserService/entity"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var allData = make(map[int]entity.User)

var (
	db  *sql.DB
	err error
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "db-go-sql"
)

func connectDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success connect DB")

	return db
}

func main() {

	fmt.Println("TEST")

	router := mux.NewRouter()

	router.HandleFunc("/users", getData).Methods(http.MethodGet)
	router.HandleFunc("/users", register).Methods(http.MethodPost)
	router.HandleFunc("/users/{id}", getDataById).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", deleteData).Methods(http.MethodDelete)
	router.HandleFunc("/users/{id}", putData).Methods(http.MethodPut)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getData(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Get Data")

	conn := connectDB()
	var allData []entity.User
	sqlSelect := `SELECT id,username,password,email,age from "users"`

	rows, err := conn.Query(sqlSelect)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var getData entity.User

		err = rows.Scan(&getData.Id, &getData.Username, &getData.Password, &getData.Email, &getData.Age)
		if err != nil {
			panic(err)
		}

		allData = append(allData, getData)
	}

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

	conn := connectDB()
	var allData entity.User
	sqlSelect := `SELECT id,username,password,email,age from "users" where id=$1`

	rows, err := conn.Query(sqlSelect, idData)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var getData entity.User

		err = rows.Scan(&getData.Id, &getData.Username, &getData.Password, &getData.Email, &getData.Age)
		if err != nil {
			panic(err)
		}

		allData = getData
	}

	u, _ := json.Marshal(&allData)
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

	conn := connectDB()

	sqlInsert := `
	insert into "users"(username,password,email,age)
	values($1,$2,$3,$4)
	`

	_, err = conn.Exec(sqlInsert, user.Username, user.Password, user.Email, user.Age)

	if err != nil {
		panic(err)
	}

	fmt.Println("Sucess insert data")

	conn.Close()
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

	conn := connectDB()

	sqlUpdate := `
	UPDATE "users" set username=$2, password=$3, email=$4, age=$5
	where id=$1
	`
	_, err := conn.Exec(sqlUpdate, id, user.Username, user.Password, user.Email, user.Age)

	if err != nil {
		panic(err)
	}
	fmt.Println("Success Update data")

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
	conn := connectDB()
	sqlSelect := `DELETE from "users" where id=$1`

	_, err := conn.Query(sqlSelect, idData)
	if err != nil {
		panic(err)
	}

}
