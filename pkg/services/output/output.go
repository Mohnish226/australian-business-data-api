package output

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/mohnish226/australian-business-data-api/pkg/config"
)

func CSVWriter(data []map[string]interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers using friendly names
	friendlyHeaders := make([]string, len(config.Headers))
	for i, header := range config.Headers {
		friendlyHeaders[i] = config.HeadersMap[header]
	}
	if err := writer.Write(friendlyHeaders); err != nil {
		return err
	}

	// Write data
	for _, record := range data {
		row := make([]string, len(config.Headers))
		for i, header := range config.Headers {
			if value, ok := record[header]; ok {
				row[i] = fmt.Sprintf("%v", value)
			}
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func TerminalTablePrint(data []map[string]interface{}, filename string) error {
	var writer *os.File
	var err error

	if filename != "" {
		writer, err = os.Create(filename)
		if err != nil {
			return err
		}
		defer writer.Close()
	} else {
		writer = os.Stdout
	}

	// Calculate column widths using friendly headers
	widths := make([]int, len(config.Headers))
	for i, header := range config.Headers {
		friendlyHeader := config.HeadersMap[header]
		widths[i] = len(friendlyHeader)
	}

	for _, record := range data {
		for i, header := range config.Headers {
			if value, ok := record[header]; ok {
				valueStr := fmt.Sprintf("%v", value)
				if len(valueStr) > widths[i] {
					widths[i] = len(valueStr)
				}
			}
		}
	}

	// Print headers using friendly names
	headerLine := ""
	for i, header := range config.Headers {
		if i > 0 {
			headerLine += " | "
		}
		friendlyHeader := config.HeadersMap[header]
		headerLine += fmt.Sprintf("%-*s", widths[i], friendlyHeader)
	}
	fmt.Fprintln(writer, headerLine)

	// Print separator
	separator := ""
	for i, width := range widths {
		if i > 0 {
			separator += "-+-"
		}
		separator += strings.Repeat("-", width)
	}
	fmt.Fprintln(writer, separator)

	// Print data
	for _, record := range data {
		line := ""
		for i, header := range config.Headers {
			if i > 0 {
				line += " | "
			}
			value := ""
			if v, ok := record[header]; ok {
				value = fmt.Sprintf("%v", v)
			}
			line += fmt.Sprintf("%-*s", widths[i], value)
		}
		fmt.Fprintln(writer, line)
	}

	return nil
}
