package services

import (
	"auth/errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"time"
)

var MySigningKey = []byte(os.Getenv("JWT_KEY"))

type MyCustomClaims struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

var CustomJWTMiddleware = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey:  MySigningKey,
	TokenLookup: "cookie:jwt",
	Claims: &MyCustomClaims{},
	ErrorHandlerWithContext: func(err error, context echo.Context) error {
		unAuth := errors.NotAuthorizedError{Message: "Not Authorized"}
		return context.JSON(
			unAuth.Status(),
			unAuth.Errors(),
		)
	},
})

func GenerateToken(user MyCustomClaims) (string, error) {
	claims := MyCustomClaims{
		user.ID,
		user.Email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(MySigningKey)
	if err != nil {
		return "", err
	}

	return signedString, nil
}

func SetToken(signedString string, ctx echo.Context) {
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    signedString,
		Expires:  time.Now().Add(15 * time.Minute),
		Secure:   true,
		HttpOnly: true,
	}
	ctx.SetCookie(&cookie)
}

func InvalidateToken(ctx echo.Context) {
	cookie := http.Cookie{
		Name:     "jwt",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
	}
	ctx.SetCookie(&cookie)
}
