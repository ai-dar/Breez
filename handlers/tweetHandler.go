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
	cookie, err := r.Cookie("user_id")
	if err != nil {
		logger.WithField("error", err).Warn("Unauthorized access - missing cookie")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		logger.WithField("error", err).Warn("Invalid user ID in cookie")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var tweet models.Tweet
	if err := json.NewDecoder(r.Body).Decode(&tweet); err != nil {
		logger.WithField("error", err).Warn("Failed to decode JSON payload")
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if len(tweet.Content) == 0 {
		logger.Warn("Attempt to create tweet with empty content")
		http.Error(w, "Tweet content cannot be empty", http.StatusBadRequest)
		return
	}

	tweet.UserID = uint(userID)

	if err := db.Create(&tweet).Error; err != nil {
		logger.WithField("error", err).Error("Failed to create tweet in database")
		http.Error(w, "Failed to create tweet", http.StatusInternalServerError)
		return
	}

	logger.WithFields(logrus.Fields{
		"user_id": userID,
		"tweet":   tweet.Content,
	}).Info("Tweet created successfully")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tweet)
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
	if sort != "" {
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

// LikeTweet handles liking and unliking a tweet
func LikeTweet(w http.ResponseWriter, r *http.Request) {
	var likeRequest struct {
		TweetID uint `json:"tweetId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&likeRequest); err != nil {
		logger.WithField("error", err).Warn("Failed to decode like request payload")
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// Получить user_id из cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		logger.WithField("error", err).Warn("Unauthorized access - missing cookie")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		logger.WithField("error", err).Warn("Invalid user ID in cookie")
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Проверить, существует ли лайк
	var existingLike models.Like
	if err := db.Where("tweet_id = ? AND user_id = ?", likeRequest.TweetID, uint(userID)).First(&existingLike).Error; err == nil {
		// Если лайк существует, удаляем его
		if err := db.Delete(&existingLike).Error; err != nil {
			logger.WithField("error", err).Error("Failed to unlike tweet")
			http.Error(w, "Failed to unlike tweet", http.StatusInternalServerError)
			return
		}

		logger.WithFields(logrus.Fields{
			"action":   "unlike",
			"user_id":  userID,
			"tweet_id": likeRequest.TweetID,
		}).Info("Tweet unliked successfully")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Tweet unliked!"})
		return
	}

	// Если лайк не существует, добавляем его
	newLike := models.Like{
		TweetID: likeRequest.TweetID,
		UserID:  uint(userID),
	}

	if err := db.Create(&newLike).Error; err != nil {
		logger.WithField("error", err).Error("Failed to like tweet")
		http.Error(w, "Failed to like tweet", http.StatusInternalServerError)
		return
	}

	logger.WithFields(logrus.Fields{
		"action":   "like",
		"user_id":  userID,
		"tweet_id": likeRequest.TweetID,
	}).Info("Tweet liked successfully")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Tweet liked!"})
}
