package handlers

import (
	"breez/models"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var db *gorm.DB
var Log = logrus.New()

func init() {
	Log.SetFormatter(&logrus.JSONFormatter{})
}

func InitDB(database *gorm.DB) {
	db = database
}

func isValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

// RegisterUser handles user registration
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		Log.Warn("Invalid Content-Type in request")
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var user models.User

	if user.Role != "admin" && user.Role != "user" {
		user.Role = "user" // По умолчанию обычный пользователь
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		Log.WithField("error", err).Warn("Failed to decode JSON payload")
		http.Error(w, "Invalid payload. Ensure JSON structure matches User model.", http.StatusBadRequest)
		return
	}

	if user.Name == "" {
		Log.Warn("Attempt to register user with missing name")
		http.Error(w, "Name is required.", http.StatusBadRequest)
		return
	}

	if !isValidEmail(user.Email) {
		Log.WithField("email", user.Email).Warn("Invalid email format")
		http.Error(w, "Invalid email format.", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		Log.WithField("error", err).Error("Failed to hash password")
		http.Error(w, "Failed to hash password.", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	if err := db.Create(&user).Error; err != nil {
		Log.WithField("error", err).Error("Failed to create user in database")
		http.Error(w, "Failed to create user. Ensure the email is unique.", http.StatusInternalServerError)
		return
	}

	Log.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("User registered successfully")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully!"})
}

// LoginUser handles user login
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials models.User
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		Log.WithField("error", err).Warn("Failed to decode login payload")
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.Where("email = ?", credentials.Email).First(&user).Error; err != nil {
		Log.WithField("email", credentials.Email).Warn("User not found during login")
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		Log.WithField("email", credentials.Email).Warn("Invalid password attempt")
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Set a cookie for the session
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    fmt.Sprintf("%d", user.ID),
		Path:     "/",
		HttpOnly: true,
	})

	logger.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("User logged in successfully")

	// Redirect to home.html
	http.Redirect(w, r, "/static/home.html", http.StatusSeeOther)
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Получаем cookie с user_id
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Ищем пользователя в базе данных
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Возвращаем информацию о пользователе
	json.NewEncoder(w).Encode(map[string]interface{}{
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}

// CheckAuth verifies if the user is authenticated
func CheckAuth(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		Log.Warn("Unauthorized access attempt - missing cookie")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var user models.User
	if err := db.First(&user, cookie.Value).Error; err != nil {
		Log.WithField("user_id", cookie.Value).Warn("Unauthorized access - invalid user ID")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	Log.WithField("user_id", user.ID).Info("User authentication successful")
	w.WriteHeader(http.StatusOK)
}
