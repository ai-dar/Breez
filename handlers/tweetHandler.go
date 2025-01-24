package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"breez/models"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func init() {
	logger.SetFormatter(&logrus.JSONFormatter{})
}

// CreateTweet handles creating a new tweet
func CreateTweet(w http.ResponseWriter, r *http.Request) {
	// Получаем cookie с user_id
	userCookie, err := r.Cookie("user_id")
	if err != nil {
		logger.WithField("error", err).Warn("Unauthorized access - missing cookie")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(userCookie.Value)
	if err != nil || userID <= 0 {
		logger.WithField("error", err).Warn("Invalid user ID in cookie")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Определяем источник входа
	sourceCookie, err := r.Cookie("source")
	source := "account" // По умолчанию считаем, что вход выполнен через аккаунт
	if err == nil && sourceCookie.Value == "github" {
		source = "github"
	}

	// Инициализируем имя пользователя
	var userName string

	if source == "github" {
		// Извлекаем имя пользователя из базы данных
		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			logger.WithField("user_id", userID).Warn("User not found")
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}
		userName = user.Name
	} else {
		// Если вход выполнен через аккаунт, берём имя из cookie (или можно запросить в базе)
		userNameCookie, err := r.Cookie("user_name")
		if err != nil {
			logger.Warn("Missing user_name cookie for account login")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userName = userNameCookie.Value
	}

	// Декодируем JSON-полезную нагрузку для создания твита
	var tweet models.Tweet
	if err := json.NewDecoder(r.Body).Decode(&tweet); err != nil {
		logger.WithField("error", err).Warn("Failed to decode JSON payload")
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Проверяем, что содержание твита не пустое
	if len(tweet.Content) == 0 {
		logger.Warn("Attempt to create tweet with empty content")
		http.Error(w, "Tweet content cannot be empty", http.StatusBadRequest)
		return
	}

	// Устанавливаем идентификатор пользователя в твите
	tweet.UserID = uint(userID)

	// Сохраняем твит в базе данных
	if err := db.Create(&tweet).Error; err != nil {
		logger.WithField("error", err).Error("Failed to create tweet in database")
		http.Error(w, "Failed to create tweet", http.StatusInternalServerError)
		return
	}

	// Логируем успешное создание твита
	logger.WithFields(logrus.Fields{
		"user_id":   userID,
		"user_name": userName,
		"tweet":     tweet.Content,
	}).Info("Tweet created successfully")

	// Возвращаем успешный ответ с именем пользователя и созданным твитом
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tweet": tweet,
		"user":  userName,
	})
}

// GetTweetsWithFilters handles filtering, sorting, and pagination
func GetTweetsWithFilters(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры фильтрации, сортировки и пагинации
	filter := r.URL.Query().Get("filter")
	sort := r.URL.Query().Get("sort")
	page := r.URL.Query().Get("page")

	limit := 5 // Максимальное количество твитов на одной странице
	offset := 0
	pageInt := 1
	var err error

	if page != "" {
		pageInt, err = strconv.Atoi(page)
		if err != nil || pageInt < 1 {
			logger.WithField("page", page).Warn("Invalid page number")
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			return
		}
		offset = (pageInt - 1) * limit
	}

	// Формируем запрос к базе данных
	query := db.Preload("User").Model(&models.Tweet{})

	// Применяем фильтрацию
	if filter != "" {
		query = query.Where("content LIKE ? OR user_id IN (SELECT id FROM users WHERE name LIKE ?)", "%"+filter+"%", "%"+filter+"%")
	}

	// Подсчёт общего количества твитов
	var totalTweets int64
	query.Count(&totalTweets)

	// Определяем общее количество страниц
	totalPages := (totalTweets + int64(limit) - 1) / int64(limit)

	// Если пользователь запрашивает несуществующую страницу
	if pageInt > int(totalPages) {
		logger.WithFields(logrus.Fields{
			"action":      "pagination",
			"requested":   pageInt,
			"total_pages": totalPages,
		}).Warn("Requested page exceeds total pages")
		http.Error(w, "Page does not exist", http.StatusNotFound)
		return
	}

	// Применяем сортировку
	validSortFields := map[string]bool{
		"created_at": true,
		"likes":      true,
		"user_id":    true,
	}

	if sort != "" {
		if !validSortFields[sort] {
			logger.WithField("sort", sort).Warn("Invalid sort field")
			http.Error(w, "Invalid sort field", http.StatusBadRequest)
			return
		}
		query = query.Order(sort)
	} else {
		query = query.Order("created_at DESC")
	}

	// Применяем пагинацию
	query = query.Limit(limit).Offset(offset)

	// Выполняем запрос
	var tweets []models.Tweet
	if err := query.Find(&tweets).Error; err != nil {
		logger.WithField("error", err).Error("Failed to fetch tweets")
		http.Error(w, "Failed to fetch tweets", http.StatusInternalServerError)
		return
	}

	// Подсчёт лайков и формирование ответа
	var result []map[string]interface{}
	for _, tweet := range tweets {
		var likeCount int64
		db.Model(&models.Like{}).Where("tweet_id = ?", tweet.ID).Count(&likeCount)

		result = append(result, map[string]interface{}{
			"id":      tweet.ID,
			"content": tweet.Content,
			"user":    tweet.User,
			"likes":   likeCount,
		})
	}

	logger.WithFields(logrus.Fields{
		"action":      "fetch_tweets",
		"count":       len(result),
		"filter":      filter,
		"sort":        sort,
		"page":        pageInt,
		"total_pages": totalPages,
	}).Info("Tweets fetched successfully")

	// Отправка результата в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tweets":      result,
		"total_pages": totalPages,
	})
}

