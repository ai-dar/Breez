package main

import (
	"log"
	"net/http"

	"breez/handlers"
	"breez/models"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectDB() *gorm.DB {
	dsn := "user=postgres password=1234 dbname=breez sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}
	database.AutoMigrate(&models.User{})
	log.Println("Database connected and migrated")
	return database
}

func main() {
	db := connectDB()
	handlers.InitDB(db)

	r := mux.NewRouter()

	// Обслуживание статических файлов
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// CRUD маршруты
	r.HandleFunc("/create", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/update", handlers.UpdateUser).Methods("PUT")
	r.HandleFunc("/delete", handlers.DeleteUser).Methods("DELETE")
	r.HandleFunc("/user", handlers.GetUserByID).Methods("GET")

	log.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
