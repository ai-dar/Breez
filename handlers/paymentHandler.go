package handlers

import (
	"breez/models"
	"breez/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

// HandlePayment обрабатывает запрос на оплату
func HandlePayment(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем `user_id` из cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Ошибка: пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		http.Error(w, "Ошибка: некорректный user_id", http.StatusBadRequest)
		return
	}

	// Получаем email из БД по `user_id`
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		http.Error(w, "Ошибка: пользователь не найден", http.StatusNotFound)
		return
	}

	userEmail := user.Email
	fmt.Println("Оплата для пользователя:", userID, userEmail)

	var req struct {
		Amount    float64 `json:"amount"`
		Currency  string  `json:"currency"`
		ServiceID string  `json:"service_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка парсинга запроса", http.StatusBadRequest)
		return
	}

	// Отправляем запрос в платежную систему
	paymentResp, err := utils.SendPaymentRequest(int(user.ID), req.Amount, req.Currency, req.ServiceID, user.Email, user.Name)
	if err != nil {
		http.Error(w, "Ошибка обработки платежа: "+err.Error(), http.StatusInternalServerError)
		return
	}

	redirectURL := fmt.Sprintf("http://localhost:8081/pay/%s", paymentResp.TransactionID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"redirect_url": redirectURL})
}
