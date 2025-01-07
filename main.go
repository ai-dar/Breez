package main

import (
	"breez/handlers"
	"breez/models"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var limiter = rate.NewLimiter(1, 3)
var log = logrus.New()

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			log.WithField("ip", r.RemoteAddr).Warn("Rate limit exceeded")
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func connectDB() *gorm.DB {
	dsn := "user=postgres password=1234 dbname=breez sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}
	database.AutoMigrate(&models.User{}, &models.Tweet{}, &models.Like{})
	log.Println("Database connected and migrated")
	return database
}

func setupLogger() {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
	log.Info("Logger initialized")
}

func main() {
	setupLogger()
	db := connectDB()
	handlers.InitDB(db)

	r := mux.NewRouter()
	r.Use(rateLimitMiddleware)
	r.HandleFunc("/admin/register", handlers.RegisterAdmin).Methods("POST")
	r.HandleFunc("/admin/send-emails", handlers.SendEmailToAllUsers).Methods("POST")
	r.HandleFunc("/user/me", handlers.GetCurrentUser).Methods("GET")
	r.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	r.HandleFunc("/login", handlers.LoginUser).Methods("POST")
	r.HandleFunc("/tweets", handlers.GetTweetsWithFilters).Methods("GET")
	r.HandleFunc("/tweets", handlers.CreateTweet).Methods("POST")
	r.HandleFunc("/like", handlers.LikeTweet).Methods("POST")
	r.HandleFunc("/check-auth", handlers.CheckAuth).Methods("GET")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("static", "index.html"))
	})

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
