package middlewares

import (
	"log"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnvLocal() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading Env local")
	}
}

// func SeparatePhoneString(phoneString string) (int int) {

// }

func GetPhoneString(codePhone, phone int) string {
	codePhoneStr, phoneStr := strconv.Itoa(codePhone), strconv.Itoa(phone)
	return codePhoneStr + phoneStr
}
