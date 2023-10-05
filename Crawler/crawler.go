package Crawler

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/tebeka/selenium"
	"golang.org/x/net/html"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	InitWebSite     = "https://scon.stj.jus.br/SCON/jt/toc.jsp?tipo=JT&b=TEMA&p=true&thesaurus=JURIDICO&l=20&i=1&operador=E&ordenacao=MAT,@NUM"
	XpathTableLine  = "//*[@id=\"corpopaginajurisprudencia\"]/div[3]/div[2]/div/div/form/table/tbody/tr"
	XpathDocLink    = "td[2]/span/a"
	XpathPageNumber = "/html/body/div[1]/section[2]/div[2]/div[3]/div[2]/div/div/div[2]/div/span[3]"
	BASEURL         = "https://www.stj.jus.br/docs_internet/jurisprudencia/jurisprudenciaemteses/Jurisprudencia%20em%20Teses%20"
	EXTENSION       = ".pdf"
)

type CrawElement struct {
	FileName string
	URL      string
	Download bool
}

// Craw scrapes data from a website using a Selenium WebDriver and returns a list of CrawElement.
// It takes the WebDriver as input, navigates through multiple pages, extracts data, and downloads files.
// It returns a list of CrawElement containing file information and download status, along with any errors.
func Craw(driver selenium.WebDriver) ([]CrawElement, error) {
	// Initialize the search link
	searchLink := InitWebSite

	// Navigate to the initial search link
	err := driver.Get(searchLink)
	if err != nil {
		return nil, errors.New("failed to get the search link, err: " + err.Error())
	}

	// Get the HTML page source
	htmlPgSrc, err := getPageSource(driver)
	if err != nil {
		return nil, err
	}

	// Extract the final page number
	finalPage, err := totalPages(htmlquery.InnerText(htmlquery.FindOne(htmlPgSrc, XpathPageNumber)))
	if err != nil {
		return nil, err
	}

	// Initialize the list to store CrawElements
	var crawList []CrawElement

	// Initialize the magic number
	magicNumber := 2

	// Loop through pages
	for i := 1; i <= finalPage; i++ {
		// Navigate to the next page and update the magic number
		driver2, ma, err := nextPage(driver, i, magicNumber)
		if err != nil {
			return nil, err
		}
		magicNumber = ma

		// Get the HTML page source of the new page
		htmlPgSrc2, err := getPageSource(driver2)
		if err != nil {
			return nil, err
		}

		// Find table lines in the HTML
		tableLines := htmlquery.Find(htmlPgSrc2, XpathTableLine)

		// Loop through table lines
		for _, line := range tableLines {
			// Check if a document link exists
			exist := htmlquery.Find(line, XpathDocLink)
			if len(exist) > 0 {
				// Extract the edition number and title
				number, title := extractEditionAndTitle(htmlquery.InnerText(exist[0]))

				// Append the CrawElement to the list
				crawList = append(crawList, CrawElement{
					FileName: title,
					URL:      parseToURL(BASEURL, number, title, EXTENSION),
				})
			}
		}

		driver2.Close()
	}

	// Loop through the list of CrawElements and download files
	for _, s := range crawList {
		err = download(s.URL, s.FileName)
		if err != nil {
			s.Download = false
		} else {
			s.Download = true
		}
	}

	// Return the list of CrawElements
	return crawList, nil
}

// getPageSource retrieves the page source from a Selenium WebDriver and parses it into an HTML node.
// It takes a WebDriver as input and returns the parsed HTML node.
// If there is an error while getting the page source or parsing it, it returns an error.
func getPageSource(driver selenium.WebDriver) (*html.Node, error) {
	// Get the page source from the WebDriver
	pageSource, err := driver.PageSource()
	if err != nil {
		return nil, errors.New("could not get page source, err: " + err.Error())
	}

	// Parse the page source into an HTML node
	htmlPgSrc, err := htmlquery.Parse(strings.NewReader(pageSource))
	if err != nil {
		return nil, errors.New("could not convert string to HTML node, err: " + err.Error())
	}

	// Return the parsed HTML node
	return htmlPgSrc, nil
}

