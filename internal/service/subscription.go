package service

import (
	"database/sql"
	"errors"
	"subscriptions-service/internal/logger"
	"subscriptions-service/internal/model"
	"subscriptions-service/internal/repository"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SubscriptionService struct {
	repo *repository.SubscriptionRepo
}

func NewSubscriptionService(r *repository.SubscriptionRepo) *SubscriptionService {
	return &SubscriptionService{repo: r}
}

func (s *SubscriptionService) Create(sub model.Subscription) error {
	logger.Log.Info("create subscription called",
	zap.String("service", sub.ServiceName),
	zap.String("user_id", sub.UserID.String()),
	zap.Int("price", sub.Price),
	zap.String("start_date", sub.StartDate.Format("2006-01-02")),
)
	return s.repo.Create(sub)
}

func (s *SubscriptionService) GetAll() ([]model.Subscription, error) {
	logger.Log.Info("GetAll service called")
	return s.repo.GetAll()
}

func (s *SubscriptionService) GetByID(id int) (model.Subscription, error) {
	logger.Log.Info("GetByID service called", zap.Int("id", id))
	return s.repo.GetByID(id)
}

func (s *SubscriptionService) Update(sub model.Subscription) error {
	logger.Log.Info("Update service called", zap.String("service", sub.ServiceName))
	_, err := s.repo.GetByID(sub.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Warn("subscription not found", zap.Int("id", sub.ID))
			return errors.New("not found")
		}

		logger.Log.Error("db error", zap.Error(err))
		return err
	}
	return s.repo.Update(sub)
}

func (s *SubscriptionService) Delete(id int) error {
	logger.Log.Info("Delete service called", zap.Int("id", id))
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Warn("subscription not found", zap.Int("id", id))
			return errors.New("not found")
		}

		logger.Log.Error("db error", zap.Error(err))
		return err
	}

	return s.repo.Delete(id)
}

func (s *SubscriptionService) GetTotalSum(userID uuid.UUID, serviceName string, from, to *time.Time) (int, error) {
	logger.Log.Info("GetTotalSum service called", zap.String("user_id", userID.String()), zap.String("service_name", serviceName))
	return s.repo.GetTotalSum(userID, serviceName, from, to)
}