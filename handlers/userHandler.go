package handlers

import (
	"breez/models"
	"breez/utils"
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

	if user.Name == "" || !isValidEmail(user.Email) || len(user.Password) < 8 {
		Log.Warn("Attempt to register user with missing name")
		http.Error(w, "Invalid input. Ensure all fields are correct.", http.StatusBadRequest)
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

	// Generate verification token
	verificationToken := utils.GenerateToken(user.Email)
	verificationURL := fmt.Sprintf("http://localhost:8080/api/verify?token=%s", verificationToken)

	// Send verification emai
	emailBody := fmt.Sprintf("Hello %s, please confirm your email by clicking the link: %s", user.Name, verificationURL)
	go utils.SendEmailWithAttachments(user.Email, "Email Verification", emailBody, nil)

	Log.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("User registered successfully and verification email sent")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully!"})
}

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	email := utils.VerifyToken(token)
	if email == "" {
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if user.IsVerified {
		http.Error(w, "Email already verified", http.StatusBadRequest)
		return
	}

	user.IsVerified = true
	db.Save(&user)

	Log.WithField("email", email).Info("Email verified successfully")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Email verified successfully!"})
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

	if !user.IsVerified {
		Log.WithField("email", user.Email).Warn("Unverified email login attempt")
		http.Error(w, "Email not verified. Please verify your email before logging in.", http.StatusForbidden)
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
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
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
		"user_id": userID,
		"name":    user.Name,
		"email":   user.Email,
		"role":    user.Role,
	})
}

func GetUserID(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"user_id": cookie.Value})
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
