package main

import (
	"origin-challenge/controller"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}

}

func run() error {
	insuranceParser, err := controller.NewInsuranceParser("./client.json")
	if err != nil {
		return err
	}
	survey, err := insuranceParser.UnmarshallSurvey()
	if err != nil {
		return err
	}
	results, err := insuranceParser.ParseSurvey(survey)
	if err != nil {
		return err
	}
	assignements, err := insuranceParser.SetAssignmentResults(results)
	if err != nil {
		return err
	}
	insuranceJson, err := insuranceParser.MarshallAssignment(assignements)
	if err != nil {
		return err
	}
	spew.Dump(insuranceJson)
	return nil
}