func isAdmin(userID uint) (bool, error) {
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return false, err
	}
	return user.Role == "admin", nil
}

func UpdateTweet(w http.ResponseWriter, r *http.Request) {
	var userID uint
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID64, err := strconv.ParseUint(cookie.Value, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	userID = uint(userID64)

	isAdmin, err := isAdmin(userID)
	if err != nil {
		http.Error(w, "Failed to verify user role", http.StatusInternalServerError)
		return
	}
	if !isAdmin {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	var updateRequest struct {
		TweetID uint   `json:"tweet_id"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(updateRequest.Content) == 0 {
		http.Error(w, "Tweet content cannot be empty", http.StatusBadRequest)
		return
	}

	var tweet models.Tweet
	if err := db.First(&tweet, updateRequest.TweetID).Error; err != nil {
		http.Error(w, "Tweet not found", http.StatusNotFound)
		return
	}

	tweet.Content = updateRequest.Content
	if err := db.Save(&tweet).Error; err != nil {
		http.Error(w, "Failed to update tweet", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Tweet updated successfully"})
}

func DeleteTweet(w http.ResponseWriter, r *http.Request) {
	var userID uint
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID64, err := strconv.ParseUint(cookie.Value, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	userID = uint(userID64)

	isAdmin, err := isAdmin(userID)
	if err != nil {
		http.Error(w, "Failed to verify user role", http.StatusInternalServerError)
		return
	}
	if !isAdmin {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	var deleteRequest struct {
		TweetID uint `json:"tweet_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&deleteRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := db.Delete(&models.Tweet{}, deleteRequest.TweetID).Error; err != nil {
		http.Error(w, "Failed to delete tweet", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Tweet deleted successfully"})
}

// LikeTweet handles liking and unliking a tweet
func LikeTweet(w http.ResponseWriter, r *http.Request) {
	var likeRequest struct {
		TweetID uint `json:"tweetId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&likeRequest); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// Получить user_id из cookie
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

	// Проверить, существует ли лайк
	var existingLike models.Like
	if err := db.Where("tweet_id = ? AND user_id = ?", likeRequest.TweetID, uint(userID)).First(&existingLike).Error; err == nil {
		if err := db.Delete(&existingLike).Error; err != nil {
			http.Error(w, "Failed to unlike tweet", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Tweet unliked!"})
		return
	}

	newLike := models.Like{
		TweetID: likeRequest.TweetID,
		UserID:  uint(userID),
	}

	if err := db.Create(&newLike).Error; err != nil {
		http.Error(w, "Failed to like tweet", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Tweet liked!"})
}
