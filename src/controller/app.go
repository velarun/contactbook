package controller

import (
	"github.com/joho/godotenv"
	"os"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/gorilla/mux"
	"net/http"
	"model"
	"log"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type App struct {
	Router *mux.Router
	Conn *gorm.DB
}

func (a *App) Initialize() {
	
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	dbName := os.Getenv("dbname");
	dbPassword := os.Getenv("dbpass")
	dbUsername := os.Getenv("dbuser")
	dbHost := os.Getenv("dbhost")
	dbPort := os.Getenv("dbport")

	conString := fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s", dbUsername, ":" , dbPassword, "@tcp(", dbHost, ":", dbPort, ")/", dbName, "?parseTime=true")

	fmt.Println(conString)
	var db, errr = gorm.Open("mysql", conString)
	if errr != nil {
		panic(errr)
	}

	db.Debug().AutoMigrate(&model.User{}, &model.Contact{})
	a.Conn = db
	
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/user", a.CreateUser).Methods("POST")
	a.Router.HandleFunc("/user", a.DeleteUser).Methods("DELETE")
	a.Router.HandleFunc("/login", a.Login).Methods("POST")

	a.Router.HandleFunc("/contact", a.CreateContact).Methods("POST")
	a.Router.HandleFunc("/contact", a.DeleteContact).Methods("DELETE")
	a.Router.HandleFunc("/contact", a.UpdateContact).Methods("PUT")
	a.Router.HandleFunc("/contact", a.SearchContact).Methods("GET")

	a.Router.Use(a.BasicAuth)
}

