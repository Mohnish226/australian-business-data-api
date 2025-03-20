package models

// Business represents a business entity from the API
type Business struct {
	Name           string `json:"BN_NAME"`
	StateOfReg     string `json:"BN_STATE_OF_REG"`
	Status         string `json:"BN_STATUS"`
	RegistrationDt string `json:"BN_REG_DT"`
	CancellationDt string `json:"BN_CANCEL_DT"`
	ABN            string `json:"ABN"`
	EntityType     string `json:"ENTITY_TYPE"`
	EntityStatus   string `json:"ENTITY_STATUS"`
	EntityTypeCode string `json:"ENTITY_TYPE_CODE"`
	Postcode       string `json:"POSTCODE"`
}

// BusinessStats represents business statistics
type BusinessStats struct {
	TotalCount           int
	ActiveCount          int
	AverageAge           float64
	StatusDistribution   map[string]int
	StateDistribution    map[string]int
	RegistrationByMonth  map[string]int
	RegistrationByYear   map[string]int
	RegistrationByState  map[string]int
	RegistrationByStatus map[string]int
}
