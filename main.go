package main

import (
	_ "model"
	"controller"
	"github.com/gorilla/mux"
	"net/http"
	"app"
)
func main()  {

	router := mux.NewRouter()
	router.HandleFunc("/user", controller.CreateUser).Methods("POST")
	router.HandleFunc("/login", controller.Login).Methods("POST")

	router.HandleFunc("/contact", controller.CreateContact).Methods("POST")
	router.HandleFunc("/contact", controller.DeleteContact).Methods("DELETE")
	router.HandleFunc("/contact", controller.UpdateContact).Methods("PUT")
	router.HandleFunc("/contact", controller.SearchContact).Methods("GET")

	router.Use(app.JwtMiddleWare)

	err := http.ListenAndServe("127.0.0.1:9000", router)
	if err != nil {
		panic(err)
	}
}