package ports

import (
	"github.com/mohammad-siraj/hexarchgo/internal/libs/http"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/logger"
	user "github.com/mohammad-siraj/hexarchgo/internal/user/driving/adapters"
)

func RegisterRequestHandlers(h http.IHttpClient, l logger.ILogger) {
	// register request handlers for all the available endpoints here
	userHandler := user.NewUserHandler(h, l)

	userSubRoute := h.NewSubGroup("/auth")
	{
		userSubRoute.Post("/register", userHandler.RegisterUserHandler)
	}
}