// nextPage closes the old Selenium WebDriver, creates a new one, and navigates to the next page based on the given parameters.
// It takes the old WebDriver, page index (i), and magic number (magic) as input.
// If i is 1 or 2, it navigates to the first or second page respectively; otherwise, it calculates the next page using the magic number.
// It returns the new WebDriver, updated page index, and any error that may occur during the process.
func nextPage(driverOld selenium.WebDriver, i int, magic int) (selenium.WebDriver, int, error) {
	// Close the old WebDriver
	_ = driverOld.Close()

	// Create a new WebDriver
	driver, err := SeleniumWebDriver()
	if err != nil {
		return nil, 2, errors.New("error creating new driver on next page, err:" + err.Error())
	}

	// Navigate to the next page based on the page index (i)
	if i == 1 || i == 2 {
		if i == 1 {
			// Navigate to the first page
			err = driver.Get("https://scon.stj.jus.br/SCON/jt/toc.jsp?tipo=JT&b=TEMA&p=true&thesaurus=JURIDICO&l=20&i=1&operador=E&ordenacao=MAT,@NUM")
			if err != nil {
				return nil, i, errors.New("failed to get the search link, err: " + err.Error())
			}
		} else {
			// Navigate to the second page
			err = driver.Get("https://scon.stj.jus.br/SCON/jt/toc.jsp?tipo=JT&b=TEMA&p=true&thesaurus=JURIDICO&l=20&i=21&operador=E&ordenacao=MAT,@NUM")
			if err != nil {
				return nil, i, errors.New("failed to get the search link, err: " + err.Error())
			}
		}

		return driver, i, nil
	} else {
		// Calculate the next page index using the magic number (n)
		n := newMagic(magic)
		err = driver.Get("https://scon.stj.jus.br/SCON/jt/toc.jsp?tipo=JT&b=TEMA&p=true&thesaurus=JURIDICO&l=20&i=" + strconv.Itoa(n) + "1&operador=E&ordenacao=MAT,@NUM")
		if err != nil {
			return nil, i, errors.New("failed to get the search link, err: " + err.Error())
		}
		return driver, n, nil
	}
}

// totalPages extracts the total number of pages (Y) from a string in the format "Page X of Y".
// It takes the input text and returns the total number of pages as an integer.
// If the input format is not as expected or if there are conversion errors, it returns an error.
func totalPages(textPageNumber string) (int, error) {
	// Split the string into separate words using space as a delimiter
	words := strings.Fields(textPageNumber)

	// Check if the string has at least three words (Page X of Y)
	if len(words) < 3 {
		return 0, fmt.Errorf("the input does not have the expected format")
	}

	// Extract the first number (X)
	_, err := strconv.Atoi(words[1])
	if err != nil {
		return 0, errors.New("error converting string to int (Initial Page), err: " + err.Error())
	}

	// Extract the second number (Y)
	secondNumber, err := strconv.Atoi(words[3])
	if err != nil {
		return 0, errors.New("error converting string to int (Final Page), err: " + err.Error())
	}

	// Return the extracted total number of pages (Y)
	return secondNumber, nil
}

// parseToURL constructs a URL by combining the base URL, file number, cleaned file name, and file extension.
// It takes the base URL, file number, file name, and file extension as input and returns the constructed URL as a string.
func parseToURL(baseURL string, fileNumber string, fileName string, fileExtension string) string {
	return baseURL + fileNumber + "%20-%20" + cleanString(removeSpecialCharacters(capitalizeWords(strings.TrimSpace(fileName)))) + fileExtension
}

// cleanString removes special characters and replaces spaces with "%20" in a given string.
// It takes a string as input and returns a cleaned string.
func cleanString(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, " ", "%20"), "/", "%20"), "N.", "N"), ":", ""), ".", "%20"), "(", ""), ")", "")
}

// newMagic increments the value of oldMagic by 2.
// It takes an integer (oldMagic) as input and returns the updated integer.
// I use this to conform with TJs logic on url code
func newMagic(oldMagic int) int {
	return oldMagic + 2
}

// extractEditionAndTitle extracts the edition number and title from a given input text using regular expressions.
// It expects the input text to be in the format "EDIÇÃO N. {edition_number}: {title}" and returns both values.
// If no match is found, it returns empty strings for both edition number and title.
func extractEditionAndTitle(text string) (string, string) {
	// Define a regular expression pattern to match the edition number and title.
	pattern := `EDIÇÃO N\. (\d+): (.+)`

	// Compile the regular expression.
	regex := regexp.MustCompile(pattern)

	// Find the submatches in the input text.
	matches := regex.FindStringSubmatch(text)

	// Check if we have at least two submatches (edition number and title).
	if len(matches) >= 3 {
		editionNumber := matches[1]
		title := matches[2]
		return editionNumber, title
	}

	// Return empty strings if no match was found.
	return "", ""
}

