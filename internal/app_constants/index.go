package app_constants

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

var (
	Auth_User_ID  = "Auth_User_ID"
	Flash_Error   = "flash_error"
	Jwt_Name      = "jwt_token"
	Jwt_Expiration = 1 * time.Hour
	SecureCookies = true
	X_CSRF_Token  = "X-CSRF-Token"
)

func GetSecureCookies() (bool, error) {
	secureCookies := true
	if secureEnv := os.Getenv("SECURE_COOKIES"); secureEnv != "" {
		parsedBool, err := strconv.ParseBool(secureEnv)
		if err != nil {
			return true, fmt.Errorf("warning: SECURE_COOKIES environment variable '%s' is not a valid boolean (expected 'true' or 'false'). Defaulting to secure cookies (true)", secureEnv)
		} else {
			secureCookies = parsedBool
		}
	}

	return secureCookies, nil
}
