package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoggir/re-path/redirect-service/internal/dto"
	"github.com/hoggir/re-path/redirect-service/internal/repository"
	"github.com/hoggir/re-path/redirect-service/internal/service"
)

type RedirectHandler struct {
	redirectService   service.RedirectService
	clickEventService service.ClickEventService
}

func NewRedirectHandler(
	redirectService service.RedirectService,
	clickEventService service.ClickEventService,
) *RedirectHandler {
	return &RedirectHandler{
		redirectService:   redirectService,
		clickEventService: clickEventService,
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
		dto.ErrorResponse(c, http.StatusBadRequest, "Short url is required", nil)
		return
	}

	url, err := h.redirectService.GetURL(c.Request.Context(), shortUrl)
	if err != nil {
		log.Printf("❌ Failed to get URL for shortUrl %s: %v", shortUrl, err)

		// Handle specific errors dengan HTTP status yang sesuai
		if errors.Is(err, repository.ErrURLExpired) {
			dto.ErrorResponse(c, http.StatusGone, err.Error(), nil) // 410 Gone
			return
		}
		if errors.Is(err, repository.ErrURLInactive) {
			dto.ErrorResponse(c, http.StatusForbidden, err.Error(), nil) // 403 Forbidden
			return
		}
		if errors.Is(err, repository.ErrURLNotFound) {
			dto.ErrorResponse(c, http.StatusNotFound, err.Error(), nil) // 404 Not Found
			return
		}

		dto.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve URL", nil)
		return
	}

	go func() {
		ctx := context.Background()

		if err := h.redirectService.IncrementClickCount(ctx, shortUrl); err != nil {
			log.Printf("⚠️  Failed to increment click count for shortUrl %s: %v", shortUrl, err)
		}

		if err := h.clickEventService.TrackClick(ctx, c, shortUrl); err != nil {
			log.Printf("⚠️  Failed to track click for shortUrl %s: %v", shortUrl, err)
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
		dto.ErrorResponse(c, http.StatusBadRequest, "Short code is required", nil)
		return
	}

	url, err := h.redirectService.GetURL(c.Request.Context(), shortCode)
	if err != nil {
		// Handle specific errors dengan HTTP status yang sesuai
		if errors.Is(err, repository.ErrURLExpired) {
			dto.ErrorResponse(c, http.StatusGone, err.Error(), nil) // 410 Gone
			return
		}
		if errors.Is(err, repository.ErrURLInactive) {
			dto.ErrorResponse(c, http.StatusForbidden, err.Error(), nil) // 403 Forbidden
			return
		}
		if errors.Is(err, repository.ErrURLNotFound) {
			dto.ErrorResponse(c, http.StatusNotFound, err.Error(), nil) // 404 Not Found
			return
		}

		// Generic error
		dto.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve URL", nil)
		return
	}

	response := dto.RedirectResponse{
		// ShortCode:   url.ShortCode,
		OriginalURL: url.OriginalURL,
		// ClickCount:  url.ClickCount,
	}

	dto.SuccessResponse(c, http.StatusOK, "URL info retrieved successfully", response)
}
