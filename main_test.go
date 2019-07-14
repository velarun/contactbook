package main

import (
	"testing"
	"os"
	"fmt"
	"net/http"
	"bytes"
	"net/http/httptest"
	"controller"
    //"github.com/stretchr/testify/assert"
)

var a controller.App

func TestMain(m *testing.M) {
	a = controller.App{}
	a.Initialize()

	code := m.Run()

	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    a.Router.ServeHTTP(rr, req)

    return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
    if expected != actual {
        t.Errorf("Expected response code %d. Got %d\n", expected, actual)
    }
}

func TestCreateUser(t *testing.T) {

	payload := []byte(`{"username": "tester", "email": "tester@testing.com", "password": "tester"}`)
    req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(payload))
    response := executeRequest(req)

	fmt.Println(response)
}
