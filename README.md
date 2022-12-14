# Simple web application having CRUD operations under customer object

# Stack
*  Go (GoLang)
*  HTML, JavaScript
*  Bootstrap
*  PostgreSQL


## Setup Guide

* clone the repo
```
https://github.com/vincentvanpro/go-web-app.git
```
* Build Docker
```
in root folder `docker-compose up --build`
```
* Run
```
in root folder `go run main.go`
```
* Open in browser (**NB** Preferably use any browser other than *FireFox* bc of conflicts)
```
http://localhost:8080/
```


## Testing
* Build Docker
```
in root folder `docker-compose up --build`
```
* Run Tests
```
in root `go test ./...`
```