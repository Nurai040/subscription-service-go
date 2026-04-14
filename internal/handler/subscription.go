package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"subscriptions-service/internal/logger"
	"subscriptions-service/internal/model"
	"subscriptions-service/internal/service"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	service *service.SubscriptionService
}

func NewSubscriptionHandler(s *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: s}
}

// @Summary Create subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body model.Subscription true "subscription"
// @Success 201
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(c *gin.Context) {
	logger.Log.Info("Create handler called")
	var sub model.Subscription

	if err := c.BindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.Create(sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// @Summary Get all subscriptions
// @Tags subscriptions
// @Produce json
// @Success 200
// @Router /subscriptions [get]
func (h *SubscriptionHandler) GetAll(c *gin.Context) {
	logger.Log.Info("GetAll handler called")
	subs, err := h.service.GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, subs)
}

// @Summary Get subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "subscription id"
// @Success 200
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	logger.Log.Info("GetByID handler called")
	idParam := c.Param("id")

	var id int
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	sub, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, sub)
}

// @Summary Update subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "subscription id"
// @Param input body model.Subscription true "subscription"
// @Success 200
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(c *gin.Context) {
	logger.Log.Info("Update handler called")

	var sub model.Subscription

	if err := c.BindJSON(&sub); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	idParam := c.Param("id")
	var id int
	fmt.Sscanf(idParam, "%d", &id)

	sub.ID = id

	err := h.service.Update(sub)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(200)
}

// @Summary Delete subscription
// @Tags subscriptions
// @Param id path int true "subscription id"
// @Success 204
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c *gin.Context) {
	logger.Log.Info("Delete handler called")

	idParam := c.Param("id")

	var id int
	fmt.Sscanf(idParam, "%d", &id)

	err := h.service.Delete(id)
	if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
        c.JSON(404, gin.H{"error": "not found"})
        return
    }

    c.JSON(500, gin.H{"error": "internal error"})
    return
}

	c.Status(204)
}

// @Summary Get total sum
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "user id"
// @Param service_name query string false "service name"
// @Param from query string false "from date"
// @Param to query string false "to date"
// @Success 200
// @Router /subscriptions/summary [get]
func (h *SubscriptionHandler) GetSummary(c *gin.Context) {
	logger.Log.Info("GetSummary handler called")

	userIDStr := c.Query("user_id")
	serviceName := c.Query("service_name")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	var userID uuid.UUID
	if userIDStr != "" {
		parsed, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
			return
		}
		userID = parsed
	}

	var from, to *time.Time

	if fromStr != "" {
		t, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date"})
			return
		}
		from = &t
	}

	if toStr != "" {
		t, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to date"})
			return
		}
		to = &t
	}

	sum, err := h.service.GetTotalSum(userID, serviceName, from, to)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"total_sum": sum})
}