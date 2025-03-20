package charts

import (
	"fmt"
	"time"
)

func GetAverageAgeOfBusinesses(data []map[string]interface{}) (float64, error) {
	var totalAge float64
	var count int

	for _, record := range data {
		regDate, ok := record["BN_REG_DT"].(string)
		if !ok {
			continue
		}

		parsedDate, err := time.Parse("02/01/2006", regDate)
		if err != nil {
			continue
		}

		age := time.Since(parsedDate).Hours() / 24 / 365.25
		totalAge += age
		count++
	}

	if count == 0 {
		return 0, fmt.Errorf("no valid registration dates found")
	}

	return totalAge / float64(count), nil
}

func GetRegistrationStatusChart(data []map[string]interface{}) {
	statusCount := make(map[string]int)

	for _, record := range data {
		status, ok := record["BN_STATUS"].(string)
		if !ok {
			continue
		}
		statusCount[status]++
	}

	fmt.Println("\nRegistration Status Distribution:")
	fmt.Println("--------------------------------")
	for status, count := range statusCount {
		fmt.Printf("%s: %d\n", status, count)
	}
}

func GetRegistrationDistributionChart(data []map[string]interface{}) {
	yearCount := make(map[int]int)

	for _, record := range data {
		regDate, ok := record["BN_REG_DT"].(string)
		if !ok {
			continue
		}

		parsedDate, err := time.Parse("02/01/2006", regDate)
		if err != nil {
			continue
		}

		yearCount[parsedDate.Year()]++
	}

	fmt.Println("\nRegistration Distribution by Year:")
	fmt.Println("--------------------------------")
	for year := 2000; year <= time.Now().Year(); year++ {
		if count, exists := yearCount[year]; exists {
			fmt.Printf("%d: %d\n", year, count)
		}
	}
}

func GetRegistrationStateChart(data []map[string]interface{}) {
	stateCount := make(map[string]int)

	for _, record := range data {
		state, ok := record["BN_STATE_OF_REG"].(string)
		if !ok {
			continue
		}
		stateCount[state]++
	}

	fmt.Println("\nRegistration Distribution by State:")
	fmt.Println("--------------------------------")
	for state, count := range stateCount {
		fmt.Printf("%s: %d\n", state, count)
	}
}
