package controller

import (
	"net/http"
	"model"
	"encoding/json"
	"log"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"os"
	"context"
)

func Respond(w http.ResponseWriter, data map[string]interface{})  {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (a *App) CreateUser(w http.ResponseWriter, r *http.Request) {

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

	if resp, ok := user.Create(a.Conn); !ok {
		Respond(w, resp)
	}else {
		Respond(w, resp)
	}
}

func (a *App) Login(w http.ResponseWriter, r *http.Request) {

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

	resp = model.Login(user.Username, user.Password, a.Conn)
	Respond(w, resp)
}

func TokenAuth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		passThrough := []string{"/login", "/user"}

		requestPath := r.URL.Path;
		for _, value := range passThrough {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string] interface{})
		headerString := r.Header.Get("Authorization")
		log.Println(headerString)
		parts := strings.Split(headerString, " ")

		if strings.TrimSpace(headerString) == "" {
			response["status"] = false
			response["message"] = "Missing auth token"
			Respond(w, response)
			return
		}

		if len(parts) != 2 {
			response["status"] = false
			response["message"] = "Invalid token"
			Respond(w, response)
			return
		}

		tokenString := parts[1]

		tk := &model.Token{}

		token, err := jwt.ParseWithClaims(tokenString, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("dbtoken")), nil
		})

		if err != nil {
			response["status"] = false
			response["message"] = "Malformed token"
			response["err"] = err
			Respond(w, response)
			return
		}

		if !token.Valid {
			response["status"] = false
			response["message"] = "Token is invalid"
			Respond(w, response)
			return
		}

		log.Println("Token Authorization succeed for user.")
		log.Println("User ID: ", tk.User, "Username: ", tk.Username)
		ctx := context.WithValue(r.Context(), "user", tk.User)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
