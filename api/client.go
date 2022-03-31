package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/ismaelpereira/origin-challenge/controller"
	"github.com/ismaelpereira/origin-challenge/types"

	"github.com/davecgh/go-spew/spew"
)

func HandleSurvey(w http.ResponseWriter, rq *http.Request) {
	insuranceParser, err := controller.NewInsuranceParser()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		os.Exit(1)
	}

	var survey types.Survey

	err = json.NewDecoder(rq.Body).Decode(&survey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		os.Exit(1)
	}

	spew.Dump(survey)

	results, err := insuranceParser.ParseSurvey(&survey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		os.Exit(1)
	}
	assignements, err := insuranceParser.SetAssignmentResults(&results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		os.Exit(1)
	}

	err = json.NewEncoder(w).Encode(assignements)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		os.Exit(1)
	}

}
