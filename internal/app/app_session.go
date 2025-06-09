package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"

	appConstants "dessert-ordering-go-system/internal/app_constants"
	utils "dessert-ordering-go-system/internal/utils"
)

type ApplicationSession struct {
	*scs.SessionManager
}

func NewApplicationSession(s *scs.SessionManager) *ApplicationSession {
	return &ApplicationSession{s}
}

func openSession(loggers *ApplicationLoggers, redisStore *redisstore.RedisStore) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.HttpOnly = true
	secureCookies := true
	if secureEnv := os.Getenv("SECURE_COOKIES"); secureEnv != "" {
		parsedBool, err := strconv.ParseBool(secureEnv)
		if err != nil {
			loggers.Info.Printf("Warning: SECURE_COOKIES environment variable '%s' is not a valid boolean (expected 'true' or 'false'). Defaulting to secure cookies (true).", secureEnv)
		} else {
			secureCookies = parsedBool
		}
	}
	sessionManager.Cookie.Secure = secureCookies
	sessionManager.Store = redisStore
	return sessionManager
}

func (s *ApplicationSession) GetAuthUserID(ctx context.Context) int {
	return s.GetInt(ctx, appConstants.Auth_User_ID)
}
func (s *ApplicationSession) RemoveAuthUserID(ctx context.Context) {
	s.Remove(ctx, appConstants.Auth_User_ID)
}
func (s *ApplicationSession) SetAuthUserID(ctx context.Context, userID int) {
	s.Put(ctx, appConstants.Auth_User_ID, userID)
}

func (s *ApplicationSession) GetCsrfToken(ctx context.Context) string {
	token := s.GetString(ctx, appConstants.X_CSRF_Token)
	if token == "" {
		csrfToken, err := utils.GenerateRandomString(32)
		if err != nil {
			log.Printf("ERROR: ApplicationSession.GetCsrfToken - utils.GenerateRandomString: %v", err)
		} else {
			s.Put(ctx, appConstants.X_CSRF_Token, csrfToken)
			token = csrfToken
		}
	}
	return token
}
func (s *ApplicationSession) RemoveCsrfToken(ctx context.Context) {
	s.Remove(ctx, appConstants.X_CSRF_Token)
}
func (s *ApplicationSession) SetCsrfToken(ctx context.Context, token string) {
	s.Put(ctx, appConstants.X_CSRF_Token, token)
}

func (s *ApplicationSession) PopFlashError(ctx context.Context) string {
	return s.PopString(ctx, appConstants.Flash_Error)
}

func (s *ApplicationSession) SetFlashError(ctx context.Context, err string) {
	s.Put(ctx, appConstants.Flash_Error, err)
}
