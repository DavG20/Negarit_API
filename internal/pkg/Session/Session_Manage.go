package session

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Session struct {
	jwt.StandardClaims
	Email string
}

type CookieHandler struct{}

func NewCookieHanedler() *CookieHandler {
	return &CookieHandler{}
}

func (cookieHandler *CookieHandler) GetSession(session *Session) (cookie http.Cookie, errs error) {
	expirtionTime := time.Now().Add(24 * time.Hour)

	session = &Session{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirtionTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, session)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWTKEY")))
	if err != nil {
		log.Fatal("error while token signing")
	}
	cookie = http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirtionTime,
	}
	return cookie, nil
}

func ValidateCookie() bool {
	return true
}
