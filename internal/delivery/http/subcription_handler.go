package http

import (
	"eff-subscriptions/internal/domain/models"
	"eff-subscriptions/internal/repository"
	"eff-subscriptions/internal/validator"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

// createSubscription godoc
// @Summary Create a new subscription
// @Description Create a new subscription with the input payload
// @Tags subscriptions
// @Accept  json
// @Produce  json
// @Param input body models.CreateSubscriptionRequest true "Subscription object"
// @Success 201 {object} models.SubscriptionResponse
// @Failure 400 {object} errorResponse
// @Failure 422 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /v1/subscriptions [post]
func (h *Handler) createSubscription(c *gin.Context) {
	var input models.CreateSubscriptionRequest

	err := c.BindJSON(&input)
	if err != nil {
		h.badRequestResponse(c, err)
		return
	}

	subscription := &models.Subscription{
		ServiceName: input.ServiceName,
		Price:       input.Price,
		UserID:      input.UserID,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
	}

	v := validator.New()

	if models.ValidateSubscription(v, subscription); !v.Valid() {
		h.failedValidationResponse(c, v.Errors)
		return
	}

	err = h.subscriptionService.Insert(subscription)
	if err != nil {
		h.serverErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, models.SubscriptionResponse{Subscription: subscription})
}

// readSubscription godoc
// @Summary Get subscription
// @Description Return subscription by id
// @Tags subscriptions
// @Accept  json
// @Produce  json
// @Param id path int true "ID subscription"
// @Success 200 {object} models.SubscriptionResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /v1/subscriptions/{id} [get]
func (h *Handler) readSubscription(c *gin.Context) {
	id, err := readIDParam(c)
	if err != nil {
		h.badRequestResponse(c, err)
		return
	}

	subscription, err := h.subscriptionService.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			h.notFoundResponse(c)
		default:
			h.serverErrorResponse(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, models.SubscriptionResponse{Subscription: subscription})
}

// updateSubscription godoc
// @Summary Update subscription
// @Description Update subscription by id
// @Tags subscriptions
// @Accept  json
// @Produce  json
// @Param id path int true "ID subscription"
// @Param input body models.UpdateSubscriptionRequest true "New data"
// @Success 200 {object} models.SubscriptionResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 422 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /v1/subscriptions/{id} [patch]
func (h *Handler) updateSubscription(c *gin.Context) {
	id, err := readIDParam(c)
	if err != nil {
		h.badRequestResponse(c, err)
		return
	}

	subscription, err := h.subscriptionService.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			h.notFoundResponse(c)
		default:
			h.serverErrorResponse(c, err)
		}
		return
	}

	if c.GetHeader("X-Expected-Version") != "" {
		if strconv.Itoa(subscription.Version) != c.GetHeader("X-Expected-Version") {
			h.editConflictResponse(c)
			return
		}
	}

	var input models.UpdateSubscriptionRequest

	err = c.BindJSON(&input)
	if err != nil {
		h.badRequestResponse(c, err)
		return
	}

	if input.ServiceName != nil {
		subscription.ServiceName = *input.ServiceName
	}
	if input.Price != nil {
		subscription.Price = input.Price
	}
	if input.UserID != nil {
		subscription.UserID = *input.UserID
	}
	if input.StartDate != nil {
		subscription.StartDate = *input.StartDate
	}
	if input.EndDate != nil {
		subscription.EndDate = input.EndDate
	}

	v := validator.New()

	if models.ValidateSubscription(v, subscription); !v.Valid() {
		h.failedValidationResponse(c, v.Errors)
		return
	}

	err = h.subscriptionService.Update(subscription)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrEditConflict):
			h.editConflictResponse(c)
		default:
			h.serverErrorResponse(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, models.SubscriptionResponse{Subscription: subscription})
}

