package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"

	_ "subscription-service/docs"
)

// @title Subscription Service API
// @version 1.0
// @description REST API для управления онлайн-подписками пользователей
// @host localhost:8080
// @BasePath /
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env не найден, используем переменные окружения")
	}

	//БД
	conn, err := ConnectDB()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer conn.Close(context.Background())

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

	port := os.Getenv("SERVER_PORT")
	log.Printf("Сервер запущен на http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
