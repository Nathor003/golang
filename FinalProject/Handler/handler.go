package handler

import (
	"CRUD_POSTGRE_GITIGNORE/entity"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func emailChecker(email string) (bool, string) {
	if email == "" {
		return false, "email is empty"
	}

	if !strings.Contains(email, "@mail.com") {
		return false, "Invalid email format"
	}
	selectSql := `select 1 from "users" where email=$1`
	var checkUsed int
	rows, err := conn.Query(selectSql, email)
	if err != nil {
		return false, "error"
	}
	for rows.Next() {
		err := rows.Scan(
			&checkUsed,
		)
		if err != nil {
			return false, "error"
		}
	}
	if checkUsed < 1 {
		return true, "OK"
	} else {
		return false, "email already used"
	}

}

func passwordChecker(password string) (bool, string) {

	if password == "" {
		return false, "Password is empty"
	}
	if len(password) < 6 {
		return false, "Password not long enough"
	}

	return true, "OK"
}

func usernameChecker(username string) (bool, string) {

	if username == "" {
		return false, "Username is empty"
	}

	querySelect := `select 1 from "users" where username = $1`
	rows, err := conn.Query(querySelect, username)
	if err != nil {
		return false, "error"
	}

	var checkAvailable int
	for rows.Next() {
		err := rows.Scan(
			&checkAvailable,
		)
		if err != nil {
			return false, "error"
		}
	}

	rows.Scan(&checkAvailable)
	if checkAvailable < 1 {
		return true, "OK"
	} else {
		return false, "Username already used"
	}
}

func ageChecker(age int) (bool, string) {

	if age == 0 {
		return false, "Age is empty"
	}

	if age < 8 {
		return false, "Not Older Enough"
	}

	return true, "OK"
}

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

type Claims struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func tokenCheck() (string, error) {
	var jsonFile struct {
		Token string
	}

	content, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(content, &jsonFile)
	if err != nil {
		return "", err
	}

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(jsonFile.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secretKey"), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", err
		}

		return "", err
	}
	if !tkn.Valid {
		return "", errors.New("token invalid")
	}

	fmt.Println("Claims : ", claims)
	return claims.Username, nil
}

var conn *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "db-go-sql"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	conn = connectDB()
	defer conn.Close()

	fmt.Println("Path : ", r.URL.Path, r.Method)

	switch r.Method {
	case "POST":
		switch r.URL.Path {
		case "/users/register":
			registerUser(w, r)
		case "/users/login":
			loginUser(w, r)
		case "/photos":
			uploadPhoto(w, r)
		case "/comments":
			postComments(w, r)
		case "/socialmedias":
			postSocialMedia(w, r)
		}
	case "PUT":
		if r.URL.Path == "/users" {
			updateDataUser(w, r)
		} else if strings.Contains(r.URL.Path, "/photos/") {
			updatePhoto(w, r)
		} else if strings.Contains(r.URL.Path, "/comments/") {
			updateComments(w, r)
		} else if strings.Contains(r.URL.Path, "/socialmedias/") {
			updateSocialMedia(w, r)
		}
	case "DELETE":
		if r.URL.Path == "/users" {
			deleteUser(w, r)
		} else if strings.Contains(r.URL.Path, "/photos/") {
			deletePhoto(w, r)
		} else if strings.Contains(r.URL.Path, "/comments/") {
			deleteComment(w, r)
		} else if strings.Contains(r.URL.Path, "/socialmedias/") {
			deleteSocialMedia(w, r)
		}
	case "GET":
		switch r.URL.Path {
		case "/photos":
			getPhoto(w, r)
		case "/comments":
			getCommnets(w, r)
		case "/socialmedias":
			getSocialMedia(w, r)
		}
	}

}

