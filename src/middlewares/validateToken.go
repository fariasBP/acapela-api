package middlewares

import (
	"time"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	Id  string `json:"id"`
	Rol uint8  `json:"rol"`
	jwt.StandardClaims
}

const SECRET = "secreto"

func CreateToken(id string, rol uint8) (string, time.Time, error) {
	expiresJWT := time.Now().Add(30 * time.Minute)
	secret := []byte(SECRET)
	claims := &JwtCustomClaims{
		id,
		rol,
		jwt.StandardClaims{
			ExpiresAt: expiresJWT.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, tokenErr := token.SignedString(secret)
	return tokenString, expiresJWT, tokenErr
}

func ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Verificando que la cookie token exista
		tkn, err := c.Cookie("token")
		if err != nil {
			return c.JSON(400, "error: token not provided")
		}
		receivedToken := tkn.Value
		if receivedToken == "" {
			return c.JSON(400, "error: token not provided")
		}
		// verificando token
		secret := []byte(SECRET)
		claims := &JwtCustomClaims{}
		token, err := jwt.ParseWithClaims(receivedToken, claims, func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				return c.JSON(400, config.SetRes(400, "error: token unauthorized"))
			}
			return c.JSON(400, config.SetRes(400, "error: token invalidate"))
		}
		if !token.Valid {
			return c.JSON(400, config.SetRes(400, "error: token unauthorized"))
		}
		// creando supervariables echo
		c.Set("id", claims.Id)
		c.Set("rol", claims.Rol)
		// siguiente
		return next(c)
	}
}

func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		rol := c.Get("rol").(uint8)
		if rol == 1 {
			return next(c)
		}
		return c.JSON(400, config.SetRes(400, "error: user is not admin or empl"))
	}
}
func IsAdminOrEmpl(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		rol := c.Get("rol").(uint8)
		if rol == 1 || rol == 2 {
			return next(c)
		}
		return c.JSON(400, config.SetRes(400, "error: user is not admin or empl"))
	}
}
