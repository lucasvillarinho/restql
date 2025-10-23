package main

import "github.com/labstack/echo/v4"

// setupRoutes configures all application routes
func setupRoutes(e *echo.Echo) {
	h := NewHandler()

	e.GET("/users", h.GetUsers)
}
