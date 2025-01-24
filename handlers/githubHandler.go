package handlers

import (
	"breez/models"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"gorm.io/gorm"
)

// Объект конфигурации OAuth2
var githubOauthConfig *oauth2.Config

func InitGitHubOAuth() {
	githubOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     github.Endpoint,
	}
}

// GitHubLoginHandler - перенаправляет пользователя на GitHub для авторизации
func GitHubLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := githubOauthConfig.AuthCodeURL("randomstate", oauth2.AccessTypeOffline)
	log.WithField("url", url).Info("Redirecting to GitHub login")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GitHubCallbackHandler - обрабатывает обратный вызов после авторизации через GitHub
func GitHubCallbackHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	// Получаем токен доступа
	token, err := githubOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	// Получаем информацию о пользователе
	client := githubOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	// Извлекаем email и имя пользователя
	email, _ := userInfo["email"].(string)
	name, _ := userInfo["name"].(string)

	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			user = models.User{Email: email, Name: name}
			if err := db.Create(&user).Error; err != nil {
				http.Error(w, "Failed to save user", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	}

	// Устанавливаем cookie с user_id и source
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    strconv.Itoa(int(user.ID)),
		HttpOnly: true,
		Path:     "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name:  "source",
		Value: "github",
		Path:  "/",
	})

	// Перенаправляем на главную страницу
	http.Redirect(w, r, "/static/home.html", http.StatusSeeOther)
}