// capitalizeWords capitalizes words in a given input string while handling exceptions for prepositions and numbers.
func capitalizeWords(input string) string {
	// Split the input into individual words
	words := strings.Fields(input)

	// Create a slice to store the capitalized words
	capitalizedWords := make([]string, len(words))

	// Define a map of prepositions for which the first letter should not be capitalized
	prepositions := map[string]struct{}{
		"do": {}, "da": {}, "de": {}, "e": {}, "das": {},
		// Add more prepositions as needed
	}

	// Define a map of numbers that should be fully capitalized
	numbers := map[string]struct{}{
		"i": {}, "ii": {}, "iii": {}, "iv": {}, "v": {}, "vi": {}, "vii": {}, "viii": {}, "ix": {}, "x": {},
		// Add more numbers as needed
	}

	// Loop through the words in the input string
	for i, word := range words {
		if len(word) > 0 {
			// Check if the word is a preposition
			_, isPreposition := prepositions[strings.ToLower(word)]
			// Check if the word is a number
			_, isNumbers := numbers[strings.ToLower(word)]

			if !isPreposition && !isNumbers {
				// Capitalize the first letter and make the rest of the word lowercase
				firstLetter := strings.ToUpper(word[0:1])
				restOfWord := word[1:]
				capitalizedWords[i] = firstLetter + strings.ToLower(restOfWord)
			} else if isNumbers {
				// Capitalize the entire word if it is a number
				capitalizedWords[i] = strings.ToUpper(word)
			} else {
				// Handle prepositions and exceptions for the first word
				if i == 0 && isPreposition {
					firstLetter := strings.ToUpper(word[0:1])
					restOfWord := word[1:]
					capitalizedWords[i] = firstLetter + strings.ToLower(restOfWord)
				} else {
					// For other prepositions, make them lowercase
					capitalizedWords[i] = strings.ToLower(word)
				}
			}
		}
	}

	// Join the capitalized words back into a single string with spaces
	return strings.Join(capitalizedWords, " ")
}

// removeSpecialCharacters replaces special characters in a given text with their simple counterparts.
// It uses a predefined map to perform character substitutions.
func removeSpecialCharacters(text string) string {
	// Mapping special characters to their simple counterparts
	charMap := map[string]string{
		"á": "a", "à": "a", "ã": "a", "â": "a",
		"é": "e", "è": "e", "ê": "e",
		"í": "i", "ì": "i", "î": "i",
		"ó": "o", "ò": "o", "õ": "o", "ô": "o",
		"ú": "u", "ù": "u", "û": "u",
		"ç": "c",
	}

	// Replace special characters with their simple counterparts
	for special, simple := range charMap {
		text = strings.ReplaceAll(text, special, simple)
	}

	// Return the modified text with special characters replaced
	return text
}

// download a file from the specified URL and saves it with the given fileName.
// It returns an error if any issues occur during the download.
func download(url string, fileName string) error {
	// Open a new WebDriver
	driver, err := SeleniumWebDriver()
	if err != nil {
		return errors.New("error creating WebDriver on download, err: " + err.Error())
	}
	defer driver.Close()

	// Send an HTTP GET request to the URL
	err = driver.Get(url)
	if err != nil {
		return errors.New("error getting the URL on download, err:" + err.Error())
	}

	// Wait for the PDF page to load (adjust wait time as needed)
	time.Sleep(1 * time.Second)

	// Simulate keyboard shortcuts to save the page as PDF
	_, err = driver.ExecuteScript(`window.print();`, nil)
	if err != nil {
		log.Fatalf("Failed to send printToPDF command: %v", err)
	}

	// Wait for a few seconds to allow the print dialog to appear
	time.Sleep(5 * time.Second)

	// Simulate user interaction to select the file save location and confirm
	err = savePDFUsingKeyboard()
	if err != nil {
		return err
	}

	// Print a message to indicate successful download
	fmt.Printf("Downloaded: %s\n", fileName)

	return nil
}

// savePDFUsingKeyboard simulates keyboard interaction to save the PDF
func savePDFUsingKeyboard() error {
	// Simulate keyboard shortcuts for saving the PDF (platform-dependent)
	// These shortcuts might vary based on the operating system and browser version
	// Below is an example for macOS with Firefox:
	cmd := exec.Command("osascript", "-e", `
		tell application "System Events"
			keystroke "s" using {command down}
			delay 1
			keystroke return
		end tell
	`)
	err := cmd.Run()
	if err != nil {
		return errors.New("Failed to simulate keyboard interaction: " + err.Error())
	}
	return nil
}
