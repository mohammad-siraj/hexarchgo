package adapters

import (
	"net/http"

	"github.com/mohammad-siraj/hexarchgo/internal/libs/logger"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/middleware"
)

func (u *UserHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	token, err := middleware.CreateToken("hello")
	if err != nil {
		if _, err := w.Write([]byte(err.Error())); err != nil {
			u.log.Error(r.Context(), "Failed to write response", logger.NewLogFieldInput("error", err))
			return
		}
		return
	}
	w.Header().Add("Authorization", "Bearer "+token)
	if _, err := w.Write([]byte("Registered user")); err != nil {
		u.log.Error(r.Context(), "Failed to write response", logger.NewLogFieldInput("error", err))
		return
	}
}
