package controller

import (
	"encoding/json"
	"io/ioutil"
	"origin-challenge/types"
	"os"
	"time"
)

type insuranceParser struct {
	path string
}

type InsuranceParser interface {
	UnmarshallSurvey() (*types.Survey, error)
	MarshallAssignment(a *types.Assignment) ([]byte, error)
	ParseSurvey(survey *types.Survey) (*types.SurveyResults, error)
	SetAssignmentResults(sr *types.SurveyResults) (*types.Assignment, error)
}

func NewInsuranceParser(path string) (InsuranceParser, error) {
	return &insuranceParser{
		path: path,
	}, nil
}

func (t *insuranceParser) UnmarshallSurvey() (*types.Survey, error) {
	file, err := os.Open(t.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	surveyFile, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var surveyDecoded types.Survey
	if err = json.Unmarshal(surveyFile, &surveyDecoded); err != nil {
		return nil, err
	}
	return &surveyDecoded, nil
}

func (t *insuranceParser) MarshallAssignment(a *types.Assignment) ([]byte, error) {

	parsedAssignment, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return parsedAssignment, nil
}

func (t *insuranceParser) ParseSurvey(s *types.Survey) (*types.SurveyResults, error) {
	var results types.SurveyResults
	var riskPoints int

	for _, riskAnswers := range s.RiskQuestions {
		riskPoints += riskAnswers
	}

	results.DisabilityPoints = riskPoints
	results.HomePoints = riskPoints
	results.LifePoints = riskPoints
	results.VehiclePoints = riskPoints

	if s.Income == 0 {
		results.DisabilityPoints = 500
	}

	if s.Vehicle == (types.Vehicle{}) {
		results.VehiclePoints = 500
	}

	if s.House == (types.House{}) {
		results.HomePoints = 500
	}

	if s.Age > 60 {
		results.DisabilityPoints = 500
		results.LifePoints = 500
	}

	if s.Age < 30 {
		results.VehiclePoints -= 2
		results.DisabilityPoints -= 2
		results.HomePoints -= 2
		results.LifePoints -= 2
	}

	if s.Age >= 30 && s.Age <= 40 {
		results.VehiclePoints -= 1
		results.DisabilityPoints -= 1
		results.HomePoints -= 1
		results.LifePoints -= 1
	}

	if s.Income > 200000 {
		results.VehiclePoints -= 1
		results.DisabilityPoints -= 1
		results.HomePoints -= 1
		results.LifePoints -= 1
	}

	if s.House.OwnershipStatus == "mortgaged" {
		results.HomePoints += 1
		results.DisabilityPoints += 1
	}

	if s.Dependents > 0 {
		results.DisabilityPoints += 1
		results.LifePoints += 1
	}

	if s.MaritialStatus == "married" {
		results.LifePoints += 1
		results.DisabilityPoints -= 1
	}

	if s.Vehicle.Year-time.Now().Year() <= 5 {
		results.VehiclePoints += 1
	}

	return &results, nil
}

func (t *insuranceParser) SetAssignmentResults(sr *types.SurveyResults) (*types.Assignment, error) {
	var assignment types.Assignment

	if sr.DisabilityPoints <= 0 {
		assignment.Disability = "economic"
	}

	if sr.DisabilityPoints > 0 && sr.DisabilityPoints < 3 {
		assignment.Disability = "regular"
	}

	if sr.DisabilityPoints >= 3 {
		assignment.Disability = "responsible"
	}

	if sr.DisabilityPoints >= 400 {
		assignment.Disability = "ineligible"
	}

	if sr.HomePoints <= 0 {
		assignment.Home = "economic"
	}

	if sr.HomePoints > 0 && sr.HomePoints < 3 {
		assignment.Home = "regular"
	}

	if sr.HomePoints >= 3 {
		assignment.Home = "responsible"
	}

	if sr.HomePoints >= 400 {
		assignment.Home = "ineligible"
	}

	if sr.LifePoints <= 0 {
		assignment.Life = "economic"
	}

	if sr.LifePoints > 0 && sr.LifePoints < 3 {
		assignment.Life = "regular"
	}

	if sr.LifePoints >= 3 {
		assignment.Life = "responsible"
	}
	if sr.LifePoints >= 400 {
		assignment.Life = "ineligible"
	}

	if sr.VehiclePoints <= 0 {
		assignment.Vehicle = "economic"
	}

	if sr.VehiclePoints > 0 && sr.VehiclePoints < 3 {
		assignment.Vehicle = "regular"
	}

	if sr.VehiclePoints >= 3 {
		assignment.Vehicle = "responsible"
	}
	if sr.VehiclePoints >= 400 {
		assignment.Vehicle = "ineligible"
	}

	return &assignment, nil
}
