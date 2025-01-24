package handlers

import (
	"github.com/gorilla/mux"
)

// SetupRouter создает и настраивает маршруты
func SetupRouter() *mux.Router {
	r := mux.NewRouter()

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
