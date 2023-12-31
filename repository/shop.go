package repository

import (
	"context"
	"github.com/restore/shop/entity"
	"gorm.io/gorm"
	"time"
)

type Shop struct {
	db *gorm.DB
}

func NewShop(db *gorm.DB) *Shop {
	return &Shop{
		db: db,
	}
}

func (s *Shop) CreateRequest(ctx context.Context, request *entity.Request) error {
	result := s.db.Create(request)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *Shop) UpdateRequest(ctx context.Context, id int, request *entity.Request) error {
	result := entity.Request{ID: id}
	res := s.db.First(&result)
	if res.Error != nil {
		return res.Error
	}

	result.Status = request.Status
	result.Track = request.Track

	res = s.db.Save(result)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (s *Shop) ConfirmRequests(ctx context.Context, paymentID string) error {
	res := s.db.Model(&entity.Request{}).
		Where("payment_id = ?", paymentID).
		Update("status", "preparing")
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (s *Shop) GetRequestByPayment(ctx context.Context, paymentID string) ([]entity.Request, error) {
	var result []entity.Request
	query := s.db.Where("payment_id = ?", paymentID)
	res := query.Find(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	return result, nil
}

func (s *Shop) SearchRequest(ctx context.Context, id int, status string, init, end time.Time) ([]entity.Request, error) {
	var result []entity.Request
	query := s.db.Where("store_id = ? AND status != ?", id, "created")

	if status != "" {
		query.Where("status = ?", status)
	}

	t := time.Time{}
	if init != t {
		query.Where("created_at > ?", init)
	}
	if end != t {
		query.Where("created_at < ?", end)
	}

	res := query.Find(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	return result, nil
}

func (s *Shop) SearchProfileRequest(ctx context.Context, id int, status string, init, end time.Time) ([]entity.Request, error) {
	var result []entity.Request
	query := s.db.Where("user_id = ? AND status != ?", id, "created")

	if status != "" {
		query.Where("status = ?", status)
	}

	t := time.Time{}
	if init != t {
		query.Where("created_at > ?", init)
	}
	if end != t {
		query.Where("created_at < ?", end)
	}

	res := query.Find(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	return result, nil
}

func (s *Shop) CreatePayment(ctx context.Context, payment *entity.Payment) (int, error) {
	result := s.db.Create(payment)
	if result.Error != nil {
		return 0, result.Error
	}
	return payment.ID, nil
}

func (s *Shop) UpdatePayment(ctx context.Context, id int, payment *entity.Payment) error {
	result := entity.Payment{ID: id}
	res := s.db.First(&result)
	if res.Error != nil {
		return res.Error
	}

	result.Status = payment.Status

	res = s.db.Save(result)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (s *Shop) GetPayments(ctx context.Context, id int) ([]entity.Payment, error) {
	var result []entity.Payment
	query := s.db.Where("store_id = ?", id)

	res := query.Find(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	return result, nil
}

func (s *Shop) SearchPayment(ctx context.Context, status string, init, end time.Time) ([]entity.Payment, error) {
	var result []entity.Payment
	query := s.db.Where("")

	if status != "" {
		query.Where("status = ?", status)
	}

	t := time.Time{}
	if init != t {
		query.Where("created_at > ?", init)
	}
	if end != t {
		query.Where("created_at < ?", end)
	}

	res := query.Find(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	return result, nil
}
