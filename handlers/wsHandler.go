package handlers

import (
	"net/http"
	"sync"
	"time"

	"breez/models"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]string) // WebSocket -> Chat ID
var broadcast = make(chan models.Message)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var mu sync.Mutex

// Обработчик WebSocket
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Upgrade") != "websocket" {
		http.Error(w, "Expected WebSocket handshake", http.StatusUpgradeRequired)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("❌ Ошибка подключения WebSocket:", err)
		return
	}
	defer ws.Close()

	chatID := r.URL.Query().Get("chat_id")
	if chatID == "" {
		log.Println("⚠️ Ошибка: chat_id не передан")
		return
	}
	cookie, err := r.Cookie("user_id")
	if err != nil {
		log.Println("❌ Ошибка получения user_id:", err)
		return
	}
	userID := cookie.Value
	// Получаем роль пользователя из БД
	var role string
	if err := db.Raw("SELECT role FROM users WHERE id = ?", userID).Scan(&role).Error; err != nil {
		log.Println("❌ Ошибка получения роли:", err)
		return
	}

	mu.Lock()
	clients[ws] = chatID
	mu.Unlock()

	log.Println("🟢 Клиент подключен к чату:", chatID)

	var messages []models.Message
	if err := db.Where("chat_id = ?", chatID).Order("created_at ASC").Find(&messages).Error; err != nil {
		log.Println("❌ Ошибка загрузки истории чата:", err)
	} else {
		for _, msg := range messages {
			msgFormatted := map[string]string{
				"time":    msg.Timestamp, // Часы:Минуты:Секунды
				"sender":  msg.Sender,
				"content": msg.Content,
			}
			log.Println("📤 Отправка в WebSocket:", msgFormatted)
			err := ws.WriteJSON(msgFormatted)
			if err != nil {
				log.Println("❌ Ошибка отправки истории:", err)
			}
		}
	}

	for {
		var msg models.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("❌ Ошибка чтения WebSocket:", err)
			mu.Lock()
			delete(clients, ws)
			mu.Unlock()
			break
		}
		if role == "admin" {
			msg.Sender = "Админ"
		} else {
			msg.Sender = "Клиент"
		}

		// 🔹 Сохраняем сообщение в БД
		msg.ChatID = chatID
		msg.Timestamp = time.Now().Format("15:04:05")

		if err := db.Create(&msg).Error; err != nil {
			log.Println("❌ Ошибка сохранения сообщения в БД:", err)
		}

		log.Println("📩 Получено сообщение:", msg.Content)
		broadcast <- msg
	}
}

// Отправка сообщений всем клиентам
func HandleMessages() {
	for {
		msg := <-broadcast

		// ✅ Форматируем перед отправкой
		msgFormatted := map[string]string{
			"time":    msg.Timestamp,
			"sender":  msg.Sender,
			"content": msg.Content,
		}
		log.Println("📩 Отправка сообщения:", msg)
		mu.Lock()
		for client, chatID := range clients {
			if chatID == msg.ChatID {
				err := client.WriteJSON(msgFormatted)
				if err != nil {
					log.Println("❌ Ошибка отправки сообщения:", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
		log.Println("📤 Отправка в WebSocket:", msgFormatted)
		mu.Unlock()
	}
}
