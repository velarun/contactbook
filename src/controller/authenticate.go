package controller

import (
	"net/http"
	"model"
	"encoding/json"
	"log"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	user := &model.User{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("Error Creating User:", err)
	}

	if resp, ok := user.Validate(); !ok {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	if resp, ok := user.Create(model.GetConn()); !ok {
		Respond(w, resp)
	}else {
		Respond(w, resp)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {

	resp := make(map[string]interface{})

	user := &model.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		log.Println("Error Login User:", err)
		resp["status"] = false
		resp["message"] = "Malformed payload"
		Respond(w, resp)
		return
	}

	resp = model.Login(user.Username, user.Password)
	Respond(w, resp)
}

func Respond(w http.ResponseWriter, data map[string]interface{})  {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
