package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "subscription-service/docs"
)

// @title Subscription Service API
// @version 1.0
// @description REST API для управления онлайн-подписками пользователей
// @host localhost:8080
// @BasePath /
func main() {

	// .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Ошибка при загрузке .env: %v", err)
	}

	// БД
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer conn.Close(context.Background())

	log.Println("База данных подключена")

	r := mux.NewRouter()

	// Ручки
	r.HandleFunc("/subscriptions/summary", GetSubscriptionsSummaryHandler(conn)).Methods("GET") // Сначала конкретный маршрут
	r.HandleFunc("/subscriptions/{id}", GetSubscriptionByIDHandler(conn)).Methods("GET")
	r.HandleFunc("/subscriptions/{id}", UpdateSubscriptionHandler(conn)).Methods("PUT")
	r.HandleFunc("/subscriptions/{id}", DeleteSubscriptionHandler(conn)).Methods("DELETE")
	r.HandleFunc("/subscriptions", CreateSubscriptionHandler(conn)).Methods("POST")
	r.HandleFunc("/subscriptions", GetSubscriptionsHandler(conn)).Methods("GET")
	// Swagger UI доступен по адресу http://localhost:8080/swagger/index.html
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	port := "8080"
	log.Printf("Сервер запущен на http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
