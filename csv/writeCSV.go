package csv

import (
	"encoding/csv"
	"github.com/DanielFillol/Crawler_STJ/Crawler"
	"os"
	"path/filepath"
	"strconv"
)

const FOLDERNAME = "Result"
const fileName = "thesis"

// WriteCSV writes the crawledElements data to a CSV file.
// It takes a slice of CrawElement as input and delegates the writing process to WriteData.
// If any error occurs during the writing process, it returns the error.
func WriteCSV(crawledElements []Crawler.CrawElement) error {
	// Delegate the writing process to WriteData
	err := WriteData(crawledElements)
	if err != nil {
		return err
	}
	return nil
}

// createFile creates a new file with the given path, including necessary directory creation.
// It takes the path as input and returns the created file or an error if any.
func createFile(p string) (*os.File, error) {
	// Create directories recursively if they don't exist
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	// Create and return the file
	return os.Create(p)
}

// WriteData writes the crawledElements data to a CSV file.
// It takes a slice of CrawElement as input, generates the CSV data, and writes it to a file.
// If any error occurs during the writing process, it returns the error.
func WriteData(crawledElements []Crawler.CrawElement) error {
	var rows [][]string

	// Generate and append header row
	rows = append(rows, generateHeaders())

	// Generate and append data rows for each CrawElement
	for _, lawsuit := range crawledElements {
		rows = append(rows, dataRows(lawsuit))
	}

	// Create a new file for writing
	cf, err := createFile(FOLDERNAME + "/" + fileName + ".csv")
	if err != nil {
		return err
	}

	// Initialize a CSV writer
	w := csv.NewWriter(cf)

	// Write all the rows to the CSV file
	err = w.WriteAll(rows)
	if err != nil {
		return err
	}

	return nil
}

// generateHeaders generates the header row for the CSV file.
func generateHeaders() []string {
	return []string{
		"Nome do Arquivo",
		"URL",
		"Download?",
	}
}

// dataRows generates data rows for a single CrawElement.
func dataRows(crawledElements Crawler.CrawElement) []string {
	return []string{
		crawledElements.FileName,
		crawledElements.URL,
		strconv.FormatBool(crawledElements.Download),
	}
}
