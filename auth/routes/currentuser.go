package routes

import (
	"auth/services"

	//TODO(Check an update to jwt of echo using V4, currently on v3
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RegisterCurrentUserHandler(e *echo.Echo) {
	e.GET("/api/users/currentuser", func(context echo.Context) error {

		user := context.Get("user").(*jwt.Token)
		claims := user.Claims.(*services.MyCustomClaims)

		return context.JSON(http.StatusOK, claims)

	}, services.CustomJWTMiddleware)
}
