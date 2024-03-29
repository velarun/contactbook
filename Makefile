PROJECT_NAME := Contact Book App

.PHONY: help
help:
	@echo "------------------------------------------------------------------------"
	@echo "${PROJECT_NAME}"
	@echo "------------------------------------------------------------------------"
	@grep -E '^[a-zA-Z0-9_/%\-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## build application binaries
	GOOS=darwin GOARCH=amd64 go build -o contactbook-darwin-amd64 .
	GOOS=linux GOARCH=amd64 go build -o contactbook-linux-amd64 .

.PHONY: dep
dep: ## install latest build of dependency manager and linters
	go get -u github.com/jinzhu/gorm
	go get -u github.com/go-sql-driver/mysql
	go get -u golang.org/x/crypto/bcrypt
	go get -u github.com/dgrijalva/jwt-go
	go get -u github.com/gorilla/mux
	go get -u github.com/gemcook/pagination-go
	go get -u github.com/joho/godotenv
	go get -u github.com/stretchr/testify/assert

.PHONY: lint
lint: ## check code for lint errors
	go vet ./...

.PHONY: test
test: ## run unit tests
	go test -v
