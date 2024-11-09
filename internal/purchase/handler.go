package purchase

import (
	"net/http"
	"trinity/internal/model"
	"trinity/pkg/logger"
	"trinity/pkg/reason"
	"trinity/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
	logger  logger.Logger
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
		logger:  logger.NewLogger("purchaseHandler"),
	}
}

// RegisterRoutes registers the purchase routes with the Gin router
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/", h.ProcessPurchase)
}

// ProcessPurchase godoc
// @Summary Process a subscription purchase
// @Description Process a subscription purchase with optional voucher code
// @Tags Purchase
// @Accept  json
// @Produce  json
// @Param request body purchase.ProcessPurchaseRequest true "Purchase data"
// @Success 200 {object} model.Purchase
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /purchases [post]
func (h *Handler) ProcessPurchase(c *gin.Context) {
	var req struct {
		UserId      string                 `json:"user_id"`
		Plan        model.SubscriptionPlan `json:"plan"`
		VoucherCode string                 `json:"voucher_code,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		msg := reason.InvalidRequestFormat.Message()
		h.logger.Errorf("%s: %v", msg, err)
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: msg})
		return
	}

	if req.UserId == "" || req.Plan == "" {
		msg := reason.InvalidRequest.Message()
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: msg})
		return
	}

	purchase, err := h.service.ProcessPurchase(req.UserId, req.Plan, req.VoucherCode)
	if err != nil {
		msg := reason.InternalServerError.Message()
		h.logger.Errorf("%s: %v", msg, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, purchase)
}