//------------------------------------------------------------
//USERS
func registerUser(w http.ResponseWriter, r *http.Request) {

	var register struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Age      int    `json:"age"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&register); err != nil {
		w.Write([]byte("Input invalid"))
		fmt.Println(err)
		return
	}
	//check is username valid
	check, message := usernameChecker(register.Username)
	if !check {
		w.Write([]byte(message))
		return
	}

	//check is password valid
	check, message = passwordChecker(register.Password)
	if !check {
		w.Write([]byte(message))
		return
	}
	pass := []byte(register.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		w.Write([]byte("ERROR"))
		return
	}
	register.Password = string(hashedPassword)

	//check is email valid
	check, message = emailChecker(register.Email)
	if !check {
		w.Write([]byte(message))
		return
	}
	check, message = ageChecker(register.Age)
	if !check {
		w.Write([]byte(message))
		return
	}
	dt := time.Now()
	//valid data input
	var idUser int
	insertQuery := `insert into "users"(username,password,email,age,create_time) values ($1,$2,$3,$4,$5) returning id`
	rows, err := conn.Query(insertQuery, register.Username, register.Password, register.Email, register.Age, dt.Format("2006-01-02 15:04:05"))
	if err != nil {
		w.Write([]byte("Error insert data"))
	}
	for rows.Next() {
		rows.Scan(&idUser)
	}

	var returnData struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Age      int    `json:"age"`
	}
	returnData.Id = idUser
	returnData.Username = register.Username
	returnData.Email = register.Email
	returnData.Age = register.Age

	u, _ := json.Marshal(&returnData)
	w.Header().Add("Content-Type", "application/json")
	w.Write(u)
}

func loginUser(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Login User")

	var userpass struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userpass); err != nil {
		w.Write([]byte("Input invalid"))
		fmt.Println(err)
		return
	}
	fmt.Println(userpass)
	sqlSelect := `select username,password from "users" where email=$1`

	rows, err := conn.Query(sqlSelect, userpass.Email)
	if err != nil {
		w.Write([]byte("Error execute query"))
		return
	}

	var hashedPass string
	var username string
	for rows.Next() {
		rows.Scan(&username, &hashedPass)
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(userpass.Password))
	if err != nil {
		w.Write([]byte("Password incorrect"))
		return
	}

	//make jwt
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email:    userpass.Email,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secretKey"))
	if err != nil {
		w.Write([]byte("Errror make Toke"))
		return
	}

	var JWTToken struct {
		Token string `json:"token"`
	}
	JWTToken.Token = tokenString
	u, _ := json.Marshal(&JWTToken)
	w.Header().Add("Content-Type", "application/json")
	w.Write(u)

}

func updateDataUser(w http.ResponseWriter, r *http.Request) {
	username, err := tokenCheck()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var Data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&Data); err != nil {
		w.Write([]byte("Input invalid"))
		fmt.Println(err)
		return
	}
	valid, mess := emailChecker(Data.Email)
	if !valid {
		w.Write([]byte(mess))
		return
	}
	valid, mess = passwordChecker(Data.Password)
	if !valid {
		w.Write([]byte(mess))
		return
	}
	pass := []byte(Data.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		w.Write([]byte("ERROR"))
		return
	}
	var returnData entity.User

	udpateQuery := `update users set email=$1, password=$2,update_time=$3 where username=$4 returning id,email,username,age,update_time`
	dt := time.Now()
	rows, err := conn.Query(udpateQuery, Data.Email, hashedPassword, dt.Format("2006-01-02 15:04:05"), username)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	for rows.Next() {
		rows.Scan(
			&returnData.Id,
			&returnData.Email,
			&returnData.Username,
			&returnData.Age,
			&returnData.UpdatedAt,
		)
	}

	u, _ := json.Marshal(&returnData)
	w.Header().Add("Content-Type", "application/json")
	w.Write(u)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	username, err := tokenCheck()
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	deleteQuery := `delete from users where username = $1`
	_, err = conn.Query(deleteQuery, username)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	var Mess struct {
		Notif string `json:"message"`
	}

	Mess.Notif = "Your account has been successfully deleted"

	u, _ := json.Marshal(&Mess)
	w.Header().Add("Content-Type", "application/json")
	w.Write(u)
}

//------------------------------------------------------------
//------------------------------------------------------------
//PHOTOS

func uploadPhoto(w http.ResponseWriter, r *http.Request) {

	username, err := tokenCheck()
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	var Input struct {
		Title    string `json:"title"`
		Caption  string `json:"caption"`
		PhotoUrl string `json:"photo_url"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&Input); err != nil {
		w.Write([]byte("Input invalid"))
		fmt.Println(err)
		return
	}

	if Input.PhotoUrl == "" || Input.Title == "" {
		w.Write([]byte("Title or Photo Url must be filled"))
		return
	}
	dt := time.Now()
	insertQuery := `insert into "Photo" (title,caption,photo_url,created_at,userid) 
				   values($1,$2,$3,$4,(select id from "users" where username=$5)) 
				   returning  id,title,caption,photo_url,created_at,userid`
	rows, err := conn.Query(insertQuery, Input.Title, Input.Caption, Input.PhotoUrl, dt.Format("2006-01-02 15:04:05"), username)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	var Result struct {
		Id         int       `json:"id"`
		Title      string    `json:"title"`
		Caption    string    `json:"caption"`
		PhotoUrl   string    `json:"photo_url"`
		Created_at time.Time `json:"date"`
		UserId     int       `json:"user_id"`
	}
	for rows.Next() {
		rows.Scan(
			&Result.Id,
			&Result.Title,
			&Result.Caption,
			&Result.PhotoUrl,
			&Result.Created_at,
			&Result.UserId,
		)
	}

	u, _ := json.Marshal(&Result)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(u)
}

