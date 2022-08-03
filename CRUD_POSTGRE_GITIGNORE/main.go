package main

import (
	"CRUD_POSTGRE_GITIGNORE/entity"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var allData = make(map[int]entity.User)

var (
	db  *sql.DB
	err error
)

type ConfigDatabase struct {
	Port     int    `yaml:"port" env:"PORT"`
	Host     string `yaml:"host" env:"HOST"`
	Name     string `yaml:"name" env:"NAME"`
	User     string `yaml:"user" env:"USER"`
	Password string `yaml:"password" env:"PASSWORD"`
	DbName   string `yaml:"db"`
}

type item struct{}

type orders struct{}

func connectDB() *sql.DB {

	var cfg ConfigDatabase

	err := cleanenv.ReadConfig("config.yaml", &cfg)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Isi File Config", cfg)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName)

	fmt.Println("Coba konek ke DB")

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success connect DB")

	return db
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var conn *sql.DB

var orderID int

func main() {

	conn = connectDB()
	defer conn.Close()

	fmt.Println("TEST")

	router := mux.NewRouter()

	router.HandleFunc("/users", getData).Methods(http.MethodGet)

	router.HandleFunc("/users", register).Methods(http.MethodPost)
	router.HandleFunc("/users/login", userLogin).Methods(http.MethodGet)

	router.HandleFunc("/users/{id}", getDataById).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", deleteData).Methods(http.MethodDelete)
	router.HandleFunc("/users/{id}", putData).Methods(http.MethodPut)

	router.HandleFunc("/orders", getRegisterOrder).Methods(http.MethodGet)
	router.HandleFunc("/orders", registerOrder).Methods(http.MethodPost)
	router.HandleFunc("/orders/{orderid}", updateRegisterOrder).Methods(http.MethodPut)
	router.HandleFunc("/orders/{orderid}", deleteRegisterOrder).Methods(http.MethodDelete)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func userLogin(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if ok == false {
		panic("err")
	}

	sqlGetPassword := `select password from users where username=$1`

	rows, err := conn.Query(sqlGetPassword, username)

	var savedPass string

	for rows.Next() {
		if err != nil {
			panic(err)

		}

		err = rows.Scan(&savedPass)

		if err != nil {
			panic(err)
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(savedPass), []byte(password))
	if err != nil {
		panic(err)
	}

	fmt.Println("Password match")
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secretKey"))
	if err != nil {
		panic(err)
	}

	fmt.Println("JWT TOKEN : ", tokenString)

	u, _ := json.Marshal(tokenString)
	w.Header().Add("Content-Type", "application/json")
	w.Write(u)
	fmt.Fprint(w)
}

func getData(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Get Data")

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

	pass := []byte(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	fmt.Println(string(hashedPassword))

	sqlInsert := `
	insert into "users"(username,password,email,age)
	values($1,$2,$3,$4)
	`

	_, err = conn.Exec(sqlInsert, user.Username, hashedPassword, user.Email, user.Age)

	if err != nil {
		panic(err)
	}

	fmt.Println("Sucess insert data")

	u, _ := json.Marshal(&user)
	w.Header().Add("Content-Type", "application/json")
	w.Write(u)
	fmt.Fprint(w)
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
	sqlSelect := `DELETE from "users" where id=$1`
	_, err := conn.Query(sqlSelect, idData)
	if err != nil {
		panic(err)
	}

}

func registerOrder(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var order1 entity.Order
	if err := decoder.Decode(&order1); err != nil {
		w.Write([]byte("error decoding json body"))
		fmt.Println(err)
		return
	}

	fmt.Println(len(order1.Barang))

	orderID++

	sqlInsertOrder := `Insert into "orders"(order_id,customer_name,ordered_at) values($1,$2,$3) `

	date, _ := time.Parse("2006-01-02 15:04", order1.OrderAt)
	fmt.Println(date)
	_, err := conn.Query(sqlInsertOrder, orderID, order1.CustomerName, date)
	if err != nil {
		panic(err)
	}
	sqlInsertItem := `INSERT INTO "item"(item_id,item_code,description,quantity,order_id) values ($1,$2,$3,$4,$5)`

	for idx, _ := range order1.Barang {
		_, err := conn.Query(sqlInsertItem,
			order1.Barang[idx].LineItemID,
			order1.Barang[idx].ItemCode,
			order1.Barang[idx].Description,
			order1.Barang[idx].Quantity,
			orderID,
		)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Success Insert")
}

func updateRegisterOrder(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, _ := vars["orderid"]
	idData, _ := strconv.Atoi(id)

	decoder := json.NewDecoder(r.Body)
	var order1 entity.Order
	if err := decoder.Decode(&order1); err != nil {
		w.Write([]byte("error decoding json body"))
		fmt.Println(err)
		return
	}

	fmt.Println(len(order1.Barang))

	orderID++

	sqlUpdateOrder := `UPDATE "orders" SET customer_name=$1,ordered_at=$2 where order_id=$3`

	date, _ := time.Parse("2006-01-02 15:04", order1.OrderAt)
	fmt.Println(date)
	_, err := conn.Query(sqlUpdateOrder, order1.CustomerName, date, idData)
	if err != nil {
		panic(err)
	}
	sqlUpdateItem := `UPDATE "item" set item_id=$1,item_code=$2,description=$3,quantity=$4 where order_id=$5 and item_code=$6`

	for idx, _ := range order1.Barang {
		_, err := conn.Query(sqlUpdateItem,
			order1.Barang[idx].LineItemID,
			order1.Barang[idx].ItemCode,
			order1.Barang[idx].Description,
			order1.Barang[idx].Quantity,
			idData,
			order1.Barang[idx].ItemCode,
		)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Success UPDATE")
}

func getRegisterOrder(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Get Data")

	var allData []entity.Order
	sqlSelect := `select 
    a.customer_name,
    a.ordered_at,
    b.item_id, 
    b.item_code,
    b.description,
    b.quantity
	from orders a, item b
	where a.order_id = b.order_id
	`

	rows, err := conn.Query(sqlSelect)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer rows.Close()

	var dataPerUser entity.Order
	var namaUser string
	checkDiffUser := false

	for rows.Next() {

		type temp struct {
			CustomerName string
			OrderAt      string
			LineItemID   int
			ItemCode     string
			Description  string
			Quantity     int
		}

		var tempData temp

		err = rows.Scan(
			&tempData.CustomerName,
			&tempData.OrderAt,
			&tempData.LineItemID,
			&tempData.ItemCode,
			&tempData.Description,
			&tempData.Quantity,
		)
		if err != nil {
			panic(err)
		}

		//insert to allData when user name different
		if namaUser == "" {
			namaUser = tempData.CustomerName
		} else {
			if namaUser != tempData.CustomerName {
				checkDiffUser = true
			}
		}

		//set inside loop?
		if checkDiffUser {
			allData = append(allData, dataPerUser)
			checkDiffUser = false
		}

		var tempItem entity.Item
		tempItem.LineItemID = tempData.LineItemID
		tempItem.ItemCode = tempData.ItemCode
		tempItem.Description = tempData.Description
		tempItem.Quantity = tempData.Quantity

		if !checkDiffUser {
			dataPerUser.CustomerName = tempData.CustomerName
			dataPerUser.OrderAt = tempData.OrderAt
			dataPerUser.Barang = append(dataPerUser.Barang, tempItem)
		}

	}

	u, _ := json.Marshal(&allData)
	w.Header().Add("Content-Type", "application/json")
	w.Write(u)
	fmt.Fprint(w)
}

func deleteRegisterOrder(w http.ResponseWriter, r *http.Request) {

	fmt.Println("DELETE DATA")

	vars := mux.Vars(r)
	id, _ := vars["orderid"]
	idData, _ := strconv.Atoi(id)
	sqlDelItem := `DELETE from "item" where id=$1`
	_, err := conn.Query(sqlDelItem, idData)
	if err != nil {
		panic(err)
	}

	sqlDelOrders := `DELETE from "orders" where id=$1`
	_, err = conn.Query(sqlDelOrders, idData)
	if err != nil {
		panic(err)
	}

}
