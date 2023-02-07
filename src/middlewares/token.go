package middlewares

import (
	"os"
	"strconv"
	"time"

	"github.com/fariasBP/acapela-api/src/config"
	"github.com/fariasBP/acapela-api/src/models"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	Id string `json:"id"`
	jwt.StandardClaims
}

func CreateToken(id string, rol uint8) (string, time.Time, error) {
	// obteniendo secret de variable de entorno
	secretVal, defined := os.LookupEnv("SECRET_JWT")
	if !defined {
		secretVal = "secreto"
	}

	// obteniendo la duracion del token y convirtiendo
	durationTknAuth, defined := os.LookupEnv("TOKEN_AUTH_DURATION")
	if !defined {
		durationTknAuth = "30"
	}
	durationTknAuthInt, err := strconv.Atoi(durationTknAuth)
	if err != nil {
		durationTknAuthInt = 30
	}

	// creando el token
	expiresJWT := time.Now().Add(time.Duration(durationTknAuthInt) * 24 * time.Hour)
	secret := []byte(secretVal)
	claims := &JwtCustomClaims{
		id,
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
		var shop string = ""
		for name, values := range c.Request().Header {
			if name == "Access-Token" {
				tkn = string(values[0])
			}
			if name == "Shop-Id" {
				shop = string(values[0])
			}
		}

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
		// creando supervariables echo (guardando id de usuario y tienda)
		c.Set("id", claims.Id)
		c.Set("shop", shop)
		// fin del middleware
		return next(c)
	}
}

// verificar si es dueño de la tienda
func IsOwnerShop(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, idShop := c.Get("id").(string), c.Get("shop").(string)

		if b := models.VerifyOwnerShop(id, idShop); b {
			return next(c)
		}
		return c.JSON(400, config.SetResError(400, "Error: Usted no es Owner.", "user id or store id are not related"))
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
