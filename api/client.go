package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"origin-challenge/controller"
	"origin-challenge/types"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
)

func ApiHandler() {
	r := mux.NewRouter()

	r.HandleFunc("/survey", handleSurvey).Methods("POST")

	fmt.Println("app is listening on port 8080")
	http.ListenAndServe(":8080", r)

}

func handleSurvey(w http.ResponseWriter, rq *http.Request) {
	insuranceParser, err := controller.NewInsuranceParser()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var survey types.Survey

	err = json.NewDecoder(rq.Body).Decode(&survey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spew.Dump(survey)

	results, err := insuranceParser.ParseSurvey(&survey)
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
