package types

type Survey struct {
	Age            int
	Dependents     int
	House          House
	Income         int
	MaritialStatus string `json:"maritial_status"`
	RiskQuestions  []int  `json:"risk_questions"`
	Vehicle        Vehicle
}

type House struct {
	OwnershipStatus string `json:"ownership_status"`
}

type Vehicle struct {
	Year int
}

type Assignment struct {
	Vehicle    string `json:"auto"`
	Disability string
	Home       string
	Life       string
}
