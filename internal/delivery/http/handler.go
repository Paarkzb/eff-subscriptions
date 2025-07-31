package http

import (
	"eff-subscriptions/internal/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log/slog"

	_ "eff-subscriptions/docs"
)

type Handler struct {
	log                 *slog.Logger
	subscriptionService *service.SubscriptionService
}

func NewHandler(log *slog.Logger, subscriptionService *service.SubscriptionService) *Handler {
	return &Handler{
		log:                 log,
		subscriptionService: subscriptionService,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	mux := gin.Default()

	mux.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	mux.GET("/v1/subscriptions", h.listSubscriptions)
	mux.POST("/v1/subscriptions", h.createSubscription)
	mux.GET("/v1/subscriptions/:id", h.readSubscription)
	mux.PATCH("/v1/subscriptions/:id", h.updateSubscription)
	mux.DELETE("/v1/subscriptions/:id", h.deleteSubscription)

	mux.GET("/v1/sum-subscriptions-price", h.sumSubscriptionsPrice)

	return mux
}
