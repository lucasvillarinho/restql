package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Handler contains the necessary services
type Handler struct {
	userService *UserService
}

// NewHandler creates a new instance of Handler
func NewHandler() *Handler {
	return &Handler{
		userService: NewUserService(),
	}
}

// GetUsers returns the list of users with RestQL filters
func (h *Handler) GetUsers(c echo.Context) error {
	// Get query parameters from Echo and pass to service
	users, sqlQuery, args, err := h.userService.GetUsers(c.Request().Context(), c.QueryParams())
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Return response with query info for debugging
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":  users,
		"count": len(users),
		"query": sqlQuery,
		"args":  args,
	})
}
