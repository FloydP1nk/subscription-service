package main

import "time"

// Subscription представляет подписку пользователя
// @Description Сущность подписки
type Subscription struct {
	ID          int        `json:"id" example:"1"`
	ServiceName string     `json:"service_name" example:"MyService"`
	Price       int        `json:"price" example:"100"`
	UserID      string     `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	StartDate   time.Time  `json:"start_date" example:"2025-08-29T14:30:00Z"`
	EndDate     *time.Time `json:"end_date" example:"2025-09-29T14:30:00Z"`
}

type CreateSubscriptionRequest struct {
	UserID      string     `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ServiceName string     `json:"service_name" example:"MyService"`
	Price       int        `json:"price" example:"100"`
	StartDate   time.Time  `json:"start_date" example:"2025-08-29T14:30:00Z"`
	EndDate     *time.Time `json:"end_date" example:"2025-09-29T14:30:00Z"`
}