func getPhoto(w http.ResponseWriter, r *http.Request) {
	username, err := tokenCheck()
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	type Users struct {
		Email    string `json:"email"`
		Username string `json:"username"`
	}
	type Result struct {
		Id        int       `json:"id"`
		Title     string    `json:"title"`
		Caption   string    `json:"caption"`
		PhotoUrl  string    `json:"photo_url"`
		Createdat time.Time `json:"created_at"`
		UserId    int       `json:"user_id"`
		UpdatedAt time.Time `json:"updated_at"`
		User      Users     `json:"User"`
	}

	selectQuery := `select b.email,b.username,a.id,a.title,a.caption,a.photo_url,a.userid,a.created_at,a.updated_at
					from "Photo" a,"users" b 
					where b.id=(select c.id from "users" c where c.username=$1)`

	rows, err := conn.Query(selectQuery, username)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	var res []Result

	for rows.Next() {
		var temp Result
		rows.Scan(
			&temp.User.Email,
			&temp.User.Username,
			&temp.Id,
			&temp.Title,
			&temp.Caption,
			&temp.PhotoUrl,
			&temp.UserId,
			&temp.Createdat,
			&temp.UpdatedAt,
		)
		res = append(res, temp)
	}

	u, _ := json.Marshal(res)
	w.Header().Add("Content-Type", "application/json")
	w.Write(u)
}

