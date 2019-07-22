package main

import (
	"testing"
	"os"
	"fmt"
	"net/http"
	"bytes"
	"net/http/httptest"
	"controller"
	"encoding/json"
	"encoding/base64"
	"time"
	"log"
	
    "github.com/stretchr/testify/assert"
)

var a controller.App
var m map[string]interface{}

func TestMain(m *testing.M) {
	a = controller.App{}
	a.Initialize()

	ensureTableExists()
	clearTable()

	code := m.Run()

	clearTable()

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

func ensureTableExists() {
	if a.Conn.HasTable("users") && a.Conn.HasTable("contacts") {
		log.Println("Table Exist")
	} else {
		log.Fatalln("Table Not Exist")
	}
}

func clearTable() {
	a.Conn.Exec("TRUNCATE TABLE users")
	a.Conn.Exec("ALTER TABLE users AUTO_INCREMENT = 1")

	a.Conn.Exec("TRUNCATE TABLE contacts")
	a.Conn.Exec("ALTER TABLE contacts AUTO_INCREMENT = 1")
}

func TestCreateUser(t *testing.T) {

	payload := []byte(`{"username": "tester", "email": "tester@testing.com", "password": "tester"}`)
    req, _ := http.NewRequest("DELETE", "/user", bytes.NewBuffer(payload))
    executeRequest(req)

	var m map[string]interface{}
	AcctCreate := "Account has been created."

	req, _ = http.NewRequest("POST", "/user", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &m)
	assert.Equal(t, AcctCreate, m["message"])

	payload = []byte(`{"username": "tester", "email": "tester@testing.com", "password": "tester"}`)
    req, _ = http.NewRequest("DELETE", "/user", bytes.NewBuffer(payload))
    executeRequest(req)
}

func TestCreateUserExistCheck(t *testing.T) {

	payload := []byte(`{"username": "tester", "email": "tester@testing.com", "password": "tester"}`)
    req, _ := http.NewRequest("DELETE", "/user", bytes.NewBuffer(payload))
    executeRequest(req)

	var m map[string]interface{}
	AcctExist := "Username already exists"

	req, _ = http.NewRequest("POST", "/user", bytes.NewBuffer(payload))
	executeRequest(req)

	time.Sleep(2)

	req, _ = http.NewRequest("POST", "/user", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &m)
	assert.Equal(t, AcctExist, m["message"])

	payload = []byte(`{"username": "tester", "email": "tester@testing.com", "password": "tester"}`)
    req, _ = http.NewRequest("DELETE", "/user", bytes.NewBuffer(payload))
    executeRequest(req)
}

func TestCreateUserExistMailCheck(t *testing.T) {

	payload := []byte(`{"username": "tester", "email": "tester@testing.com", "password": "tester"}`)
    req, _ := http.NewRequest("DELETE", "/user", bytes.NewBuffer(payload))
    executeRequest(req)

	var m map[string]interface{}
	AcctMailExist := "Email address already in use"

	req, _ = http.NewRequest("POST", "/user", bytes.NewBuffer(payload))
	executeRequest(req)

	time.Sleep(2)

	payload = []byte(`{"username": "tester1", "email": "tester@testing.com", "password": "tester1"}`)
	req, _ = http.NewRequest("POST", "/user", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &m)
	assert.Equal(t, AcctMailExist, m["message"])

	payload = []byte(`{"username": "tester", "email": "tester@testing.com", "password": "tester"}`)
    req, _ = http.NewRequest("DELETE", "/user", bytes.NewBuffer(payload))
    executeRequest(req)
}

func TestCreateUserEmptyFieldCheck(t *testing.T) {

	payload := []byte(`{"username": "tester", "email": "tester@testing.com", "password": "tester"}`)
    req, _ := http.NewRequest("DELETE", "/user", bytes.NewBuffer(payload))
    executeRequest(req)

	var m map[string]interface{}

	payload3 := []byte(`{"username": "", "email": "tester@testing.com", "password": "tester"}`)
	payload4 := []byte(`{"username": "tester2", "email": "", "password": "tester"}`)
	payload5 := []byte(`{"username": "tester2", "email": "tester2@testing.com", "password": ""}`)

	EmptyUser := "Username cannot be empty"
	EmptyEmail := "Email cannot be empty"
	EmptyPassword := "Password cannot be empty"

	tests := map[string]struct{
        payload  []byte
        expected string
    }{
		"EmptyUsername": { payload: payload3, expected: EmptyUser },
		"EmptyEmail": { payload: payload4, expected: EmptyEmail },
		"EmptyPassword": { payload: payload5, expected: EmptyPassword },
    }

	for name, test := range tests {
        t.Run(name, func(t *testing.T){
			req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(test.payload))
			response := executeRequest(req)

			checkResponseCode(t, http.StatusOK, response.Code)
			json.Unmarshal(response.Body.Bytes(), &m)
			assert.Equal(t, test.expected, m["message"])
			time.Sleep(5)
		})
    }

	payload = []byte(`{"username": "tester", "email": "tester@testing.com", "password": "tester"}`)
    req, _ = http.NewRequest("DELETE", "/user", bytes.NewBuffer(payload))
    executeRequest(req)
}

