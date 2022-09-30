package middlewares

import (
	"os"
	"time"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	Id  string `json:"id"`
	Rol uint8  `json:"rol"`
	jwt.StandardClaims
}

func CreateToken(id string, rol uint8) (string, time.Time, error) {
	// obteniendo secret de variable de entorno
	secretVal, defined := os.LookupEnv("SECRET_JWT")
	if !defined {
		secretVal = "secreto"
	}

	expiresJWT := time.Now().Add(2 * 24 * time.Hour)
	secret := []byte(secretVal)
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
		// obteniedo el header access-token
		var tkn string = ""
		for name, values := range c.Request().Header {
			if name == "Access-Token" {
				tkn = string(values[0])
				break
			}
		}
		// if tkn == "" {
		// 	return c.String(400, "error: token not provided")
		// }
		// fmt.Println("itrisdfL", tkn)

		// obteniendo secret de variable de entorno
		secretVal, defined := os.LookupEnv("SECRET_JWT")
		if !defined {
			secretVal = "secreto"
		}
		// verificando token
		secret := []byte(secretVal)
		claims := &JwtCustomClaims{}
		token, err := jwt.ParseWithClaims(tkn, claims, func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				return c.JSON(400, config.SetResError(400, "firma del token, no autorizado", err.Error()))
			}
			return c.JSON(400, config.SetResError(400, "token invalido", err.Error()))
		}
		if !token.Valid {
			return c.JSON(400, config.SetResError(400, "token no autorizado", ""))
		}
		// verificando que el usuario existe
		_, err = models.GetUserByIDStr(claims.Id)
		if err != nil {
			return c.JSON(400, config.SetResError(400, "Error: Id del token incorrecto", err.Error()))
		}
		// creando supervariables echo
		c.Set("id", claims.Id)
		c.Set("rol", claims.Rol)
		// fin del middleware
		return next(c)
	}
}

func IsBoss(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		rol := c.Get("rol").(uint8)
		if rol == 1 {
			return next(c)
		}
		return c.JSON(400, config.SetRes(400, "Error: El usuario no es AdminBoss."))
	}
}
func IsBossOrAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		rol := c.Get("rol").(uint8)
		if rol == 1 || rol == 2 {
			return next(c)
		}
		return c.JSON(400, config.SetRes(400, "Error: El usuario no es AdminBoss o AdminEmploye."))
	}
}
func IsBossOrAdminOrEmpl(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		rol := c.Get("rol").(uint8)
		if rol == 1 || rol == 2 || rol == 3 {
			return next(c)
		}
		return c.JSON(400, config.SetRes(400, "Error: El usuario no es AdminBoss o AdminEmploye o Employe."))
	}
}
func VerifyTokenWp(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obteniedo el header access-token
		var tkn string = ""
		for name, values := range c.Request().Header {
			if name == "Access-Token" {
				tkn = string(values[0])
				break
			}
		}
		// oteniendo env
		tkn2, _ := os.LookupEnv("WP_VERIFY_TOKEN")

		if tkn != tkn2 {
			return c.JSON(400, config.SetRes(400, "Error: No tienes permiso"))
		}

		// fin del middleware
		return next(c)
	}
}
