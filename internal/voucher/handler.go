package voucher

import (
	"net/http"
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
		logger:  logger.NewLogger("voucherHandler"),
	}
}

// RegisterRoutes registers the voucher routes with the Gin router
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/redeem", h.RedeemVoucher)
}

// RedeemVoucher godoc
// @Summary Redeem a voucher
// @Description Redeem a voucher using code and user ID
// @Tags Voucher
// @Accept  json
// @Produce  json
// @Param request body voucher.RedeemVoucherRequest true "Voucher redemption data"
// @Success 200 {object} model.Voucher
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /vouchers/redeem [post]
func (h *Handler) RedeemVoucher(c *gin.Context) {
	var req RedeemVoucherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		msg := reason.InvalidRequestFormat.Message()
		h.logger.Errorf("%s: %v", msg, err)
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: msg})
		return
	}
	if req.Code == "" || req.UserId == "" {
		msg := reason.InvalidRequest.Message()
		h.logger.Errorf("%s: code or user_id missing", msg)
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: msg})
		return
	}
	voucher, err := h.service.RedeemVoucher(req.Code, req.UserId)
	if err != nil {
		msg := reason.InvalidToken.Message()
		h.logger.Errorf("%s: %v", msg, err)
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: msg})
		return
	}
	c.JSON(http.StatusOK, voucher)
}
