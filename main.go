package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
	"github.com/jinzhu/gorm"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"encoding/json"
)

type User struct {
	ID   int `gorm:"primary_key" json:"id"`
	Email string `gorm:"unique_index" json:"email"`
	Name string `json:"name"`
}

//database global var error
var db *gorm.DB
var dbError error

func main(){

	//init router
	port := os.Getenv("PORT")
	router := mux.NewRouter()

	//db config vars
	var dbHost string = os.Getenv("DB_HOST")
	var dbName string = os.Getenv("DB_NAME")
	var dbUser string = os.Getenv("DB_USERNAME")
	var dbPassword string = os.Getenv("DB_PASSWORD")

	//connect to db
	db, dbError = gorm.Open("mysql", dbUser+":"+ dbPassword +"@tcp(" + dbHost+ ":3306)/"+ dbName + "?charset=utf8&parseTime=True&loc=Local")
	if dbError != nil {
		panic("Failed to connect to database")
	}

	//fix for connection timeout
	//see: https://github.com/go-sql-driver/mysql/issues/257
	db.DB().SetMaxIdleConns(0)

	//defer connection
	defer db.Close()

	//handles model updates (no deletes or changes to existing columns)
	db.AutoMigrate(&User{})

	//routes
	router.HandleFunc("/users", UsersHandler).Methods("GET")
	router.HandleFunc("/users/{id}", UserHandler).Methods("GET")
	router.HandleFunc("/users", CreateUserHandler).Methods("POST")

	//create http server
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []User
	//since we're passing a pointer to users, db.Find assigns array to the address
	db.Find(&users)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func UserHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	var user User
	db.First(&user, params["id"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request){
	var user User
	user.Email = r.FormValue("email")
	user.Name = r.FormValue("name")
	db.Create(&user)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&user)
}