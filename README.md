# Australian Business Data API

A command-line tool to search and analyze Australian business registration data from data.gov.au.

## Features

- Search businesses by name, date, state, and registration status
- Export results to CSV or formatted table output
- Generate business statistics and charts
- Caching support for improved performance
- Cross-platform support (Windows, macOS, Linux)
- Wildcard search support using '%' character

## Installation

### Prerequisites

- Go 1.21 or later
- Git

### Building from Source

1. Clone the repository:
```bash
git clone https://github.com/mohnish226/australian-business-data-api.git
cd australian-business-data-api
```

2. Build the application:
```bash
# Build for current platform
./build.sh

# Build for specific platform
./build.sh --linux    # Build for Linux
./build.sh --windows  # Build for Windows
./build.sh --mac      # Build for macOS

# Build for all platforms
./build.sh --all
```

The compiled binaries will be available in the `build` directory.

## Usage

### Basic Search

```bash
# Search by business name
./australian-business-data-api --search "ACME Corporation"

# Search by registration date
./australian-business-data-api --date "01/01/2020"

# Search by state
./australian-business-data-api --search "ACME" --state "NSW"

# Search by registration status
./australian-business-data-api --search "ACME" --status "Registered"

# Wildcard search examples
./australian-business-data-api --search "%ACME"    # Ends with ACME
./australian-business-data-api --search "ACME%"    # Starts with ACME
./australian-business-data-api --search "%ACME%"   # Contains ACME
```

### Output Options

```bash
# Export to CSV
./australian-business-data-api --search "ACME" --output "results.csv"

# Export to formatted table
./australian-business-data-api --search "ACME" --output "results.txt"

# Suppress output
./australian-business-data-api --search "ACME" --no-output
```

### Analysis Options

```bash
# Get average age of businesses
./australian-business-data-api --search "ACME" --average-age

# Generate registration status chart
./australian-business-data-api --search "ACME" --registration-chart

# Generate registration distribution chart
./australian-business-data-api --search "ACME" --registration-distribution-chart

# Generate registration state chart
./australian-business-data-api --search "ACME" --registration-state-chart
```

### Cache Management

```bash
# Clean expired cache
./australian-business-data-api --clean

# Disable cache
./australian-business-data-api --nocache

# Set cache expiration time (in minutes)
./australian-business-data-api --cache-expiration 30
```

## Output Fields

The tool provides the following information for each business:

- Business Name
- State of Registration
- Status
- Registration Date
- Cancellation Date

## Data Source

This tool uses the Australian Business Register data from [data.gov.au](https://data.gov.au/data/api/action/datastore_search?resource_id=55ad4b1c-5eeb-44ea-8b29-d410da431be3).

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Data provided by the Australian Government through data.gov.au
- Built with Go programming language 