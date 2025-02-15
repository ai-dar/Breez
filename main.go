package main

import (
	"breez/handlers"
	"breez/models"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/time/rate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var limiter = rate.NewLimiter(1, 3)
var log = logrus.New()
var githubOauthConfig *oauth2.Config

func init() {
	// Загрузка переменных окружения из .env
	if err := godotenv.Load(".env"); err != nil {
		log.Warn("No .env file found, using system environment variables")
	}

	githubOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     github.Endpoint,
	}
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/user/me" {
			next.ServeHTTP(w, r)
			return
		}

		if !limiter.Allow() {
			log.WithField("ip", r.RemoteAddr).Warn("Rate limit exceeded")
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func connectDB() *gorm.DB {
	// dsn := os.Getenv("DATABASE_URL")
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
	file, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(file)
	log.SetLevel(logrus.InfoLevel)

	log.Info("Logger initialized and writing to logs.txt")
}

func main() {
	setupLogger()
	db := connectDB()
	handlers.InitDB(db)
	handlers.InitGitHubOAuth()

	r := mux.NewRouter()
	r.Use(rateLimitMiddleware)

	// Маршруты GitHub OAuth
	r.HandleFunc("/auth/github/login", handlers.GitHubLoginHandler).Methods("GET")
	r.HandleFunc("/auth/github/callback", func(w http.ResponseWriter, r *http.Request) { handlers.GitHubCallbackHandler(db, w, r) }).Methods("GET")

	r.HandleFunc("/pay", func(w http.ResponseWriter, r *http.Request) { handlers.HandlePayment(db, w, r) }).Methods("POST")

	// Роуты для администратора
	r.HandleFunc("/admin/register", handlers.RegisterAdmin).Methods("POST")
	r.HandleFunc("/admin/send-emails", handlers.SendEmailToAllUsers).Methods("POST")
	r.HandleFunc("/admin/tweet/update", handlers.UpdateTweet).Methods("PUT")
	r.HandleFunc("/admin/tweet/delete", handlers.DeleteTweet).Methods("DELETE")
	r.HandleFunc("/ws", handlers.HandleConnections)
	go handlers.HandleMessages()
	r.HandleFunc("/api/check-active-chat", handlers.CheckActiveChat).Methods("GET")
	r.HandleFunc("/api/start-chat", handlers.StartChat).Methods("POST")
	r.HandleFunc("/api/admin/close-chat", handlers.CloseChat).Methods("POST")
	r.HandleFunc("/api/admin/active-chats", handlers.GetActiveChats).Methods("GET")

	// Роуты для пользователя
	r.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	r.HandleFunc("/login", handlers.LoginUser).Methods("POST")
	r.HandleFunc("/user/me", handlers.GetCurrentUser).Methods("GET")
	r.HandleFunc("/check-auth", handlers.CheckAuth).Methods("GET")
	r.HandleFunc("/verify", handlers.VerifyEmail).Methods("GET")
	r.HandleFunc("/api/verify", handlers.VerifyEmail).Methods("GET")

	// Роуты для твитов
	r.HandleFunc("/tweets", handlers.GetTweetsWithFilters).Methods("GET")
	r.HandleFunc("/tweets", handlers.CreateTweet).Methods("POST")
	r.HandleFunc("/like", handlers.LikeTweet).Methods("POST")

	// Статические файлы
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
