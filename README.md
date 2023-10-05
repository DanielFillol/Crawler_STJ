# Crawler_STJ
This Go program demonstrates how to scrape data from a website using Selenium WebDriver. It navigates through multiple pages, extracts data, and writes it to a CSV file. The program is designed to be a useful tool for web scraping tasks where dynamic content loading or interactions with web elements are required.

## Features
- Utilizes Selenium WebDriver to automate web interactions.
- Supports headless browser mode for web scraping.
- Navigates through paginated content.
- Extracts data from HTML elements.
- Writes scraped data to a CSV file.

## Prerequisites
- [Go](https://go.dev) (1.17 or higher)
- [Selenium WebDriver](https://www.selenium.dev/downloads/) (3.141.0 or compatible)
- Firefox WebDriver for Selenium (GeckoDriver)

## Installation
1. Install Go if you haven't already: [Download and Install Go](https://go.dev/dl/).
1. Install Selenium WebDriver:
   ```go
    go get github.com/tebeka/selenium
    ```
1. Install Firefox WebDriver (GeckoDriver):
    - Download the appropriate version of GeckoDriver for your system from the [official website](https://github.com/mozilla/geckodriver/releases).
    - Place the downloaded executable in a directory included in your system's PATH.


## Usage
1. Clone this repository:
    ```go
    git clone https://github.com/DanielFillol/Crawler_STJ
    ```
1. Update the configuration and URLs in the **'config.go'** file to match your specific web scraping requirements.
1. Run the selenium server
   ```java
   java -jar selenium-server-standalone.jar
    ```
1. Build and run the program:
   ```go
   go run main.go
   ```
1. The program will navigate through the specified website, scrape data, and write it to a CSV file.

## Configuration
- **'config.go'**: Contains configuration constants and URLs used by the web scraper. Modify this file to adapt the program to your scraping needs.
Contributing

Contributions are welcome!
