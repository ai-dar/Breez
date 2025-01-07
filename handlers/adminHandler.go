package handlers

import (
	"encoding/json"
	"net/http"

	"breez/models"
	"breez/utils"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var log = logrus.New()

func RegisterAdmin(w http.ResponseWriter, r *http.Request) {
	var user models.User // Используем общую модель User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// Проверяем, чтобы роль была либо 'admin', либо 'user'
	if user.Role != "admin" && user.Role != "user" {
		user.Role = "admin" // Если роль не указана или некорректна, задаем 'admin'
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Сохраняем администратора в базе данных
	if err := db.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create admin", http.StatusInternalServerError)
		return
	}

	log.WithField("email", user.Email).Info("Admin registered successfully")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Admin registered successfully!"})
}

func SendEmailToAllUsers(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.WithField("error", err).Warn("Invalid email payload")
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// Получение всех пользователей
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		log.WithField("error", err).Error("Failed to fetch users")
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	// Рассылка email пользователям
	for _, user := range users {
		if err := utils.SendEmail(user.Email, request.Subject, request.Body); err != nil {
			log.WithFields(logrus.Fields{
				"email": user.Email,
				"error": err,
			}).Error("Failed to send email")
		} else {
			log.WithField("email", user.Email).Info("Email sent successfully")
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Emails sent successfully!"})
}
