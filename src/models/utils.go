package models

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pwd string, c echo.Context) (string, error) {
	pwdHashingByte, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(pwdHashingByte), err
}
