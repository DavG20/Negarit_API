package session

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Session struct {
	jwt.StandardClaims
	UserName string
}

type CookieHandler struct{}

func NewCookieHanedler() *CookieHandler {
	return &CookieHandler{}
}

func (cookieHandler *CookieHandler) GetCookie(session *Session) (cookie http.Cookie, errs error) {
	expirtionTime := time.Now().Add(2400 * time.Hour)

	// session = &Session{
	// 	StandardClaims: jwt.StandardClaims{
	// 		ExpiresAt: expirtionTime.Unix(),
	// 	},
	// }
	ss := jwt.StandardClaims{ExpiresAt: expirtionTime.Unix()}

	session.StandardClaims = ss
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, session)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWTKEY")))
	if err != nil {
		log.Fatal("error while token signing")
	}
	cookie = http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirtionTime,
		HttpOnly: true,
	}
	return cookie, nil
}

func (cookieHandler *CookieHandler) ValidateCookie(request *http.Request) (*Session, bool) {
	tokenCookie, err := request.Cookie("token")
	if err != nil {
		return nil, false
	}
	token := tokenCookie.Value
	fmt.Println(token, " token")
	session := &Session{}
	tkn, err := jwt.ParseWithClaims(token, session, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return "", nil
		}
		return []byte(os.Getenv("JWTKEY")), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, false
		}
		return nil, false
	}
	if !tkn.Valid {
		return nil, false
	}
	return session, true

}

func (cookieHandler *CookieHandler) RemoveCookie() (http.Cookie, error) {
	expirationTime := time.Unix(0, 0)
	session := &Session{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
		UserName: "",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, session)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWTKEY")))
	if err != nil {
		log.Println("error signing token")
		log.Println("error signing")
	}
	cookie := http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
	}

	return cookie, nil

}


