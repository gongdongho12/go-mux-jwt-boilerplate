package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"os"
	"github.com/jpfuentes2/go-env"
	"path"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
)

func main(){

	//read env
	pwd, _ := os.Getwd()
	env.ReadEnv(path.Join(pwd, ".env"))

	//connect to db
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	db, err := gorm.Open("mysql", dbUser+":"+ dbPassword +"@tcp(" + dbHost+ ":3306)/"+ dbName)
	defer db.Close()

	//init router
	port := os.Getenv("PORT")
	router := mux.NewRouter()

	if err != nil {
		fmt.Println(err)
	}

	router.HandleFunc("/", HomeHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	//load up the env package only in local env
	for _, v := range os.Environ() {
		w.Write([]byte("env: " + v))
	}

	r.Body.Close()

}