func updatePhoto(w http.ResponseWriter, r *http.Request) {
	username, err := tokenCheck()
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	id := strings.TrimPrefix(r.URL.Path, "/photos/")
	idData, _ := strconv.Atoi(id)

	var Input struct {
		Title    string `json:"title"`
		Caption  string `json:"caption"`
		PhotoUrl string `json:"photo_url"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&Input); err != nil {
		w.Write([]byte("Input invalid"))
		fmt.Println(err)
		return
	}

	if Input.PhotoUrl == "" || Input.Title == "" {
		w.Write([]byte("Title or Photo Url must be filled"))
		return
	}
	dt := time.Now()
	insertQuery := `update "Photo" set
					title=$1,
					caption=$2,
					photo_url=$3,
					updated_at=$4
					where userid = (select id from "users" where username=$5)
					and id = $6
				   returning title,caption,photo_url,updated_at,userid,id`
	rows, err := conn.Query(insertQuery, Input.Title, Input.Caption, Input.PhotoUrl, dt.Format("2006-01-02 15:04:05"), username, idData)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	var Result struct {
		Id         int       `json:"id"`
		Title      string    `json:"title"`
		Caption    string    `json:"caption"`
		PhotoUrl   string    `json:"photo_url"`
		Updated_at time.Time `json:"updated_at"`
		UserId     int       `json:"user_id"`
	}
	for rows.Next() {
		rows.Scan(
			&Result.Title,
			&Result.Caption,
			&Result.PhotoUrl,
			&Result.Updated_at,
			&Result.UserId,
			&Result.Id,
		)
	}

	u, _ := json.Marshal(&Result)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(u)
}

func deletePhoto(w http.ResponseWriter, r *http.Request) {
	username, err := tokenCheck()
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	id := strings.TrimPrefix(r.URL.Path, "/photos/")
	idData, _ := strconv.Atoi(id)

	deleteQuery := `delete from "Photo" where userid=(select id from "users" where username = $1) and id=$2`
	_, err = conn.Query(deleteQuery, username, idData)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	var Mess struct {
		Notif string `json:"message"`
	}

	Mess.Notif = "Your Photo has been successfully deleted"

	u, _ := json.Marshal(&Mess)
	w.Header().Add("Content-Type", "application/json")
	w.Write(u)
}

//-------------------------------------------------------------
//Comments

func postComments(w http.ResponseWriter, r *http.Request) {
	username, err := tokenCheck()
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	var Input struct {
		Message  string `json:"message"`
		Photo_id int    `json:"photo_id"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&Input); err != nil {
		w.Write([]byte("Input invalid"))
		fmt.Println(err)
		return
	}
	if Input.Message == "" {
		w.Write([]byte("Message must be filled"))
		return
	}

	sqlInsert := `insert into "Comment"(message,photoid,userid,created_at)
					values(
					$1,
					$2,
					(select id from "users" where username=$3),   
					$4)
					returning id,message,photoid,userid,created_at`

	dt := time.Now()

	var Result struct {
		Id        int       `json:"id"`
		Message   string    `json:"string"`
		Photoid   int       `json:"photo_id"`
		Userid    int       `json:"user_id"`
		Createdat time.Time `json:"created_at"`
	}
	rows, err := conn.Query(sqlInsert, Input.Message, Input.Photo_id, username, dt.Format("2006-01-02 15:04:05"))
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	for rows.Next() {
		rows.Scan(
			&Result.Id,
			&Result.Message,
			&Result.Photoid,
			&Result.Userid,
			&Result.Createdat,
		)
	}

	u, _ := json.Marshal(&Result)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(u)
}

