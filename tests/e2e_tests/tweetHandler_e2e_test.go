package e2e_tests

import (
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

func TestIndexPageE2E(t *testing.T) {
	const (
		seleniumPath     = "selenium-server-4.28.0.jar"
		chromeDriverPath = "chromedriver"
		port             = 4444
	)

	// Запуск Selenium Server
	service, err := selenium.NewSeleniumService(seleniumPath, port)
	if err != nil {
		t.Fatalf("Failed to start Selenium server: %s", err)
	}
	defer service.Stop()

	// Подключение к WebDriver
	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, "http://localhost:4444/wd/hub")
	if err != nil {
		t.Fatalf("Failed to connect to WebDriver: %s", err)
	}
	defer wd.Quit()

	// Открываем главную страницу
	err = wd.Get("http://localhost:8080/static/index.html")
	if err != nil {
		t.Fatalf("Failed to load page: %s", err)
	}
	t.Log("Main page loaded successfully.")

	// Проверяем наличие заголовка "Breez"
	header, err := wd.FindElement(selenium.ByTagName, "h1")
	if err != nil {
		t.Fatalf("Header not found: %s", err)
	}
	headerText, err := header.Text()
	if err != nil {
		t.Fatalf("Failed to get header text: %s", err)
	}
	if headerText != "Breez" {
		t.Fatalf("Unexpected header text. Expected: 'Breez', Got: '%s'", headerText)
	}
	t.Logf("Header verified: %s", headerText)

	// Проверяем форму регистрации
	registerForm, err := wd.FindElement(selenium.ByID, "registerForm")
	if err != nil {
		t.Fatalf("Registration form not found: %s", err)
	}
	t.Log("Registration form found.")

	// Вводим данные в форму регистрации
	nameInput, _ := registerForm.FindElement(selenium.ByID, "name")
	emailInput, _ := registerForm.FindElement(selenium.ByID, "email")
	passwordInput, _ := registerForm.FindElement(selenium.ByID, "password")
	registerButton, _ := registerForm.FindElement(selenium.ByTagName, "button")

	nameInput.SendKeys("Test User")
	emailInput.SendKeys("testuser@example.com")
	passwordInput.SendKeys("password123")
	t.Log("Registration form filled successfully.")

	// Кликаем по кнопке регистрации
	registerButton.Click()
	t.Log("Registration button clicked.")

	// Ожидание результата
	time.Sleep(2 * time.Second)

	// Проверяем, что результат или сообщение об ошибке появилось
	alertText, err := wd.AlertText()
	if err == nil {
		t.Logf("Registration result: %s", alertText)
		wd.AcceptAlert()
	} else {
		t.Log("No alert displayed. Registration may not have completed as expected.")
	}

	// Проверяем кнопку "Sign in with GitHub"
	githubButton, err := wd.FindElement(selenium.ByTagName, "button")
	if err != nil {
		t.Fatalf("GitHub login button not found: %s", err)
	}
	githubButtonText, _ := githubButton.Text()
	if githubButtonText != "Sign in with GitHub" {
		t.Fatalf("Unexpected GitHub button text. Expected: 'Sign in with GitHub', Got: '%s'", githubButtonText)
	}
	t.Log("GitHub login button verified.")
}
