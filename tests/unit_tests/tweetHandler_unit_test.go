package unit_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Модель для твита
type Tweet struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	Content string `json:"content"`
	UserID  uint   `json:"user_id"`
}

// Инициализация тестовой базы данных
func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Создаем таблицу для Tweet
	err = db.AutoMigrate(&Tweet{})
	if err != nil {
		return nil, err
	}

	// Добавляем тестовые данные
	db.Create(&Tweet{Content: "Mock tweet 1", UserID: 1})
	db.Create(&Tweet{Content: "Mock tweet 2", UserID: 1})

	return db, nil
}

// Обработчик для создания твита
func CreateTweetHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tweet Tweet
		if err := json.NewDecoder(r.Body).Decode(&tweet); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Сохраняем твит в базе данных
		result := db.Create(&tweet)
		if result.Error != nil {
			http.Error(w, "Failed to create tweet", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"tweet": tweet,
		})
	}
}

// Обработчик для лайка твита
func LikeTweetHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			TweetID uint `json:"tweetId"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		var tweet Tweet
		result := db.First(&tweet, req.TweetID)
		if result.Error != nil {
			http.Error(w, "Tweet not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "Tweet liked!",
		})
	}
}

// Тест для создания твита
func TestCreateTweet(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	// Подготовка запроса
	tweet := map[string]interface{}{
		"content": "Unit test tweet",
		"user_id": 1,
	}
	body, _ := json.Marshal(tweet)

	req, _ := http.NewRequest("POST", "/tweets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Тестируем обработчик
	w := httptest.NewRecorder()
	handler := CreateTweetHandler(db)
	handler.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "Unit test tweet", response["tweet"].(map[string]interface{})["content"])
}

// Тест для лайка твита
func TestLikeTweet(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	// Подготовка запроса
	likeRequest := map[string]interface{}{
		"tweetId": 1,
	}
	body, _ := json.Marshal(likeRequest)

	req, _ := http.NewRequest("POST", "/like", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Тестируем обработчик
	w := httptest.NewRecorder()
	handler := LikeTweetHandler(db)
	handler.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "Tweet liked!", response["message"])
}
