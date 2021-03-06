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

type SurveyResults struct {
	VehiclePoints    int
	DisabilityPoints int
	HomePoints       int
	LifePoints       int
}

type Assignment struct {
	Vehicle    string `json:"auto"`
	Disability string `json:"disability"`
	Home       string `json:"home"`
	Life       string `json:"life"`
}
