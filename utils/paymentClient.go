package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// PaymentRequest структура запроса на оплату
type PaymentRequest struct {
	UserID    int     `json:"user_id"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	ServiceID string  `json:"service_id"`
	Email     string  `json:"email"`
	Name      string  `json:"name"`
}

// PaymentResponse структура ответа от платежного сервиса
type PaymentResponse struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	TransactionID string `json:"transaction_id"` // <-- Добавляем это поле
}

// SendPaymentRequest отправляет запрос в Breez Payment System (порт 8081)
func SendPaymentRequest(userID int, amount float64, currency, serviceID, email string, name string) (*PaymentResponse, error) {
	paymentURL := "http://localhost:8081/api/pay"

	requestBody, err := json.Marshal(PaymentRequest{
		UserID:    userID,
		Amount:    amount,
		Currency:  currency,
		ServiceID: serviceID,
		Email:     email,
		Name:      name,
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("Отправляем запрос на оплату с email:", email)

	resp, err := http.Post(paymentURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Логируем ответ сервера
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Ошибка чтения ответа от сервера: " + err.Error())
	}
	fmt.Println("Ответ от платежного сервиса:", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Ошибка при создании платежа: " + string(body))
	}

	var paymentResp PaymentResponse
	if err := json.Unmarshal(body, &paymentResp); err != nil {
		return nil, err
	}

	return &paymentResp, nil
}
