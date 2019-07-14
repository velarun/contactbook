package app

import (
	"net/http"
	"controller"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"os"
	"model"
	"context"
	"log"
)

func JwtMiddleWare(next http.Handler) http.Handler {

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
		parts := strings.Split(headerString, " ")

		if strings.TrimSpace(headerString) == "" {
			response["status"] = false
			response["message"] = "Missing auth token"
			controller.Respond(w, response)
			return
		}

		if len(parts) != 2 {
			response["status"] = false
			response["message"] = "Invalid token"
			controller.Respond(w, response)
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
			controller.Respond(w, response)
			return
		}

		if !token.Valid {
			response["status"] = false
			response["message"] = "Token is invalid"
			controller.Respond(w, response)
			return
		}

		log.Println("Token Authorization succeed for user.")
		log.Println("User ID: ", tk.User, "Username: ", tk.Username)
		ctx := context.WithValue(r.Context(), "user", tk.User)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
