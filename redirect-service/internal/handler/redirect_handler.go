package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoggir/re-path/redirect-service/internal/config"
	"github.com/hoggir/re-path/redirect-service/internal/domain"
	"github.com/hoggir/re-path/redirect-service/internal/dto"
	"github.com/hoggir/re-path/redirect-service/internal/logger"
	"github.com/hoggir/re-path/redirect-service/internal/service"
)

type RedirectHandler struct {
	redirectService   service.RedirectService
	clickEventService service.ClickEventService
	config            *config.Config
	logger            logger.Logger
}

func NewRedirectHandler(
	redirectService service.RedirectService,
	clickEventService service.ClickEventService,
	cfg *config.Config,
	log logger.Logger,
) *RedirectHandler {
	return &RedirectHandler{
		redirectService:   redirectService,
		clickEventService: clickEventService,
		config:            cfg,
		logger:            log,
	}
}

// Redirect returns the original URL without redirecting
// @Summary Get original URL from short url
// @Description Returns the original URL information without performing a redirect
// @Tags Redirect
// @Accept json
// @Produce json
// @Param shortUrl path string true "Short Url"
// @Success 200 {object} dto.Response{data=dto.RedirectResponse}
// @Failure 404 {object} dto.Response
// @Router /r/{shortUrl} [get]
func (h *RedirectHandler) Redirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")

	if shortUrl == "" {
		dto.HandleError(c, domain.ErrMissingRequired.WithContext("field",
			"shortUrl"))
		return
	}
	if len(shortUrl) > 50 {
		dto.HandleError(c, domain.ErrInvalidInput.WithContext("field",
			"shortUrl").WithMessage("Short URL too long"))
		return
	}

	url, err := h.redirectService.GetURL(c.Request.Context(), shortUrl)
	if err != nil {
		h.logger.ErrorContext(c.Request.Context(), "failed to get URL", "shortUrl", shortUrl, "error", err)
		dto.HandleError(c, err)
		return
	}

	metadata := domain.ClickMetadata{
		ClientIP:  c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Referrer:  c.Request.Referer(),
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), h.config.Service.ClickTrackingTimeout)
		defer cancel()

		if err := h.clickEventService.TrackClick(ctx, metadata, shortUrl); err != nil {
			h.logger.WarnContext(ctx, "failed to track click", "shortUrl", shortUrl, "error", err)
		}
	}()

	response := dto.RedirectResponse{
		OriginalURL: url.OriginalURL,
	}

	dto.SuccessResponse(c, http.StatusOK, "URL retrieved successfully", response)
}

// GetURLInfo returns URL information without redirecting
// @Summary Get URL information
// @Tags Redirect
// @Param shortCode path string true "Short Url"
// @Success 200 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Router /api/info/{shortUrl} [get]
func (h *RedirectHandler) GetURLInfo(c *gin.Context) {
	shortCode := c.Param("shortCode")

	if shortCode == "" {
		dto.HandleError(c, domain.ErrMissingRequired.WithContext("field", "shortCode"))
		return
	}

	url, err := h.redirectService.GetURL(c.Request.Context(), shortCode)
	if err != nil {
		dto.HandleError(c, err)
		return
	}

	response := dto.RedirectResponse{
		OriginalURL: url.OriginalURL,
	}

	dto.SuccessResponse(c, http.StatusOK, "URL info retrieved successfully", response)
}
