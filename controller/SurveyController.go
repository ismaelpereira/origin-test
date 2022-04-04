package controller

import (
	"encoding/json"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/ismaelpereira/origin-challenge/types"
)

type insuranceParser struct {
}

type InsuranceParser interface {
	UnmarshallSurvey(jsonSurvey []byte) (types.Survey, error)
	MarshallAssignment(a *types.Assignment) ([]byte, error)
	ParseSurvey(survey *types.Survey) (types.SurveyResults, error)
	SetAssignmentResults(sr *types.SurveyResults) types.Assignment
}

func NewInsuranceParser() (InsuranceParser, error) {
	return &insuranceParser{}, nil
}

func (t *insuranceParser) UnmarshallSurvey(jsonSurvey []byte) (types.Survey, error) {
	var surveyDecoded types.Survey
	if err := json.Unmarshal(jsonSurvey, &surveyDecoded); err != nil {
		return types.Survey{}, err
	}
	return surveyDecoded, nil
}

func (t *insuranceParser) MarshallAssignment(a *types.Assignment) ([]byte, error) {
	parsedAssignment, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return parsedAssignment, nil
}

func (t *insuranceParser) sumRiskQuestions(s *types.Survey) int {
	var riskPoints int

	for _, riskAnswers := range s.RiskQuestions {
		riskPoints += riskAnswers
	}
	return riskPoints
}

func (t *insuranceParser) verifyIneligibility(s *types.Survey, r types.SurveyResults) types.SurveyResults {
	if s.Income == 0 {
		r.DisabilityPoints = 500
	}

	if s.Vehicle == (types.Vehicle{}) {
		r.VehiclePoints = 500
	}

	if s.House == (types.House{}) {
		r.HomePoints = 500
	}

	if s.Age > 60 {
		r.DisabilityPoints = 500
		r.LifePoints = 500
	}

	return r
}

func (t *insuranceParser) checkAgePoints(s *types.Survey, r types.SurveyResults) types.SurveyResults {

	if s.Age < 30 {
		r.VehiclePoints -= 2
		r.DisabilityPoints -= 2
		r.HomePoints -= 2
		r.LifePoints -= 2
	}

	if s.Age >= 30 && s.Age <= 40 {
		r.VehiclePoints -= 1
		r.DisabilityPoints -= 1
		r.HomePoints -= 1
		r.LifePoints -= 1
	}

	return r
}

func (t *insuranceParser) checkIncomePoints(s *types.Survey, r types.SurveyResults) types.SurveyResults {

	if s.Income > 200000 {
		r.VehiclePoints -= 1
		r.DisabilityPoints -= 1
		r.HomePoints -= 1
		r.LifePoints -= 1
	}

	return r
}

func (t *insuranceParser) checkHomePoints(s *types.Survey, r types.SurveyResults) types.SurveyResults {

	if s.House.OwnershipStatus == "mortgaged" {
		r.HomePoints += 1
		r.DisabilityPoints += 1
	}

	return r
}

func (t *insuranceParser) checkDependents(s *types.Survey, r types.SurveyResults) types.SurveyResults {

	if s.Dependents > 0 {
		r.DisabilityPoints += 1
		r.LifePoints += 1
	}

	return r
}

func (t *insuranceParser) checkMaritialPoints(s *types.Survey, r types.SurveyResults) types.SurveyResults {

	if s.MaritialStatus == "married" {
		r.LifePoints += 1
		r.DisabilityPoints -= 1
	}

	return r
}

func (t *insuranceParser) checkVehiclePoints(s *types.Survey, r types.SurveyResults) types.SurveyResults {

	if s.Vehicle.Year-time.Now().Year() <= 5 {
		r.VehiclePoints += 1
	}

	return r
}

func (t *insuranceParser) ParseSurvey(s *types.Survey) (types.SurveyResults, error) {
	var results types.SurveyResults

	riskPoints := t.sumRiskQuestions(s)

	results.DisabilityPoints = riskPoints
	results.HomePoints = riskPoints
	results.LifePoints = riskPoints
	results.VehiclePoints = riskPoints

	spew.Dump(results)

	results = t.verifyIneligibility(s, results)
	results = t.checkAgePoints(s, results)
	results = t.checkIncomePoints(s, results)
	results = t.checkHomePoints(s, results)
	results = t.checkDependents(s, results)
	results = t.checkMaritialPoints(s, results)
	results = t.checkVehiclePoints(s, results)

	return results, nil
}

func (t *insuranceParser) setDisabilityStatus(sr *types.SurveyResults) string {
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

	return assignment.Disability
}

func (t *insuranceParser) setHomeStatus(sr *types.SurveyResults) string {
	var assignment types.Assignment

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

	return assignment.Home
}

func (t *insuranceParser) setLifeStatus(sr *types.SurveyResults) string {
	var assignment types.Assignment

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

	return assignment.Life
}

func (t *insuranceParser) setVehicleStatus(sr *types.SurveyResults) string {
	var assignment types.Assignment

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

	return assignment.Vehicle
}
func (t *insuranceParser) SetAssignmentResults(sr *types.SurveyResults) types.Assignment {
	var assignment types.Assignment

	assignment.Disability = t.setDisabilityStatus(sr)
	assignment.Home = t.setHomeStatus(sr)
	assignment.Life = t.setLifeStatus(sr)
	assignment.Vehicle = t.setVehicleStatus(sr)

	return assignment
}