func getCommnets(w http.ResponseWriter, r *http.Request) {
	username, err := tokenCheck()
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	selectSql := `select  
					c.id,
					c.title,
					c.caption,
					c.photo_url,
					c.userid,
					b.id,
					b.email,
					b.username,
					a.id,
					a.message,
					a.photoid,
					a.userid,
					a.updated_at,
					a.created_at
				from "Comment" a,"users" b,"Photo" c
				where c.userid=b.id
				and c.id=a.photoid
				and b.id = (select id from "users" where username=$1)`

	type User struct {
		Id       int    `json:"id"`
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	type Photo struct {
		Id       int    `json:"id"`
		Title    string `json:"title"`
		Caption  string `json:"caption"`
		PhotoUrl string `json:"photourl"`
		UserId   int    `json:"userid"`
	}

	type Comment struct {
		Id        int       `json:"id"`
		Message   string    `json:"message"`
		Photoid   string    `json:"photo_id"`
		Userid    string    `json:"user_id"`
		UpdatedAt time.Time `json:"updated_at"`
		CreatedAt time.Time `json:"created_at"`
		User      User      `json:"User"`
		Photo     Photo     `json:"Photo"`
	}

	rows, err := conn.Query(selectSql, username)
	var Result []Comment
	for rows.Next() {
		var temp Comment
		rows.Scan(
			&temp.Photo.Id,
			&temp.Photo.Title,
			&temp.Photo.Caption,
			&temp.Photo.PhotoUrl,
			&temp.Photo.UserId,
			&temp.User.Id,
			&temp.User.Email,
			&temp.User.Username,
			&temp.Id,
			&temp.Message,
			&temp.Photoid,
			&temp.Userid,
			&temp.UpdatedAt,
			&temp.CreatedAt,
		)
		Result = append(Result, temp)
	}
	u, _ := json.Marshal(&Result)
	w.Header().Add("Content-Type", "application/json")
	w.Write(u)
}

func updateComments(w http.ResponseWriter, r *http.Request) {
	username, err := tokenCheck()
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	var Input struct {
		Message string `json:"message"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&Input); err != nil {
		w.Write([]byte("Input invalid"))
		fmt.Println(err)
		return
	}
	if Input.Message == "" {
		w.Write([]byte("Message must be filled"))
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/comments/")
	idData, _ := strconv.Atoi(id)

	dt := time.Now()
	sqlUpdate := `update "Comment" set message=$1,updated_at=$2 where id=$3 and userid=(select id from "users" where username=$4)
					returning photoid`
	rows, err := conn.Query(sqlUpdate, Input.Message, dt.Format("2006-01-02 15:04:05"), idData, username)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	var photoId int
	for rows.Next() {
		rows.Scan(
			&photoId,
		)
	}
	fmt.Println("ID Photo : ", photoId)
	fmt.Println("Select data")

	sqlSelect := `select id,title,caption,photo_url,updated_at,userid from "Photo" where id=$1 and userid=(select id from "users" where username=$2)`
	rows, err = conn.Query(sqlSelect, photoId, username)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	var Result struct {
		Id         int       `json:"id"`
		Title      string    `json:"title"`
		Caption    string    `json:"caption"`
		PhotoUrl   string    `json:"photo_url"`
		Updated_at time.Time `json:"updated_at"`
		UserId     int       `json:"user_id"`
	}
	for rows.Next() {
		rows.Scan(
			&Result.Id,
			&Result.Title,
			&Result.Caption,
			&Result.PhotoUrl,
			&Result.Updated_at,
			&Result.UserId,
		)
	}

	u, _ := json.Marshal(&Result)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(u)
}

func deleteComment(w http.ResponseWriter, r *http.Request) {
	username, err := tokenCheck()
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	id := strings.TrimPrefix(r.URL.Path, "/comments/")
	idData, _ := strconv.Atoi(id)

	deleteQuery := `delete from "Comment" where userid=(select id from "users" where username = $1) and id=$2`
	_, err = conn.Query(deleteQuery, username, idData)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var Mess struct {
		Notif string `json:"message"`
	}

	Mess.Notif = "Your Comments has been successfully deleted"

	u, _ := json.Marshal(&Mess)
	w.Header().Add("Content-Type", "application/json")
	w.Write(u)
}

//----------------------------------------------------------
//SOCIAL MEDIA

func postSocialMedia(w http.ResponseWriter, r *http.Request) {
	username, err := tokenCheck()
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	var Input struct {
		Name        string `json:"name"`
		SocialMedia string `json:"social_media_url"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&Input); err != nil {
		w.Write([]byte("Input invalid"))
		fmt.Println(err)
		return
	}
	if Input.SocialMedia == "" {
		w.Write([]byte("SocialMedia must be filled"))
		return
	}
	dt := time.Now()
	insertQuery := `insert into "SocialMedia"(name,social_media_url,created_at,userid) values
	(
		$1,
		$2,
		$3,
		(select id from "users" where username=$4)
	) returning id,name,social_media_url,userid,created_at`

	rows, err := conn.Query(insertQuery, Input.Name, Input.SocialMedia, dt.Format("2006-01-02 15:04:05"), username)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var Result struct {
		Id             int       `json:"id"`
		Name           string    `json:"name"`
		SocialMediaUrl string    `json:"social_media_url"`
		Userid         int       `json:"user_id"`
		CreatedAt      time.Time `json:"created_at"`
	}

	for rows.Next() {
		rows.Scan(
			&Result.Id,
			&Result.Name,
			&Result.SocialMediaUrl,
			&Result.Userid,
			&Result.CreatedAt,
		)
	}

	u, _ := json.Marshal(&Result)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(u)
}

func getSocialMedia(w http.ResponseWriter, r *http.Request) {
	username, err := tokenCheck()
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	sqlSelect := `select 
	(select photo_url from "Photo" where userid=(select id from "users" where username='Orvin') limit 1), 
    b.id,
    b.username,
    a.id,
    a.name,
    a.social_media_url,
    a.userid,
    a.created_at,
    a.updated_at	
from "SocialMedia" a,"users" b
where a.userid=b.id
and b.username = $1`

	rows, err := conn.Query(sqlSelect, username)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	type User struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
		Pp       string `json:"profile_image_url"`
	}
	type SocialMedia struct {
		Id        int       `json:"id"`
		Name      string    `json:"name"`
		Smu       string    `json:"social_media_url"`
		Userid    int       `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		User      User      `json:"User"`
	}

	var ReturnResult struct {
		Result []SocialMedia `json:"social_media"`
	}
	for rows.Next() {
		var temp SocialMedia
		rows.Scan(
			&temp.User.Pp,
			&temp.User.Id,
			&temp.User.Username,
			&temp.Id,
			&temp.Name,
			&temp.Smu,
			&temp.Userid,
			&temp.CreatedAt,
			&temp.UpdatedAt,
		)
		ReturnResult.Result = append(ReturnResult.Result, temp)
	}

	u, _ := json.Marshal(&ReturnResult)
	w.Header().Add("Content-Type", "application/json")
	w.Write(u)
}

func updateSocialMedia(w http.ResponseWriter, r *http.Request) {
	username, err := tokenCheck()
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	var Input struct {
		Name        string `json:"name"`
		SocialMedia string `json:"social_media_url"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&Input); err != nil {
		w.Write([]byte("Input invalid"))
		fmt.Println(err)
		return
	}
	if Input.SocialMedia == "" {
		w.Write([]byte("SocialMedia must be filled"))
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/socialmedias/")
	idData, _ := strconv.Atoi(id)

	dt := time.Now()
	insertQuery := `Update "SocialMedia"
					set
						name=$1,
						social_media_url=$2,
						updated_at = $3
					where id = $4
					and userid = (select id from "users" where username=$5)
					returning id,name,social_media_url,userid,created_at`

	rows, err := conn.Query(insertQuery, Input.Name, Input.SocialMedia, dt.Format("2006-01-02 15:04:05"), idData, username)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var Result struct {
		Id             int       `json:"id"`
		Name           string    `json:"name"`
		SocialMediaUrl string    `json:"social_media_url"`
		Userid         int       `json:"user_id"`
		UpdatedAt      time.Time `json:"updated_at"`
	}

	for rows.Next() {
		rows.Scan(
			&Result.Id,
			&Result.Name,
			&Result.SocialMediaUrl,
			&Result.Userid,
			&Result.UpdatedAt,
		)
	}

	u, _ := json.Marshal(&Result)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(u)
}

func deleteSocialMedia(w http.ResponseWriter, r *http.Request) {
	username, err := tokenCheck()
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	id := strings.TrimPrefix(r.URL.Path, "/socialmedias/")
	idData, _ := strconv.Atoi(id)

	deleteQuery := `delete from "SocialMedia" where userid=(select id from "users" where username = $1) and id=$2`
	_, err = conn.Query(deleteQuery, username, idData)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	var Mess struct {
		Notif string `json:"message"`
	}

	Mess.Notif = "Your SocialMedia has been successfully deleted"

	u, _ := json.Marshal(&Mess)
	w.Header().Add("Content-Type", "application/json")
	w.Write(u)
}
