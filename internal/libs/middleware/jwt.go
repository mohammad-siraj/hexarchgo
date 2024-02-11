package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	secretKeyPath    = "./data/secret/tokensecret" //os.Getenv("JWT_SECRET")
	hmacSampleSecret []byte
)

const (
	AUTH_FAIL_HEADER = "authFail"
)

func init() {
	if keyData, e := os.ReadFile(secretKeyPath); e == nil {
		hmacSampleSecret = keyData
	} else {
		panic(e)
	}
}

func AuthMiddleware(w http.ResponseWriter, r *http.Request) {
	// Check if the user is authenticated
	authToken := r.Header.Get("Authorization")
	if authToken == "" {
		w.WriteHeader(http.StatusProxyAuthRequired)
		r.Header.Add(AUTH_FAIL_HEADER, "true")
		w.Write([]byte(`{"error":"No Authorization header provided"}`))
		return
	}

	splitToken := strings.Split(authToken, " ")
	if len(splitToken) != 2 {
		w.WriteHeader(http.StatusProxyAuthRequired)
		r.Header.Add(AUTH_FAIL_HEADER, "true")
		w.Write([]byte(`{"error":"No Authorization header invalid format"}`))
		return
	}

	err := verifyToken(splitToken[1])
	if err != nil {
		w.WriteHeader(http.StatusProxyAuthRequired)
		r.Header.Add(AUTH_FAIL_HEADER, "true")
		w.Write([]byte(`{"error":"No Authorization header invalid"}`))
		return
	}
}

func CheckAuthenticationHeader(r *http.Request) bool {
	return (r.Header.Get(AUTH_FAIL_HEADER) != "false")
}

func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return hmacSampleSecret, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
