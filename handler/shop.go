package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/restore/shop/config"
	"github.com/restore/shop/entity"
	"net/http"
)

type controller interface {
	CreateRequest(ctx context.Context, request *entity.Create) (string, error)
	UpdateRequest(ctx context.Context, id string, request *entity.Request) error
	SearchRequest(ctx context.Context, storeID, status, initialDate, endDate string) ([]entity.Request, error)

	CreatePayment(ctx context.Context, payment *entity.Payment) (int, error)
	UpdatePayment(ctx context.Context, id string, payment *entity.Payment) error
	GetPayments(ctx context.Context, storeID string) ([]entity.Payment, error)
	SearchPayment(ctx context.Context, status, initialDate, endDate string) ([]entity.Payment, error)
}

type Shop struct {
	controller controller
}

func NewShop(c controller) *Shop {
	return &Shop{
		controller: c,
	}
}

// CreateRequest creates a new Request.
func (s *Shop) CreateRequest(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), config.EmailHeader, c.GetHeader(config.EmailHeader))

	var request entity.Create
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	id, err := s.controller.CreateRequest(ctx, &request)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, struct {
		ID string
	}{
		ID: id,
	})
}

// UpdateRequest updates a Request.
func (s *Shop) UpdateRequest(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), config.EmailHeader, c.GetHeader(config.EmailHeader))

	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			"invalid ID",
		})
		return
	}

	var request entity.Request
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	err := s.controller.UpdateRequest(ctx, id, &request)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, struct{}{})
}

// SearchRequest searches for Requests.
func (s *Shop) SearchRequest(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), config.EmailHeader, c.GetHeader(config.EmailHeader))

	storeID := c.Param("storeID")
	if storeID == "" {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			"invalid ID",
		})
		return
	}

	result, err := s.controller.SearchRequest(
		ctx,
		storeID,
		c.Query("status"),
		c.Query("initialDate"),
		c.Query("endDate"),
	)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}

// CreatePayment creates a new Payment.
func (s *Shop) CreatePayment(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), config.EmailHeader, c.GetHeader(config.EmailHeader))

	var payment entity.Payment
	if err := c.BindJSON(&payment); err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	id, err := s.controller.CreatePayment(ctx, &payment)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, struct {
		ID int
	}{
		id,
	})
}

// UpdatePayment updates a Payment.
func (s *Shop) UpdatePayment(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), config.EmailHeader, c.GetHeader(config.EmailHeader))

	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			"invalid ID",
		})
		return
	}

	var payment entity.Payment
	if err := c.BindJSON(&payment); err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	err := s.controller.UpdatePayment(ctx, id, &payment)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, struct{}{})
}

// GetPayments searches for Payments of a store.
func (s *Shop) GetPayments(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), config.EmailHeader, c.GetHeader(config.EmailHeader))

	storeID := c.Param("storeID")
	if storeID == "" {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			"invalid ID",
		})
		return
	}

	result, err := s.controller.GetPayments(ctx, storeID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}

// SearchPayments searches for Payments.
func (s *Shop) SearchPayments(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), config.EmailHeader, c.GetHeader(config.EmailHeader))

	result, err := s.controller.SearchPayment(
		ctx,
		c.Query("status"),
		c.Query("initialDate"),
		c.Query("endDate"),
	)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}
