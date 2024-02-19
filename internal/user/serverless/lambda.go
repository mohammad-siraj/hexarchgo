package serverless

import (
	"context"
	"fmt"

	user "github.com/mohammad-siraj/hexarchgo/internal/user/driving/adapters/proto/service"
)

func HandlerUserCreation(ctx context.Context, userRegister *user.UserRegisterRequest) (string, error) {
	if userRegister == nil {
		return "", fmt.Errorf("received nil event")
	}
	fmt.Println("Is invoked successfully")
	message := fmt.Sprintf("Hello %s!", userRegister.Name)
	return message, nil
}