// deleteSubscription godoc
// @Summary Delete subscription
// @Description Delete subscription by id
// @Tags subscriptions
// @Accept  json
// @Produce  json
// @Param id path int true "ID subscription"
// @Success 200 {object} models.DataResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /v1/subscriptions/{id} [delete]
func (h *Handler) deleteSubscription(c *gin.Context) {
	id, err := readIDParam(c)
	if err != nil {
		h.badRequestResponse(c, err)
		return
	}

	err = h.subscriptionService.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			h.notFoundResponse(c)
		default:
			h.serverErrorResponse(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: "subscription successfully deleted"})
}

// listSubscriptions godoc
// @Summary Subscriptions list
// @Description Return subscriptions list with pagination and search
// @Tags subscriptions
// @Accept  json
// @Produce  json
// @Param service_name query string false "service name"
// @Param price query int false "price"
// @Param user_id query string false "user id"
// @Param start_date query string false "start date"
// @Param page query int false "page number"
// @Param page_size query int false "items limit on page"
// @Param sort query string false "sort field" Enums(id, service_name, price, start_date, -id, -service_name, -year, -price, -start_date) default(id)
// @Success 200 {object} models.SubscriptionsListResponse
// @Failure 422 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /v1/subscriptions [get]
func (h *Handler) listSubscriptions(c *gin.Context) {
	var input struct {
		ServiceName string            `json:"service_name"`
		Price       int               `json:"price"`
		UserID      uuid.UUID         `json:"user_id"`
		StartDate   models.CustomDate `json:"start_date"`
		models.Filters
	}

	v := validator.New()

	input.ServiceName = readString(c, "service_name", "")
	input.Price = readInt(c, "price", -1, v)
	input.UserID = readUUID(c, "user_id", uuid.Nil, v)
	input.StartDate = readDate(c, "start_date", models.CustomDate(time.Time{}), v)

	input.Filters.Page = readInt(c, "page", 1, v)
	input.Filters.PageSize = readInt(c, "page_size", 20, v)

	input.Filters.Sort = readString(c, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "service_name", "price", "start_date", "-id", "-service_name", "-year", "-price", "-start_date"}

	if models.ValidateFilters(v, input.Filters); !v.Valid() {
		h.failedValidationResponse(c, v.Errors)
		return
	}

	subscriptions, metadata, err := h.subscriptionService.GetAll(input.ServiceName, input.Price, input.UserID,
		input.StartDate, input.Filters)
	if err != nil {
		h.serverErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, models.SubscriptionsListResponse{Subscription: subscriptions, Metadata: metadata})

}

// sumSubscriptionsPrice godoc
// @Summary Sums up subscriptions prices
// @Description Sums up subscription prices over a date range
// @Tags subscriptions
// @Accept  json
// @Produce  json
// @Param service_name query string false "service name"
// @Param user_id query string false "user id"
// @Param start_date query string true "start date"
// @Param end_date query string true "end date"
// @Success 200 {object} models.DataResponse
// @Failure 422 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /v1/sum-subscriptions-price [get]
func (h *Handler) sumSubscriptionsPrice(c *gin.Context) {
	var input struct {
		UserID      uuid.UUID         `json:"user_id"`
		ServiceName string            `json:"service_name"`
		StartDate   models.CustomDate `json:"start_date"`
		EndDate     models.CustomDate `json:"end_date"`
	}

	v := validator.New()

	input.UserID = readUUID(c, "user_id", uuid.Nil, v)
	input.ServiceName = readString(c, "service_name", "")
	input.StartDate = readDate(c, "start_date", models.CustomDate(time.Time{}), v)
	input.EndDate = readDate(c, "end_date",
		models.CustomDate(time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)),
		v)

	if input.StartDate.Time().After(input.EndDate.Time()) {
		v.AddError("start_date", "must be before end_date")
	}

	if !v.Valid() {
		h.failedValidationResponse(c, v.Errors)
		return
	}

	sum, err := h.subscriptionService.GetSubscriptionsSum(input.UserID, input.ServiceName, input.StartDate, input.EndDate)
	if err != nil {
		h.serverErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, models.DataResponse{Data: sum})
}
