package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/ismaelpereira/origin-challenge/api"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/survey", api.HandleSurvey).Methods("POST")

	fmt.Println("app is listening on port 8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		os.Exit(1)
	}
}
