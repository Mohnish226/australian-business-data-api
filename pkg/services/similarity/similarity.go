package similarity

import (
	"fmt"
	"math"
	"strings"
)

func calculateSimilarity(str1, str2 string) float64 {
	str1 = strings.ToLower(str1)
	str2 = strings.ToLower(str2)

	if str1 == str2 {
		return 100.0
	}

	words1 := strings.Fields(str1)
	words2 := strings.Fields(str2)

	if len(words1) == 0 || len(words2) == 0 {
		return 0.0
	}

	// Calculate word overlap
	overlap := 0
	for _, word1 := range words1 {
		for _, word2 := range words2 {
			if word1 == word2 {
				overlap++
				break
			}
		}
	}

	// Calculate similarity percentage
	maxWords := float64(math.Max(float64(len(words1)), float64(len(words2))))
	if maxWords == 0 {
		return 0.0
	}

	return (float64(overlap) / maxWords) * 100.0
}

func SortName(data []map[string]interface{}, searchTerm string) []map[string]interface{} {
	type recordWithSimilarity struct {
		record     map[string]interface{}
		similarity float64
	}

	records := make([]recordWithSimilarity, len(data))
	for i, record := range data {
		name, ok := record["BN_NAME"].(string)
		if !ok {
			name = ""
		}
		records[i] = recordWithSimilarity{
			record:     record,
			similarity: calculateSimilarity(name, searchTerm),
		}
	}

	// Sort records by similarity in descending order
	for i := 0; i < len(records)-1; i++ {
		for j := i + 1; j < len(records); j++ {
			if records[i].similarity < records[j].similarity {
				records[i], records[j] = records[j], records[i]
			}
		}
	}

	// Add similarity percentage to each record
	result := make([]map[string]interface{}, len(records))
	for i, record := range records {
		record.record["Match_Percent"] = fmt.Sprintf("%.2f%%", record.similarity)
		result[i] = record.record
	}

	return result
}
