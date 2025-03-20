package config

import (
	"os"
	"path/filepath"
	"time"
)

const (
	Host                   = "https://data.gov.au"
	RestPath               = "/data/api/action/datastore_search"
	SQLPath                = "/data/api/action/datastore_search_sql"
	ResourceID             = "55ad4b1c-5eeb-44ea-8b29-d410da431be3"
	RequestLimit           = 50
	APIToken               = ""
	DefaultCacheExpiration = time.Minute * 10 // 10 minutes
)

var (
	// APIKey is the key for the Australian Business Data API
	APIKey = os.Getenv("ABN_API_KEY")

	// CacheDir is the directory where cache files are stored
	CacheDir = filepath.Join(os.TempDir(), "abn-cache")

	// CacheExpiration is the duration for which cached data remains valid
	CacheExpiration = 24 * time.Hour

	// APIBaseURL is the base URL for the Australian Business Data API
	APIBaseURL = "https://abr.business.gov.au/json/AbnDetails.aspx"

	// MaxResults is the maximum number of results to return per request
	MaxResults = 100
)

var Headers = []string{
	"BN_NAME",
	"BN_STATE_OF_REG",
	"BN_STATUS",
	"BN_REG_DT",
	"BN_CANCEL_DT",
}

var HeadersMap = map[string]string{
	"BN_NAME":         "Business Name",
	"BN_STATE_OF_REG": "State of Registration",
	"BN_STATUS":       "Status",
	"BN_REG_DT":       "Registration Date",
	"BN_CANCEL_DT":    "Cancellation Date",
}

var ValidStates = []string{
	"ACT",
	"NSW",
	"NT",
	"QLD",
	"SA",
	"TAS",
	"VIC",
	"WA",
}

var StatusAutoCorrect = map[string]string{
	"register":     "Registered",
	"registered":   "Registered",
	"reg":          "Registered",
	"de-register":  "Deregistered",
	"deregister":   "Deregistered",
	"deregistered": "Deregistered",
	"dereg":        "Deregistered",
	"de-reg":       "Deregistered",
	"cancel":       "Deregistered",
	"cancelled":    "Deregistered",
	"canceled":     "Deregistered",
}
