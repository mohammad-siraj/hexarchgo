package main

import (
	"github.com/mohammad-siraj/hexarchgo/internal/libs/cloud"
	user "github.com/mohammad-siraj/hexarchgo/internal/user/serverless"
)

func main() {
	cloud.StartLambda(user.HandlerUserCreation)
}
