package Crawler

import (
	"errors"
	"github.com/tebeka/selenium"
)

// SeleniumWebDriver creates a new headless Firefox Selenium WebDriver and returns it.
// It sets the necessary capabilities and initializes the WebDriver.
// If any error occurs during the WebDriver creation, it returns an error.
func SeleniumWebDriver() (selenium.WebDriver, error) {
	// Set the desired capabilities for the Firefox browser in headless mode
	caps := selenium.Capabilities(map[string]interface{}{"browserName": "firefox", "Args": "--headless"})

	// Create a new remote WebDriver instance
	driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		return nil, errors.New("could not create WebDriver")
	}

	// Optionally, you can resize the window if needed
	// driver.ResizeWindow("", 0, 0)

	// Return the initialized WebDriver
	return driver, nil
}
