package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"origin-challenge/controller"

	"github.com/gorilla/mux"
)

func ApiHandler() {
	r := mux.NewRouter()

	r.HandleFunc("/survey", handleSurvey).Methods("POST")

	fmt.Println("app is listening on port 8080")
	http.ListenAndServe(":8080", r)

}

func handleSurvey(w http.ResponseWriter, rq *http.Request) {
	insuranceParser, err := controller.NewInsuranceParser("./client.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	survey, err := insuranceParser.UnmarshallSurvey()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	results, err := insuranceParser.ParseSurvey(survey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	assignements, err := insuranceParser.SetAssignmentResults(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(assignements)

}
