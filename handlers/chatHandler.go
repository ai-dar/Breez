package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ChatResponse struct {
	ChatID string `json:"chat_id"`
}

// Проверка активного чата
func CheckActiveChat(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	var chat struct {
		ChatID string `json:"chat_id"`
	}

	// Проверяем, есть ли активный чат в таблице чатов
	err := db.Raw(`
		SELECT chat_id FROM chats 
		WHERE user_id = ? AND active = true
		ORDER BY created_at DESC 
		LIMIT 1
	`, userID).Scan(&chat).Error

	if err != nil || chat.ChatID == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{"active": false})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"active": true, "chat_id": chat.ChatID})
}

// Создание нового чата
func StartChat(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	chatID := fmt.Sprintf("chat_%s", userID) // Генерируем ID чата

	// Создаем запись в таблице chats
	err := db.Exec("INSERT INTO chats (user_id, chat_id, active) VALUES (?, ?, true) ON CONFLICT DO NOTHING", userID, chatID).Error
	if err != nil {
		http.Error(w, "Failed to create chat", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"chat_id": chatID})
}

func GetActiveChats(w http.ResponseWriter, r *http.Request) {
	var activeChats []map[string]interface{}

	err := db.Raw(`
		SELECT c.chat_id, u.name AS user_name, c.created_at 
		FROM chats c 
		LEFT JOIN users u ON c.user_id::int = u.id 
		WHERE c.active = true
	`).Scan(&activeChats).Error

	if err != nil {
		http.Error(w, "Failed to fetch active chats", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(activeChats)
}

func CloseChat(w http.ResponseWriter, r *http.Request) {
	chatID := r.URL.Query().Get("chat_id")
	if chatID == "" {
		http.Error(w, "Missing chat_id", http.StatusBadRequest)
		return
	}

	// Начинаем транзакцию
	tx := db.Begin()

	// Обновляем `active = false`
	if err := tx.Exec("UPDATE chats SET active = false WHERE chat_id = ?", chatID).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to close chat", http.StatusInternalServerError)
		return
	}

	// Удаляем все сообщения этого чата
	if err := tx.Exec("DELETE FROM messages WHERE chat_id = ?", chatID).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to delete chat messages", http.StatusInternalServerError)
		return
	}

	// Фиксируем транзакцию
	tx.Commit()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Chat closed and messages deleted successfully"})
}
