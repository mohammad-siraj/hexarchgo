package adapters

import (
	"net/http"

	"github.com/mohammad-siraj/hexarchgo/internal/libs/middleware"
)

func (u *UserHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	token, err := middleware.CreateToken("hello")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("Authorization", "Bearer "+token)
	w.Write([]byte("Registered user"))
}
