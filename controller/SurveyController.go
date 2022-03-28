package controller

import (
	"encoding/json"
	"io/ioutil"
	"origin-challenge/types"
	"os"
	"time"
)

/**
Ler um JSON >
Calcular o risco para o seguro >
Retornar um json com a sugestão do tipo de seguro >



Linhas de seguro
- carmotiva
- Residencial
- Invalidez
- Vida

Verificações


if renda == 0 &&
carmóvel == {} &&
casa == 0 => invalidez = -1, carmovel = -1, casa = -1

if idade > 60 => invalidez = -1, vida = -1

if idade < 30 =>
riskQuestionsLife, riskQuestionsHouse, riskQuestionsInvalidez, riskQuestionsCar = -2

if idade > 30  && idade < 40 =>
riskQuestionsLife, riskQuestionsHouse, riskQuestionsInvalidez, riskQuestionsCar = -2

if renda >  200k =>
riskQuestionsLife, riskQuestionsHouse, riskQuestionsInvalidez, riskQuestionsLife = -1

if house == mortgaged => riskQuestionsHouse += 1 && riskQuestionsInvalidez += 1

if dependents > 0 => riskQuestionsInvalidez += 1 && riskQuestionsLife += 1


if married => riskQuestionsLife += 1 && riskQuestionsInvalidez -= 1

if veihcle.year - thisYear => 5 risQuestionsCar += 1


Resultados
 =< 0 - Economico
 1 ou 2 - Regular
 >= 3 - Responsável

*/

type insuranceParser struct {
	path string
}

type InsuranceParser interface {
	UnmarshallSurvey() (*types.Survey, error)
	MarshallAssignment(a types.Assignment) error
	ParseSurvey(survey types.Survey) (*types.SurveyResults, error)
	SetAssignmentResults(sr types.SurveyResults) (*types.Assignment, error)
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

func (t *insuranceParser) MarshallAssignment(a types.Assignment) error {

	return nil
}

func (t *insuranceParser) ParseSurvey(s types.Survey) (*types.SurveyResults, error) {
	var carRiskPoints, disabilityRiskPoints, homeRiskPoints, lifeRiskPoints *int

	survey, err := t.UnmarshallSurvey()

	if err != nil {
		return nil, err
	}

	if survey.Age > 60 {
		carRiskPoints = nil
		disabilityRiskPoints = nil
		homeRiskPoints = nil
	}

	if survey.Age < 30 {
		*carRiskPoints -= 2
		*disabilityRiskPoints -= 2
		*homeRiskPoints -= 2
		*lifeRiskPoints -= 2
	}

	if survey.Age > 30 && survey.Age < 40 || survey.Income > 200000 {
		*carRiskPoints -= 1
		*disabilityRiskPoints -= 1
		*homeRiskPoints -= 1
		*lifeRiskPoints -= 1
	}

	if survey.House.OwnershipStatus == "mortgaged" || survey.Dependents > 0 {
		*disabilityRiskPoints += 1
		*homeRiskPoints += 1
	}

	if survey.Dependents > 0 {
		*disabilityRiskPoints -= 1
		*lifeRiskPoints -= 1
	}

	if survey.MaritialStatus == "married" {
		*lifeRiskPoints += 1
		*disabilityRiskPoints -= 1
	}

	if survey.Vehicle.Year-time.Now().Year() >= 5 {
		*carRiskPoints += 1
	}

	var results types.SurveyResults

	results.Disability = *disabilityRiskPoints
	results.Home = *homeRiskPoints
	results.Life = *lifeRiskPoints
	results.Vehicle = *carRiskPoints

	return &results, nil
}

func (t *insuranceParser) SetAssignmentResults(sr types.SurveyResults) (*types.Assignment, error) {
	var assignment types.Assignment

	return &assignment, nil
}
