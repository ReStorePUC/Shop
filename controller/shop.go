package controller

import (
	"context"
	"errors"
	paymentpb "github.com/ReStorePUC/protobucket/payment"
	pb "github.com/ReStorePUC/protobucket/user"
	"github.com/restore/shop/config"
	"github.com/restore/shop/entity"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type repository interface {
	CreateRequest(ctx context.Context, request *entity.Request) error
	UpdateRequest(ctx context.Context, id int, request *entity.Request) error
	SearchRequest(ctx context.Context, id int, status string, init, end time.Time) ([]entity.Request, error)

	CreatePayment(ctx context.Context, payment *entity.Payment) (int, error)
	UpdatePayment(ctx context.Context, id int, payment *entity.Payment) error
	GetPayments(ctx context.Context, id int) ([]entity.Payment, error)
	SearchPayment(ctx context.Context, status string, init, end time.Time) ([]entity.Payment, error)
}

type Shop struct {
	repo    repository
	service pb.UserClient
	payment paymentpb.PaymentClient
}

func NewShop(r repository, s pb.UserClient, p paymentpb.PaymentClient) *Shop {
	return &Shop{
		repo:    r,
		service: s,
		payment: p,
	}
}

func (s *Shop) CreateRequest(ctx context.Context, request *entity.Create) (string, error) {
	log := zap.NewNop()

	items := []*paymentpb.Item{}
	for _, item := range request.Items {
		items = append(items, &paymentpb.Item{
			Title:     item.ItemName,
			Quantity:  1,
			UnitPrice: float32(item.Price + item.Tax),
		})
	}

	payment, err := s.payment.CreatePayment(ctx, &paymentpb.CreatePaymentRequest{
		Items: items,
	})
	if err != nil {
		log.Error(
			"error to create payment",
			zap.Error(err),
		)
		return "", err
	}

	for _, item := range request.Items {
		item.PaymentID = payment.Id
		item.Status = "preparing"
		err = s.repo.CreateRequest(ctx, &item)
		if err != nil {
			log.Error(
				"error to create request",
				zap.Error(err),
			)
			return "", err
		}
	}

	return payment.Id, nil
}

func (s *Shop) UpdateRequest(ctx context.Context, id string, request *entity.Request) error {
	log := zap.NewNop()

	admin := ctx.Value(config.EmailHeader)
	result, err := s.service.GetUser(ctx, &pb.GetUserRequest{
		Email: admin.(string),
	})
	if err != nil {
		log.Error(
			"error getting admin",
			zap.Error(err),
		)
		return err
	}
	if !result.IsAdmin {
		log.Error(
			"unauthorized action",
		)
		return errors.New("unauthorized action")
	}

	requestID, err := strconv.Atoi(id)
	if err != nil {
		log.Error(
			"error validating id",
			zap.Error(err),
		)
		return err
	}

	err = s.repo.UpdateRequest(ctx, requestID, request)
	if err != nil {
		log.Error(
			"error to update request",
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (s *Shop) SearchRequest(ctx context.Context, storeID, status, initialDate, endDate string) ([]entity.Request, error) {
	log := zap.NewNop()

	admin := ctx.Value(config.EmailHeader)
	user, err := s.service.GetUser(ctx, &pb.GetUserRequest{
		Email: admin.(string),
	})
	if err != nil {
		log.Error(
			"error getting admin",
			zap.Error(err),
		)
		return nil, err
	}
	if !user.IsAdmin {
		log.Error(
			"unauthorized action",
		)
		return nil, errors.New("unauthorized action")
	}

	id, err := strconv.Atoi(storeID)
	if err != nil {
		log.Error(
			"error validating id",
			zap.Error(err),
		)
		return nil, err
	}

	var init time.Time
	if initialDate != "" {
		init, err = time.Parse(time.RFC3339, initialDate)
		if err != nil {
			log.Error(
				"error validating initial date",
				zap.Error(err),
			)
			return nil, err
		}
	}

	var end time.Time
	if endDate != "" {
		end, err = time.Parse(time.RFC3339, endDate)
		if err != nil {
			log.Error(
				"error validating end date",
				zap.Error(err),
			)
			return nil, err
		}
	}

	result, err := s.repo.SearchRequest(ctx, id, status, init, end)
	if err != nil {
		log.Error(
			"error to search requests",
			zap.Error(err),
		)
		return nil, err
	}

	return result, nil
}

func (s *Shop) CreatePayment(ctx context.Context, payment *entity.Payment) (int, error) {
	log := zap.NewNop()

	admin := ctx.Value(config.EmailHeader)
	user, err := s.service.GetUser(ctx, &pb.GetUserRequest{
		Email: admin.(string),
	})
	if err != nil {
		log.Error(
			"error getting admin",
			zap.Error(err),
		)
		return 0, err
	}
	if !user.IsAdmin {
		log.Error(
			"unauthorized action",
		)
		return 0, errors.New("unauthorized action")
	}

	id, err := s.repo.CreatePayment(ctx, payment)
	if err != nil {
		log.Error(
			"error to create payment",
			zap.Error(err),
		)
		return 0, err
	}

	return id, nil
}

func (s *Shop) UpdatePayment(ctx context.Context, id string, payment *entity.Payment) error {
	log := zap.NewNop()

	admin := ctx.Value(config.EmailHeader)
	user, err := s.service.GetUser(ctx, &pb.GetUserRequest{
		Email: admin.(string),
	})
	if err != nil {
		log.Error(
			"error getting admin",
			zap.Error(err),
		)
		return err
	}
	if !user.IsAdmin {
		log.Error(
			"unauthorized action",
		)
		return errors.New("unauthorized action")
	}

	paymentID, err := strconv.Atoi(id)
	if err != nil {
		log.Error(
			"error validating id",
			zap.Error(err),
		)
		return err
	}

	err = s.repo.UpdatePayment(ctx, paymentID, payment)
	if err != nil {
		log.Error(
			"error to update payment",
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (s *Shop) GetPayments(ctx context.Context, storeID string) ([]entity.Payment, error) {
	log := zap.NewNop()

	admin := ctx.Value(config.EmailHeader)
	user, err := s.service.GetUser(ctx, &pb.GetUserRequest{
		Email: admin.(string),
	})
	if err != nil {
		log.Error(
			"error getting admin",
			zap.Error(err),
		)
		return nil, err
	}
	if !user.IsAdmin {
		log.Error(
			"unauthorized action",
		)
		return nil, errors.New("unauthorized action")
	}

	id, err := strconv.Atoi(storeID)
	if err != nil {
		log.Error(
			"error validating id",
			zap.Error(err),
		)
		return nil, err
	}

	result, err := s.repo.GetPayments(ctx, id)
	if err != nil {
		log.Error(
			"error to get payments",
			zap.Error(err),
		)
		return nil, err
	}

	return result, nil
}

func (s *Shop) SearchPayment(ctx context.Context, status, initialDate, endDate string) ([]entity.Payment, error) {
	log := zap.NewNop()

	admin := ctx.Value(config.EmailHeader)
	user, err := s.service.GetUser(ctx, &pb.GetUserRequest{
		Email: admin.(string),
	})
	if err != nil {
		log.Error(
			"error getting admin",
			zap.Error(err),
		)
		return nil, err
	}
	if !user.IsAdmin {
		log.Error(
			"unauthorized action",
		)
		return nil, errors.New("unauthorized action")
	}

	var init time.Time
	if initialDate != "" {
		init, err = time.Parse(time.RFC3339, initialDate)
		if err != nil {
			log.Error(
				"error validating initial date",
				zap.Error(err),
			)
			return nil, err
		}
	}

	var end time.Time
	if endDate != "" {
		end, err = time.Parse(time.RFC3339, endDate)
		if err != nil {
			log.Error(
				"error validating end date",
				zap.Error(err),
			)
			return nil, err
		}
	}

	result, err := s.repo.SearchPayment(ctx, status, init, end)
	if err != nil {
		log.Error(
			"error to search payments",
			zap.Error(err),
		)
		return nil, err
	}

	return result, nil
}
