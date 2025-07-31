package models

import (
	"eff-subscriptions/internal/validator"
	"github.com/google/uuid"
	"time"
)

type Subscription struct {
	ID          int         `json:"id"`
	ServiceName string      `json:"service_name"`
	Price       *int        `json:"price"`
	UserID      uuid.UUID   `json:"user_id"`
	StartDate   CustomDate  `json:"start_date"`
	EndDate     *CustomDate `json:"end_date,omitzero"`
	CreatedAt   time.Time   `json:"-"`
	Version     int         `json:"version"`
}

func ValidateSubscription(v *validator.Validator, subscription *Subscription) {
	v.Check(subscription.ServiceName != "", "service_name", "must be provided")
	v.Check(len(subscription.ServiceName) <= 500, "service_name", "must not be more than 500 bytes long")

	v.Check(subscription.Price != nil, "price", "must be provided")
	v.Check(*subscription.Price > -1, "price", "must be a positive integer")

	v.Check(subscription.UserID != uuid.Nil, "user_id", "must not be empty")

	v.Check(subscription.StartDate != CustomDate{}, "start_date", "must be provided")
	if subscription.EndDate != nil {
		v.Check(!subscription.StartDate.Time().After(subscription.EndDate.Time()), "start_date", "must be before end_date")
	}
}

// CreateSubscriptionRequest subscription request struct
// @Description subscription
type CreateSubscriptionRequest struct {
	ServiceName string      `json:"service_name"`
	Price       *int        `json:"price"`
	UserID      uuid.UUID   `json:"user_id"`
	StartDate   CustomDate  `json:"start_date"`
	EndDate     *CustomDate `json:"end_date,omitzero"`
}

// UpdateSubscriptionRequest subscription request struct for update
// @Description update subscription struct
type UpdateSubscriptionRequest struct {
	ServiceName *string     `json:"service_name"`
	Price       *int        `json:"price"`
	UserID      *uuid.UUID  `json:"user_id"`
	StartDate   *CustomDate `json:"start_date"`
	EndDate     *CustomDate `json:"end_date,omitzero"`
}

// SubscriptionResponse subscription response struct
// @Description subscription
type SubscriptionResponse struct {
	Subscription *Subscription `json:"subscription"`
}

// SubscriptionsListResponse subscription list response struct
// @Description subscription list with metadata for pagination
type SubscriptionsListResponse struct {
	Metadata     Metadata        `json:"metadata"`
	Subscription []*Subscription `json:"subscriptions"`
}