func TestDeleteUser(t *testing.T) {

	payload := []byte(`{"username": "tester", "email": "tester@testing.com", "password": "tester"}`)
    req, _ := http.NewRequest("DELETE", "/user", bytes.NewBuffer(payload))
    executeRequest(req)

	var m map[string]interface{}

    req, _ = http.NewRequest("POST", "/user", bytes.NewBuffer(payload))
	executeRequest(req)

	time.Sleep(2)
	
	AcctDelete := "Account has been deleted."

	req, _ = http.NewRequest("DELETE", "/user", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &m)
	assert.Equal(t, AcctDelete, m["message"])
	
}

func TestDeleteUserExistCheck(t *testing.T) {

	payload := []byte(`{"username": "tester", "email": "tester@testing.com", "password": "tester"}`)
    req, _ := http.NewRequest("DELETE", "/user", bytes.NewBuffer(payload))
	executeRequest(req)
	
	var m map[string]interface{}

    req, _ = http.NewRequest("POST", "/user", bytes.NewBuffer(payload))
	executeRequest(req)

	time.Sleep(2)
	
	AcctNotExist := "Account not exists"

	req, _ = http.NewRequest("DELETE", "/user", bytes.NewBuffer(payload))
	executeRequest(req)

	time.Sleep(2)

	req, _ = http.NewRequest("DELETE", "/user", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &m)
	assert.Equal(t, AcctNotExist, m["message"])
	
}

func TestDeleteUserEmptyFieldCheck(t *testing.T) {

	var m map[string]interface{}

	payload := []byte(`{"username": "tester", "email": "tester@testing.com", "password": "tester"}`)
    req, _ := http.NewRequest("DELETE", "/user", bytes.NewBuffer(payload))
	executeRequest(req)

	payload2 := []byte(`{"username": "newuser", "email": "noemail@testing.com", "password": "tester"}`)
	payload3 := []byte(`{"username": "", "email": "tester@testing.com", "password": "tester"}`)
	payload4 := []byte(`{"username": "tester", "email": "", "password": "tester"}`)
	payload5 := []byte(`{"username": "tester", "email": "tester@testing.com", "password": ""}`)

	AcctNotExist := "Account not exists"
	EmptyUser := "Username cannot be empty"
	EmptyEmail := "Email cannot be empty"
	EmptyPassword := "Password cannot be empty"

	tests := map[string]struct{
        payload  []byte
        expected string
    }{	
		"EmptyUsername": { payload: payload3, expected: EmptyUser },
		"EmptyEmail": { payload: payload4, expected: EmptyEmail },
		"EmptyPassword": { payload: payload5, expected: EmptyPassword },
		"EmailNotExist": { payload: payload2, expected: AcctNotExist },
    }

	for name, test := range tests {
        t.Run(name, func(t *testing.T){
			req, _ := http.NewRequest("DELETE", "/user", bytes.NewBuffer(test.payload))
			response := executeRequest(req)

			checkResponseCode(t, http.StatusOK, response.Code)
			json.Unmarshal(response.Body.Bytes(), &m)
			assert.Equal(t, test.expected, m["message"])
			time.Sleep(5)
		})
    }
}

