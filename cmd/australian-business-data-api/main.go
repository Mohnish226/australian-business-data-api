package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mohnish226/australian-business-data-api/pkg/config"
	"github.com/mohnish226/australian-business-data-api/pkg/logger"
	"github.com/mohnish226/australian-business-data-api/pkg/services/api"
	"github.com/mohnish226/australian-business-data-api/pkg/services/cache"
	"github.com/mohnish226/australian-business-data-api/pkg/services/charts"
	"github.com/mohnish226/australian-business-data-api/pkg/services/output"
	"github.com/mohnish226/australian-business-data-api/pkg/services/similarity"
)

func main() {
	// Initialize logger
	if err := logger.Init("logs"); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	logger.Logger.Printf("Starting Australian Business Data API")

	flagCleanCache := flag.Bool("clean", false, "Clean the Expired cache")
	flagNoCache := flag.Bool("nocache", false, "Do not use cache")
	flagOutput := flag.String("output", "", "Output file")
	flagNoOutput := flag.Bool("no-output", false, "Do not output of the results")
	flagCacheExpiration := flag.Int("cache-expiration", 10, "Cache expiration time in minutes")

	flagSearchTerm := flag.String("search", "", "Search term")
	flagSearchDate := flag.String("date", "", "Search date")
	flagSearchState := flag.String("state", "", "Search state")
	flagSearchRegistrationStatus := flag.String("status", "", "Search registration status")

	flagSearchLike := flag.String("searchlike", "", "Search term for SQL LIKE query")

	flagGetAverageAge := flag.Bool("average-age", false, "Get average age of businesses")
	flagGetRegistrationStatusChart := flag.Bool("registration-chart", false, "Get registration status")
	flagGetRegistrationDistributionChart := flag.Bool("registration-distribution-chart", false, "Get registration distribution chart")
	flagGetRegistrationStateChart := flag.Bool("registration-state-chart", false, "Get registration state chart")

	flag.Parse()

	// Set cache expiration from flag
	config.CacheExpiration = time.Duration(*flagCacheExpiration) * time.Minute

	results := []map[string]interface{}{}
	var err error

	flag.Usage = func() {
		fmt.Println("Usage: australian-business-data-api [options]")
		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	if *flagCleanCache {
		logger.Logger.Printf("Cleaning expired cache")
		if err := cache.RemoveExpiredCache(); err != nil {
			logger.Logger.Printf("Failed to clean cache: %v", err)
			os.Exit(1)
		}
		logger.Logger.Printf("Cache cleaned successfully")
		os.Exit(0)
	}

	if *flagNoCache {
		logger.Logger.Printf("Cache disabled")
		// TODO: Implement cache disabling
	}

	if flagSearchTerm != nil && *flagSearchTerm != "" || flagSearchDate != nil && *flagSearchDate != "" {
		query := ""
		if flagSearchTerm != nil && *flagSearchTerm != "" {
			query = *flagSearchTerm
		} else {
			query = *flagSearchDate
		}

		filter := map[string]string{}

		if flagSearchState != nil && *flagSearchState != "" {
			*flagSearchState = strings.ToUpper(*flagSearchState)
			valid := false
			for _, validState := range config.ValidStates {
				if *flagSearchState == validState {
					valid = true
					break
				}
			}
			if !valid {
				fmt.Println("Invalid state. Valid values are:", strings.Join(config.ValidStates, ", "))
				return
			}
			filter["BN_STATE_OF_REG"] = *flagSearchState
		}

		if flagSearchRegistrationStatus != nil && *flagSearchRegistrationStatus != "" {
			*flagSearchRegistrationStatus = config.StatusAutoCorrect[strings.ToLower(*flagSearchRegistrationStatus)]
			if *flagSearchRegistrationStatus != "Registered" && *flagSearchRegistrationStatus != "Deregistered" {
				fmt.Println("Invalid registration status. Valid values are: Registered, Deregistered")
				return
			}
			filter["BN_STATUS"] = *flagSearchRegistrationStatus
		}

		apiService := api.NewService()
		results, err = apiService.BasicSearch(query, filter)
		if err != nil {
			logger.Logger.Printf("Search failed: %v", err)
			os.Exit(1)
		}
		logger.Logger.Printf("Found %d results", len(results))
	}

	if *flagSearchLike != "" {
		logger.Logger.Printf("Performing SQL search with term: %s", *flagSearchLike)
		apiService := api.NewService()
		results, err = apiService.SQLSearch(*flagSearchLike)
		if err != nil {
			logger.Logger.Printf("SQL search failed: %v", err)
			os.Exit(1)
		}
		logger.Logger.Printf("Found %d results", len(results))

		results = similarity.SortName(results, *flagSearchLike)
		config.Headers = append(config.Headers, "Match_Percent")
	}

	if len(results) == 0 {
		fmt.Println("No records found")
		return
	}

	if flagOutput != nil && *flagOutput != "" {
		if strings.HasSuffix(*flagOutput, ".csv") {
			err := output.CSVWriter(results, *flagOutput)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			err := output.TerminalTablePrint(results, *flagOutput)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	} else {
		if flagNoOutput != nil && *flagNoOutput {
			// Do nothing
		} else {
			output.TerminalTablePrint(results, "")
		}
	}

	if flagGetAverageAge != nil && *flagGetAverageAge {
		averageAge, err := charts.GetAverageAgeOfBusinesses(results)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Average age of businesses:", averageAge)
	}

	if flagGetRegistrationStatusChart != nil && *flagGetRegistrationStatusChart {
		charts.GetRegistrationStatusChart(results)
	}

	if flagGetRegistrationDistributionChart != nil && *flagGetRegistrationDistributionChart {
		charts.GetRegistrationDistributionChart(results)
	}

	if flagGetRegistrationStateChart != nil && *flagGetRegistrationStateChart {
		charts.GetRegistrationStateChart(results)
	}

	logger.Logger.Printf("Application completed successfully")
}
