package main

import (
	"fmt"
	"github.com/DanielFillol/Crawler_STJ/Crawler"
	"github.com/DanielFillol/Crawler_STJ/csv"
)

// main is the entry point of the program.
func main() {
	// Call SeleniumWebDriver to create a Selenium WebDriver instance and start the main timer
	driver, err := Crawler.SeleniumWebDriver()
	if err != nil {
		fmt.Println(err)
	}
	defer driver.Close()

	// Crawl data from the website
	crawElements, err := Crawler.Craw(driver)
	if err != nil {
		fmt.Println(err)
	}

	// Write the crawled data to a CSV file
	err = csv.WriteCSV(crawElements)
	if err != nil {
		fmt.Println(err)
	}
}
