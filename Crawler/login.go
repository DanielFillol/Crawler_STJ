package Crawler

import (
	"errors"
	"github.com/tebeka/selenium"
)

const (
	SiteLogin  = "https://cpe.web.stj.jus.br/#/"
	XpathUser  = "//*[@id=\"cpf\"]"
	XpathPass  = "//*[@id=\"app\"]/div/div[2]/div[2]/div[2]/form/div[2]/input"
	XpathBtt   = "//*[@id=\"app\"]/div/div[2]/div[2]/div[2]/form/div[3]/div[1]/button"
	XpathError = "//*[@id=\"app\"]/div/div[4]/div/div/div[2]"
)

// Login performs the login operation on a website using a Selenium WebDriver.
// It takes the WebDriver, login, and password as input parameters and returns an error if any step fails.
// It navigates to the login page, enters the login and password, clicks the login button, and checks for login errors.
func Login(driver selenium.WebDriver, login string, password string) error {
	// Navigate to the login page
	err := driver.Get(SiteLogin)
	if err != nil {
		return errors.New("URL unavailable")
	}

	// Find the elements for username, password, and login button using XPaths
	userName, err := driver.FindElement(selenium.ByXPATH, XpathUser)
	if err != nil {
		return errors.New("XPath for username not found")
	}

	psw, err := driver.FindElement(selenium.ByXPATH, XpathPass)
	if err != nil {
		return errors.New("XPath for password not found")
	}

	btt, err := driver.FindElement(selenium.ByXPATH, XpathBtt)
	if err != nil {
		return errors.New("XPath for login button not found")
	}

	// Enter the login and password
	err = userName.SendKeys(login)
	if err != nil {
		return errors.New("Could not send login parameter")
	}

	err = psw.SendKeys(password)
	if err != nil {
		return errors.New("Could not send password parameter")
	}

	// Click on the login button
	err = btt.Click()
	if err != nil {
		return errors.New("Could not click on the login button")
	}

	// Check for login errors
	infoLogin, err := driver.FindElements(selenium.ByXPATH, XpathError)
	if err != nil {
		return errors.New("Could not find XPath for login error")
	}

	// If there are login errors, check the error messages
	if len(infoLogin) > 0 {
		for _, info := range infoLogin {
			innerText, err := info.Text()
			if err != nil {
				return errors.New("Could not find inner text message")
			}
			if innerText == "Usuário ou senha inválidos." {
				return errors.New("Wrong user or password parameter")
			}
		}
	}

	return nil
}
