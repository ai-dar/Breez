package handlers

import (
	"breez/models"
	"breez/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func handlePayment(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.URL.Query().Get("user_id"))
	amount, _ := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
	currency := r.URL.Query().Get("currency")
	serviceID := r.URL.Query().Get("service_id")
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		http.Error(w, "Ошибка: пользователь не найден", http.StatusBadRequest)
		return
	}
	fmt.Println("Email пользователя:", user.Email)

	paymentResp, err := utils.SendPaymentRequest(userID, amount, currency, serviceID, user.Email, user.Name)
	if err != nil {
		http.Error(w, "Ошибка платежа: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paymentResp)
}

// SetupRouter создает и настраивает маршруты
func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/pay", handlePayment).Methods("POST")

	// Маршруты для твитов
	r.HandleFunc("/tweets", CreateTweet).Methods("POST")
	r.HandleFunc("/tweets", GetTweetsWithFilters).Methods("GET")

	// Маршруты для админа
	r.HandleFunc("/admin/tweet/update", UpdateTweet).Methods("PUT")
	r.HandleFunc("/admin/tweet/delete", DeleteTweet).Methods("DELETE")

	// Маршруты для лайков
	r.HandleFunc("/like", LikeTweet).Methods("POST")

	// Добавьте здесь другие необходимые маршруты
	return r
}
