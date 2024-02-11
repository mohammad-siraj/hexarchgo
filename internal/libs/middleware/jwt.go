package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	secretKeyPath    = "./data/secret/tokensecret" //os.Getenv("JWT_SECRET")
	hmacSampleSecret []byte
	tokenInfoKey     struct{}
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
		if _, err := w.Write([]byte(`{"error":"No Authorization header provided"}`)); err != nil {
			return
		}
		return
	}

	splitToken := strings.Split(authToken, " ")
	if len(splitToken) != 2 {
		w.WriteHeader(http.StatusProxyAuthRequired)
		r.Header.Add(AUTH_FAIL_HEADER, "true")
		if _, err := w.Write([]byte(`{"error":"No Authorization header invalid format"}`)); err != nil {
			return
		}
		return
	}

	err := verifyToken(splitToken[1])
	if err != nil {
		w.WriteHeader(http.StatusProxyAuthRequired)
		r.Header.Add(AUTH_FAIL_HEADER, "true")
		if _, err := w.Write([]byte(`{"error":"No Authorization header invalid"}`)); err != nil {
			return
		}
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

func grpcAuthFunc(ctx context.Context) (context.Context, error) {
	token, err := auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	tokenInfo, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return hmacSampleSecret, nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	ctx = logging.InjectFields(ctx, logging.Fields{"auth.sub", "claim"})
	return context.WithValue(ctx, tokenInfoKey, tokenInfo), nil
}

func GrpcAuthMiddleware(ctx context.Context) grpc.UnaryServerInterceptor {
	return selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(grpcAuthFunc), selector.MatchFunc(AuthSkip))
}

func AuthSkip(_ context.Context, c interceptors.CallMeta) bool {
	fmt.Println(c.FullMethod())
	return c.FullMethod() != "/service.User/RegisterUser"
}
