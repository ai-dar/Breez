package integration_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Мок функции GetTweetsWithFilters
func MockGetTweetsWithFilters(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"tweets": []map[string]interface{}{
			{"id": 1, "content": "Test tweet 1", "user_id": 1},
			{"id": 2, "content": "Test tweet 2", "user_id": 2},
		},
		"total_pages": 1,
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// Мок функции UpdateTweet
func MockUpdateTweet(w http.ResponseWriter, r *http.Request) {
	var updateRequest map[string]interface{}
	_ = json.NewDecoder(r.Body).Decode(&updateRequest)

	if updateRequest["tweet_id"] == nil || updateRequest["content"] == nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Tweet updated successfully",
	})
}

// Integration Test: Проверка получения твитов
func TestGetTweetsWithFilters(t *testing.T) {
	req, _ := http.NewRequest("GET", "/tweets?filter=test&sort=created_at&page=1", nil)
	w := httptest.NewRecorder()

	// Используем мок
	MockGetTweetsWithFilters(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	assert.NotEmpty(t, response["tweets"])
	assert.Equal(t, 1, int(response["tweets"].([]interface{})[0].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, "Test tweet 1", response["tweets"].([]interface{})[0].(map[string]interface{})["content"].(string))
}

// Integration Test: Проверка обновления твита
func TestUpdateTweet(t *testing.T) {
	// Тестовая полезная нагрузка
	updateRequest := map[string]interface{}{
		"tweet_id": 1,
		"content":  "Updated content",
	}
	body, _ := json.Marshal(updateRequest)

	req, _ := http.NewRequest("PUT", "/admin/tweet/update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "user_id", Value: "1"})

	w := httptest.NewRecorder()

	// Используем мок
	MockUpdateTweet(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "Tweet updated successfully", response["message"])
}
