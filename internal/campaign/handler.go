package campaign

import (
	"net/http"
	"time"
	"trinity/internal/model"
	"trinity/pkg/logger"
	"trinity/pkg/reason"
	"trinity/pkg/response"

	"github.com/gin-gonic/gin"
)

// Handler handles campaign-related requests
type Handler struct {
	service Service
	logger  logger.Logger
}

// NewHandler creates a new Campaign handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
		logger:  logger.NewLogger("campaignHandler"),
	}
}

// RegisterRoutes registers the campaign routes with the Gin router
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/", h.CreateCampaign)
	rg.GET("/", h.ListCampaigns)
	rg.POST("/:id/vouchers", h.GenerateVouchers)
}

// CreateCampaign godoc
// @Summary Create a new campaign
// @Description Create a new promotional campaign
// @Tags Campaign
// @Accept  json
// @Produce  json
// @Param campaign body campaign.CreateCampaignRequest true "Campaign Data"
// @Success 201 {object} model.Campaign
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /campaigns [post]
func (h *Handler) CreateCampaign(c *gin.Context) {
	var req CreateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		msg := reason.InvalidRequestFormat.Message()
		h.logger.Errorf("%s: %v", msg, err)
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: msg})
		return
	}

	// Parse dates
	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		msg := "Invalid start date format."
		h.logger.Errorf("%s: %v", reason.InvalidRequestFormat.Message(), err)
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: msg})
		return
	}

	endDate, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		msg := "Invalid end date format."
		h.logger.Errorf("%s: %v", reason.InvalidRequestFormat.Message(), err)
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: msg})
		return
	}

	campaign := model.Campaign{
		Name:        req.Name,
		Discount:    req.Discount,
		MaxUsers:    req.MaxUsers,
		UsedUsers:   0,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: req.Description,
	}

	id, err := h.service.CreateCampaign(&campaign)
	if err != nil {
		msg := reason.InternalServerError.Message()
		h.logger.Errorf("%s: %v", msg, err)
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: msg})
		return
	}
	campaign.Id = id
	c.JSON(http.StatusCreated, campaign)
}

// GenerateVouchers godoc
// @Summary Generate vouchers for a campaign
// @Description Generate vouchers for the specified campaign
// @Tags Campaign
// @Accept  json
// @Produce  json
// @Param id path string true "Campaign ID"
// @Param request body campaign.GenerateVouchersRequest true "Number of vouchers to generate"
// @Success 200 {array} model.Voucher
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /campaigns/{id}/vouchers [post]
func (h *Handler) GenerateVouchers(c *gin.Context) {
	campaignID := c.Param("id")

	var req GenerateVouchersRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		msg := reason.InvalidRequestFormat.Message()
		h.logger.Errorf("%s: %v", msg, err)
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: msg})
		return
	}

	if req.Count <= 0 {
		msg := reason.InvalidRequest.Message()
		h.logger.Errorf("%s: count must be greater than zero", msg)
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: msg})
		return
	}

	vouchers, err := h.service.GenerateVouchers(campaignID, req.Count)
	if err != nil {
		msg := reason.InternalServerError.Message()
		h.logger.Errorf("%s: %v", msg, err)
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: msg})
		return
	}

	c.JSON(http.StatusOK, vouchers)
}

// ListCampaigns godoc
// @Summary List all campaigns
// @Description Retrieve a list of all promotional campaigns
// @Tags Campaign
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Campaign
// @Failure 500 {object} response.ErrorResponse
// @Router /campaigns [get]
func (h *Handler) ListCampaigns(c *gin.Context) {
	campaigns, err := h.service.ListCampaigns()
	if err != nil {
		msg := reason.InternalServerError.Message()
		h.logger.Errorf("%s: %v", msg, err)
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: msg})
		return
	}

	c.JSON(http.StatusOK, campaigns)
}