func CreateUser(t *testing.T) string {
	
	payload := []byte(`{"username": "tester", "email": "tester@testing.com", "password": "tester"}`)
    req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	token := fmt.Sprintf("Basic %s", basicAuth("tester","tester"))

	return token
}

func DeleteUser() {
	payload := []byte(`{"username": "tester", "email": "tester@testing.com", "password": "tester"}`)
    req, _ := http.NewRequest("DELETE", "/user", bytes.NewBuffer(payload))
	executeRequest(req)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}  

func CreateContact(token string) {
	payload := []byte(`{"contact_name": "user1", "contact_email": "user1@testing.com", "phone_number": "+91 8898883210"}`)
	req, _ := http.NewRequest("POST", "/contact", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", token)
	executeRequest(req)
}

func DeleteContact(token string) {
	payload := []byte(`{"contact_name": "user1", "contact_email": "user1@testing.com", "phone_number": "+91 8898883210"}`)
	req, _ := http.NewRequest("DELETE", "/contact", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", token)
	executeRequest(req)
}

func TestCreateContact(t *testing.T) {

	token := CreateUser(t)

	payload := []byte(`{"contact_name": "user1", "contact_email": "user1@testing.com", "phone_number": "+91 8898883210"}`)
	ContCreate := "Contact has been created."

	req, _ := http.NewRequest("POST", "/contact", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", token)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &m)
	assert.Equal(t, ContCreate, m["message"])
	
	time.Sleep(5)

	DeleteContact(token)	
	DeleteUser()
}

func TestCreateContactExistCheck(t *testing.T) {

	token := CreateUser(t)
	CreateContact(token)

	time.Sleep(5)

	payload := []byte(`{"contact_name": "user2", "contact_email": "user1@testing.com", "phone_number": "+91 8978473874"}`)
	ContExist := "This contact already exists in user's contact list"

	req, _ := http.NewRequest("POST", "/contact", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", token)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &m)
	assert.Equal(t, ContExist, m["message"])
	
	time.Sleep(5)

	DeleteContact(token)	
	DeleteUser()
}

func TestCreateContactIncorrectFieldCheck(t *testing.T) {

	token := CreateUser(t)

	payload5 := []byte(`{"contact_name": "user1user1user1user1user1user1user1user1user1user1user1user1user1user1", "contact_email": "user1@testing.com", "phone_number": "+91 8978473874"}`)
	payload6 := []byte(`{"contact_name": "user1", "contact_email": "user1testing.com", "phone_number": "+91 8978473874"}`)
	payload7 := []byte(`{"contact_name": "user1", "contact_email": "user1@testing.com", "phone_number": "+91 *&      8978473874"}`)
	
	IncorrectPhone := "Phone number has incorrect pattern."
	IncorrectEmail := "Contact email has incorrect pattern."
	IncorrectName := "Length of Contact name should be less 50 characters."

	tests := map[string]struct{
        payload  []byte
        expected string
    }{
		"IncorrectName": { payload: payload5, expected: IncorrectName },
		"IncorrectEmail": { payload: payload6, expected: IncorrectEmail },
		"IncorrectPhone": { payload: payload7, expected: IncorrectPhone },
    }

	for name, test := range tests {
        t.Run(name, func(t *testing.T){
			req, _ := http.NewRequest("POST", "/contact", bytes.NewBuffer(test.payload))
			req.Header.Set("Authorization", token)
			response := executeRequest(req)

			checkResponseCode(t, http.StatusOK, response.Code)
			json.Unmarshal(response.Body.Bytes(), &m)
			assert.Equal(t, test.expected, m["message"])
			time.Sleep(5)
		})
    }
	
	for _, test := range tests {
		req, _ := http.NewRequest("DELETE", "/contact", bytes.NewBuffer(test.payload))
		req.Header.Set("Authorization", token)
		executeRequest(req)
		time.Sleep(5)
	}
	
	DeleteUser()
}

func TestCreateContactEmptyFieldCheck(t *testing.T) {

	token := CreateUser(t)

	payload1 := []byte(`{"contact_name": "", "contact_email": "user1@testing.com", "phone_number": "+91 8978473874"}`)
	payload2 := []byte(`{"contact_name": "user1", "contact_email": "", "phone_number": "+91 8978473874"}`)
	payload3 := []byte(`{"contact_name": "user1", "contact_email": "user1@testing.com", "phone_number": ""}`)
	
	EmptyName := "Contact name is required."
	EmptyEmail := "Contact email is required."
	EmptyPhone := "Phone number is required."

	tests := map[string]struct{
        payload  []byte
        expected string
    }{
		"EmptyContactName": { payload: payload1, expected: EmptyName },
		"EmptyContactEmail": { payload: payload2, expected: EmptyEmail },
		"EmptyContactPhone": { payload: payload3, expected: EmptyPhone },
    }

	for name, test := range tests {
        t.Run(name, func(t *testing.T){
			req, _ := http.NewRequest("POST", "/contact", bytes.NewBuffer(test.payload))
			req.Header.Set("Authorization", token)
			response := executeRequest(req)

			checkResponseCode(t, http.StatusOK, response.Code)
			json.Unmarshal(response.Body.Bytes(), &m)
			assert.Equal(t, test.expected, m["message"])
			time.Sleep(5)
		})
    }
	
	for _, test := range tests {
		req, _ := http.NewRequest("DELETE", "/contact", bytes.NewBuffer(test.payload))
		req.Header.Set("Authorization", token)
		executeRequest(req)
	}
	
	DeleteUser()
}

func TestDeleteContact(t *testing.T) {

	token := CreateUser(t)
	CreateContact(token)

	payload := []byte(`{"contact_name": "user1", "contact_email": "user1@testing.com", "phone_number": "+91 8898883210"}`)
	ContDelete := "Contact has been deleted."

	req, _ := http.NewRequest("DELETE", "/contact", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", token)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &m)
	assert.Equal(t, ContDelete, m["message"])
	
	time.Sleep(5)

	DeleteContact(token)	
	DeleteUser()
}

func TestDeleteContactNotExistCheck(t *testing.T) {

	token := CreateUser(t)
	
	payload := []byte(`{"contact_name": "user2", "contact_email": "user1@testing.com", "phone_number": "+91 8978473874"}`)
	ContNotExist := "This contact not exists in user's contact list"

	req, _ := http.NewRequest("DELETE", "/contact", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", token)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &m)
	assert.Equal(t, ContNotExist, m["message"])
		
	DeleteUser()
}

func TestDeleteContactIncorrectFieldCheck(t *testing.T) {

	token := CreateUser(t)

	payload5 := []byte(`{"contact_name": "user1user1user1user1user1user1user1user1user1user1user1user1user1user1", "contact_email": "user1@testing.com", "phone_number": "+91 8978473874"}`)
	payload6 := []byte(`{"contact_name": "user1", "contact_email": "user1testing.com", "phone_number": "+91 8978473874"}`)
	payload7 := []byte(`{"contact_name": "user1", "contact_email": "user1@testing.com", "phone_number": "+91 *&      8978473874"}`)

	IncorrectPhone := "Phone number has incorrect pattern."
	IncorrectEmail := "Contact email has incorrect pattern."
	IncorrectName := "Length of Contact name should be less 50 characters."

	tests := map[string]struct{
        payload  []byte
        expected string
    }{
		"IncorrectName": { payload: payload5, expected: IncorrectName },
		"IncorrectEmail": { payload: payload6, expected: IncorrectEmail },
		"IncorrectPhone": { payload: payload7, expected: IncorrectPhone },
    }

	for name, test := range tests {
        t.Run(name, func(t *testing.T){
			req, _ := http.NewRequest("DELETE", "/contact", bytes.NewBuffer(test.payload))
			req.Header.Set("Authorization", token)
			response := executeRequest(req)

			checkResponseCode(t, http.StatusOK, response.Code)
			json.Unmarshal(response.Body.Bytes(), &m)
			assert.Equal(t, test.expected, m["message"])
			time.Sleep(5)
		})
    }

	DeleteUser()
}

func TestDeleteContactEmptyFieldCheck(t *testing.T) {

	token := CreateUser(t)

	payload1 := []byte(`{"contact_name": "", "contact_email": "user1@testing.com", "phone_number": "+91 8978473874"}`)
	payload2 := []byte(`{"contact_name": "user1", "contact_email": "", "phone_number": "+91 8978473874"}`)
	payload3 := []byte(`{"contact_name": "user1", "contact_email": "user1@testing.com", "phone_number": ""}`)
	
	EmptyName := "Contact name is required."
	EmptyEmail := "Contact email is required."
	EmptyPhone := "Phone number is required."

	tests := map[string]struct{
        payload  []byte
        expected string
    }{
		"EmptyContactName": { payload: payload1, expected: EmptyName },
		"EmptyContactEmail": { payload: payload2, expected: EmptyEmail },
		"EmptyContactPhone": { payload: payload3, expected: EmptyPhone },
    }

	for name, test := range tests {
        t.Run(name, func(t *testing.T){
			req, _ := http.NewRequest("DELETE", "/contact", bytes.NewBuffer(test.payload))
			req.Header.Set("Authorization", token)
			response := executeRequest(req)

			checkResponseCode(t, http.StatusOK, response.Code)
			json.Unmarshal(response.Body.Bytes(), &m)
			assert.Equal(t, test.expected, m["message"])
			time.Sleep(5)
		})
    }
	
	DeleteUser()
}

func TestSearchExistContactName(t *testing.T) {

	token := CreateUser(t)
	CreateContact(token)

	time.Sleep(5)

	var new map[string]interface{}
	req, _ := http.NewRequest("GET", "/contact?contact_name=user1", nil)
	req.Header.Set("Authorization", token)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &new)

	newnext := m["contact"].(map[string]interface{})
	assert.Equal(t, "user1", newnext["contact_name"])
	
	time.Sleep(5)

	DeleteContact(token)	
	DeleteUser()
}

func TestSearchNotExistContactName(t *testing.T) {
	
	token := CreateUser(t)
	CreateContact(token)
	
	time.Sleep(5)

	req, _ := http.NewRequest("GET", "/contact?contact_name=user00", nil)
	req.Header.Set("Authorization", token)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
	
	time.Sleep(5)

	DeleteContact(token)	
	DeleteUser()
}

func TestSearchExistContactEmail(t *testing.T) {

	token := CreateUser(t)
	CreateContact(token)

	time.Sleep(5)

	var new map[string]interface{}
	req, _ := http.NewRequest("GET", "/contact?contact_email=user1@testing.com", nil)
	req.Header.Set("Authorization", token)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &new)

	newnext := m["contact"].(map[string]interface{})
	assert.Equal(t, "user1@testing.com", newnext["contact_email"])
	
	time.Sleep(5)

	DeleteContact(token)	
	DeleteUser()
}

func TestSearchNotExistContactEmail(t *testing.T) {
	
	token := CreateUser(t)
	CreateContact(token)
	
	time.Sleep(5)

	req, _ := http.NewRequest("GET", "/contact?contact_email=user1esting.com", nil)
	req.Header.Set("Authorization", token)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
	
	time.Sleep(5)

	DeleteContact(token)	
	DeleteUser()
}

func TestPaginationDisplay(t *testing.T) {

	token := CreateUser(t)
	CreateContact(token)

	time.Sleep(5)

	var new []map[string]interface{}
	req, _ := http.NewRequest("GET", "/contact?page=1", nil)
	req.Header.Set("Authorization", token)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &new)

	assert.Equal(t, len(new), 1)
	
	time.Sleep(5)

	DeleteContact(token)	
	DeleteUser()
}

func TestOutofPageCheck(t *testing.T) {
	
	token := CreateUser(t)
	CreateContact(token)
	
	time.Sleep(5)

	req, _ := http.NewRequest("GET", "/contact?page=2", nil)
	req.Header.Set("Authorization", token)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
	
	time.Sleep(5)

	DeleteContact(token)	
	DeleteUser()
}

func TestUpdateContact(t *testing.T) {

	token := CreateUser(t)
	CreateContact(token)

	payload := []byte(`{"contact_name": "userupdated", "contact_email": "user1@testing.com", "phone_number": "+91 8898883210"}`)
	ContCreate := "Contact has been Updated."

	req, _ := http.NewRequest("PUT", "/contact", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", token)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &m)
	assert.Equal(t, ContCreate, m["message"])
	
	//Checks username details on updated data
	new := m["contact"].(map[string]interface{})
	assert.Equal(t, "userupdated", new["contact_name"])
	
	time.Sleep(5)

	req, _ = http.NewRequest("DELETE", "/contact", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", token)
	executeRequest(req)
	
	DeleteUser()
}

func TestUpdateContactNotExistCheck(t *testing.T) {

	token := CreateUser(t)
	
	time.Sleep(5)

	payload := []byte(`{"contact_name": "user2", "contact_email": "user2@testing.com", "phone_number": "+91 8978473874"}`)
	ContNotExist := "This contact not exists in user's contact list"

	req, _ := http.NewRequest("PUT", "/contact", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", token)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	json.Unmarshal(response.Body.Bytes(), &m)
	assert.Equal(t, ContNotExist, m["message"])
	
	time.Sleep(5)

	DeleteContact(token)	
	DeleteUser()
}

func TestUpdateContactIncorrectFieldCheck(t *testing.T) {

	token := CreateUser(t)

	payload5 := []byte(`{"contact_name": "user1user1user1user1user1user1user1user1user1user1user1user1user1user1", "contact_email": "user1@testing.com", "phone_number": "+91 8978473874"}`)
	payload6 := []byte(`{"contact_name": "user1", "contact_email": "user1testing.com", "phone_number": "+91 8978473874"}`)
	payload7 := []byte(`{"contact_name": "user1", "contact_email": "user1@testing.com", "phone_number": "+91 *&      8978473874"}`)
	
	IncorrectPhone := "Phone number has incorrect pattern."
	IncorrectEmail := "Contact email has incorrect pattern."
	IncorrectName := "Length of Contact name should be less 50 characters."

	tests := map[string]struct{
        payload  []byte
        expected string
    }{
		"IncorrectName": { payload: payload5, expected: IncorrectName },
		"IncorrectEmail": { payload: payload6, expected: IncorrectEmail },
		"IncorrectPhone": { payload: payload7, expected: IncorrectPhone },
    }

	for name, test := range tests {
        t.Run(name, func(t *testing.T){
			req, _ := http.NewRequest("PUT", "/contact", bytes.NewBuffer(test.payload))
			req.Header.Set("Authorization", token)
			response := executeRequest(req)

			checkResponseCode(t, http.StatusOK, response.Code)
			json.Unmarshal(response.Body.Bytes(), &m)
			assert.Equal(t, test.expected, m["message"])
			time.Sleep(5)
		})
    }
	
	for _, test := range tests {
		req, _ := http.NewRequest("DELETE", "/contact", bytes.NewBuffer(test.payload))
		req.Header.Set("Authorization", token)
		executeRequest(req)
		time.Sleep(5)
	}
	
	DeleteUser()
}

func TestUpdateContactEmptyFieldCheck(t *testing.T) {

	token := CreateUser(t)

	payload1 := []byte(`{"contact_name": "", "contact_email": "user1@testing.com", "phone_number": "+91 8978473874"}`)
	payload2 := []byte(`{"contact_name": "user1", "contact_email": "", "phone_number": "+91 8978473874"}`)
	payload3 := []byte(`{"contact_name": "user1", "contact_email": "user1@testing.com", "phone_number": ""}`)
	
	EmptyName := "Contact name is required."
	EmptyEmail := "Contact email is required."
	EmptyPhone := "Phone number is required."

	tests := map[string]struct{
        payload  []byte
        expected string
    }{
		"EmptyContactName": { payload: payload1, expected: EmptyName },
		"EmptyContactEmail": { payload: payload2, expected: EmptyEmail },
		"EmptyContactPhone": { payload: payload3, expected: EmptyPhone },
    }

	for name, test := range tests {
        t.Run(name, func(t *testing.T){
			req, _ := http.NewRequest("PUT", "/contact", bytes.NewBuffer(test.payload))
			req.Header.Set("Authorization", token)
			response := executeRequest(req)

			checkResponseCode(t, http.StatusOK, response.Code)
			json.Unmarshal(response.Body.Bytes(), &m)
			assert.Equal(t, test.expected, m["message"])
			time.Sleep(5)
		})
    }
	
	for _, test := range tests {
		req, _ := http.NewRequest("DELETE", "/contact", bytes.NewBuffer(test.payload))
		req.Header.Set("Authorization", token)
		executeRequest(req)
	}
	
	DeleteUser()
}

