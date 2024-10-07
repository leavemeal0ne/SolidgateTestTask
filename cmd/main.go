package main

import (
	"github.com/leavemeal0ne/SolidgateTestTask/internal/domen"
	"github.com/leavemeal0ne/SolidgateTestTask/internal/handler"
	"log"
	"net/http"
	"os"
)

func main() {
	//if the data about bank card issuers have not passed validation, the server will not start
	validator, err := domen.InitCardValidator(os.Getenv("validatorDataLocation"))
	if err != nil {
		log.Fatal(err)
	}
	mux := handler.InitHandler(validator).InitRoutes()
	log.Println("server listening on port: ", os.Getenv("serverPort"))
	err = http.ListenAndServe(":"+os.Getenv("serverPort"), mux)
	if err != nil {
		log.Fatal(err)
	}
}
