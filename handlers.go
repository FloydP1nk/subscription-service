package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

// CreateSubscriptionHandler godoc
// @Summary Создать подписку
// @Description Создаёт новую запись о подписке
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body Subscription true "Данные подписки"
// @Success 201 {object} Subscription
// @Failure 400 {string} string "Неверный формат данных"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /subscriptions [post]
func CreateSubscriptionHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("POST /subscriptions вызван")

		var sub Subscription
		if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
			log.Println("Ошибка декодирования JSON:", err)
			http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
			return
		}

		sub.StartDate = time.Date(sub.StartDate.Year(), sub.StartDate.Month(), 1, 0, 0, 0, 0, time.UTC)

		query := `
			INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`
		err := conn.QueryRow(r.Context(), query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).Scan(&sub.ID)
		if err != nil {
			log.Println("Ошибка вставки в БД:", err)
			http.Error(w, "Ошибка при сохранении подписки", http.StatusInternalServerError)
			return
		}
		log.Println("Подписка создана с ID:", sub.ID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(sub)
	}
}

// GetSubscriptionsHandler godoc
// @Summary Получить список подписок
// @Description Возвращает все подписки с возможностью фильтрации по user_id и service_name
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "Фильтр по пользователю (UUID)"
// @Param service_name query string false "Фильтр по названию сервиса"
// @Success 200 {array} Subscription
// @Failure 500 {string} string "Ошибка сервера"
// @Router /subscriptions [get]
func GetSubscriptionsHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET subscriptions вызван")

		userID := r.URL.Query().Get("user_id")
		serviceName := r.URL.Query().Get("service_name")

		query := "SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE 1=1"
		args := []interface{}{}
		argID := 1

		if userID != "" {
			query += fmt.Sprintf(" AND user_id = $%d", argID)
			args = append(args, userID)
			argID++
		}

		if serviceName != "" {
			query += fmt.Sprintf(" AND service_name = $%d", argID)
			args = append(args, serviceName)
			argID++
		}

		rows, err := conn.Query(r.Context(), query, args...)
		if err != nil {
			log.Println("Ошибка запроса к БД:", err)
			http.Error(w, "Ошибка при получении подписок", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		subscriptions := []Subscription{}
		for rows.Next() {
			var sub Subscription
			err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
			if err != nil {
				http.Error(w, "Ошибка при чтении данных", http.StatusInternalServerError)
				log.Println("DB scan error:", err)
				return
			}
			subscriptions = append(subscriptions, sub)
		}

		w.Header().Set("Content-Type", "application/json")
		log.Println("GET /subscriptions вызван")
		json.NewEncoder(w).Encode(subscriptions)
	}
}

// GetSubscriptionByIDHandler godoc
// @Summary Получить подписку по ID
// @Description Возвращает подписку по её UUID
// @Tags subscriptions
// @Produce json
// @Param id path string true "UUID подписки"
// @Success 200 {object} Subscription
// @Failure 404 {string} string "Подписка не найдена"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /subscriptions/{id} [get]
func GetSubscriptionByIDHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id := vars["id"]
		log.Println("GET subscriptions + id вызван")

		var sub Subscription
		query := "SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1"

		err := conn.QueryRow(r.Context(), query, id).Scan(
			&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate,
		)
		if err != nil {
			log.Println("Подписка не найдена или ошибка БД:", err)
			http.Error(w, "Подписка не найдена", http.StatusNotFound)
			return
		}
		log.Println("Подписка найдена с ID:", sub.ID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sub)
	}
}

// UpdateSubscriptionHandler godoc
// @Summary Обновить подписку
// @Description Обновляет данные подписки по ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "UUID подписки"
// @Param subscription body Subscription true "Новые данные подписки"
// @Success 200 {object} Subscription
// @Failure 400 {string} string "Неверный формат данных"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /subscriptions/{id} [put]
func UpdateSubscriptionHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		log.Println("PUT /subscriptions/" + id + " вызван")

		var sub Subscription
		if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
			log.Println("Ошибка декодирования JSON:", err)
			http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
			return
		}

		query := `
			UPDATE subscriptions
			SET service_name = $1, price = $2, start_date = $3, end_date = $4
			WHERE id = $5
			RETURNING id, service_name, price, user_id, start_date, end_date
		`

		err := conn.QueryRow(
			r.Context(),
			query,
			sub.ServiceName, sub.Price, sub.StartDate, sub.EndDate, id,
		).Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)

		if err != nil {
			log.Println("Ошибка обновления подписки:", err)
			http.Error(w, "Ошибка при обновлении подписки", http.StatusInternalServerError)
			return
		}
		log.Println("Подписка обновлена с ID:", sub.ID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sub)
	}
}

// DeleteSubscriptionHandler godoc
// @Summary Удалить подписку
// @Description Удаляет подписку по ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "UUID подписки"
// @Success 200 {string} string "Подписка удалена"
// @Failure 404 {string} string "Подписка не найдена"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /subscriptions/{id} [delete]
func DeleteSubscriptionHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		log.Println("DELETE /subscriptions/" + id + " вызван")

		commandTag, err := conn.Exec(r.Context(), "DELETE FROM subscriptions WHERE id = $1", id)
		if err != nil {
			log.Println("Ошибка удаления подписки:", err)
			http.Error(w, "Ошибка при удалении подписки", http.StatusInternalServerError)
			return
		}

		if commandTag.RowsAffected() == 0 {
			log.Println("Подписка не найдена для удаления:", id)
			http.Error(w, "Подписка не найдена", http.StatusNotFound)
			return
		}
		log.Println("Подписка удалена с ID:", id)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"Подписка удалена"}`))
	}
}

// GetSubscriptionsSummaryHandler godoc
// @Summary Суммарная стоимость подписок
// @Description Возвращает сумму стоимости всех подписок с фильтрацией по user_id, service_name и периоду
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "UUID пользователя"
// @Param service_name query string false "Название сервиса"
// @Param start query string false "Начало периода (YYYY-MM)"
// @Param end query string false "Конец периода (YYYY-MM)"
// @Success 200 {object} map[string]int "total"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /subscriptions/summary [get]
func GetSubscriptionsSummaryHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET /subscriptions/summary вызван")
		// Читаем query-параметры
		userID := r.URL.Query().Get("user_id")
		serviceName := r.URL.Query().Get("service_name")
		start := r.URL.Query().Get("start")
		end := r.URL.Query().Get("end")

		// Базовый SQL-запрос
		query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions WHERE 1=1`
		args := []interface{}{}
		argID := 1

		if userID != "" {
			query += fmt.Sprintf(" AND user_id = $%d", argID)
			args = append(args, userID)
			argID++
		}
		if serviceName != "" {
			query += fmt.Sprintf(" AND service_name = $%d", argID)
			args = append(args, serviceName)
			argID++
		}
		if start != "" {
			query += fmt.Sprintf(" AND to_char(start_date, 'YYYY-MM') >= $%d", argID)
			args = append(args, start)
			argID++
		}
		if end != "" {
			query += fmt.Sprintf(" AND to_char(start_date, 'YYYY-MM') <= $%d", argID)
			args = append(args, end)
			argID++
		}

		var total int
		err := conn.QueryRow(r.Context(), query, args...).Scan(&total)
		if err != nil {
			log.Println("Ошибка при подсчёте суммы:", err)
			http.Error(w, "Ошибка при подсчёте суммы", http.StatusInternalServerError)
			return
		}

		log.Println("Суммарная стоимость подписок:", total)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"total": total})
	}
}
