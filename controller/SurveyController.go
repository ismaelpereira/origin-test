package controller

import (
	"encoding/json"
	"time"

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

func (t *insuranceParser) verifyIneligibility(s *types.Survey) types.SurveyResults {
	var results types.SurveyResults

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

	return results
}

func (t *insuranceParser) checkAgePoints(s *types.Survey) types.SurveyResults {
	var results types.SurveyResults

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

	return results
}

func (t *insuranceParser) checkIncomePoints(s *types.Survey) types.SurveyResults {
	var results types.SurveyResults

	if s.Income > 200000 {
		results.VehiclePoints -= 1
		results.DisabilityPoints -= 1
		results.HomePoints -= 1
		results.LifePoints -= 1
	}

	return results
}

func (t *insuranceParser) checkHomePoints(s *types.Survey) types.SurveyResults {
	var results types.SurveyResults

	if s.House.OwnershipStatus == "mortgaged" {
		results.HomePoints += 1
		results.DisabilityPoints += 1
	}

	return results
}

func (t *insuranceParser) checkDependents(s *types.Survey) types.SurveyResults {
	var results types.SurveyResults

	if s.Dependents > 0 {
		results.DisabilityPoints += 1
		results.LifePoints += 1
	}

	return results
}

func (t *insuranceParser) checkMaritialPoints(s *types.Survey) types.SurveyResults {
	var results types.SurveyResults

	if s.MaritialStatus == "married" {
		results.LifePoints += 1
		results.DisabilityPoints -= 1
	}

	return results
}

func (t *insuranceParser) checkVehiclePoints(s *types.Survey) types.SurveyResults {
	var results types.SurveyResults

	if s.Vehicle.Year-time.Now().Year() <= 5 {
		results.VehiclePoints += 1
	}

	return results
}

func (t *insuranceParser) ParseSurvey(s *types.Survey) (types.SurveyResults, error) {
	var results types.SurveyResults

	riskPoints := t.sumRiskQuestions(s)

	results.DisabilityPoints = riskPoints
	results.HomePoints = riskPoints
	results.LifePoints = riskPoints
	results.VehiclePoints = riskPoints

	results = t.verifyIneligibility(s)
	results = t.checkAgePoints(s)
	results = t.checkIncomePoints(s)
	results = t.checkHomePoints(s)
	results = t.checkDependents(s)
	results = t.checkMaritialPoints(s)
	results = t.checkVehiclePoints(s)

	return results, nil
}

func (t *insuranceParser) setDisabilityStatus(sr *types.SurveyResults) types.Assignment {
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

	return assignment
}

func (t *insuranceParser) setHomeStatus(sr *types.SurveyResults) types.Assignment {
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

	return assignment
}

func (t *insuranceParser) setLifeStatus(sr *types.SurveyResults) types.Assignment {
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

	return assignment
}

func (t *insuranceParser) setVehicleStatus(sr *types.SurveyResults) types.Assignment {
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

	return assignment
}
func (t *insuranceParser) SetAssignmentResults(sr *types.SurveyResults) types.Assignment {
	var assignment types.Assignment

	assignment = t.setDisabilityStatus(sr)
	assignment = t.setHomeStatus(sr)
	assignment = t.setLifeStatus(sr)
	assignment = t.setVehicleStatus(sr)

	return assignment
}
