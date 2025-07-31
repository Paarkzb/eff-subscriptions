package service

import (
	"eff-subscriptions/internal/domain/models"
	"github.com/google/uuid"
	"log/slog"
)

type SubscriptionProvider interface {
	Insert(subscription *models.Subscription) error
	Get(id int) (*models.Subscription, error)
	Update(subscription *models.Subscription) error
	Delete(id int) error
	GetAll(serviceName string, price int, userID uuid.UUID, startDate models.CustomDate, filters models.Filters) ([]*models.Subscription, models.Metadata, error)
	GetSubscriptionsSum(userID uuid.UUID, serviceName string, beginDate models.CustomDate, endDate models.CustomDate) (int, error)
}

type SubscriptionService struct {
	log                  *slog.Logger
	subscriptionProvider SubscriptionProvider
}

func NewSubscriptionService(log *slog.Logger, subscriptionProvider SubscriptionProvider) *SubscriptionService {
	return &SubscriptionService{
		subscriptionProvider: subscriptionProvider,
	}
}

func (s *SubscriptionService) Insert(subscription *models.Subscription) error {
	return s.subscriptionProvider.Insert(subscription)
}
func (s *SubscriptionService) Get(id int) (*models.Subscription, error) {
	return s.subscriptionProvider.Get(id)
}
func (s *SubscriptionService) Update(subscription *models.Subscription) error {
	return s.subscriptionProvider.Update(subscription)
}
func (s *SubscriptionService) Delete(id int) error {
	return s.subscriptionProvider.Delete(id)
}
func (s *SubscriptionService) GetAll(serviceName string, price int, userID uuid.UUID, startDate models.CustomDate, filters models.Filters) ([]*models.Subscription, models.Metadata, error) {
	return s.subscriptionProvider.GetAll(serviceName, price, userID, startDate, filters)
}

func (s *SubscriptionService) GetSubscriptionsSum(userID uuid.UUID, serviceName string, beginDate models.CustomDate, endDate models.CustomDate) (int, error) {
	return s.subscriptionProvider.GetSubscriptionsSum(userID, serviceName, beginDate, endDate)
}
