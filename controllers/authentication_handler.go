package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("bebasapasaja")
var tokenName = "token"

type Claims struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserID   int    `json:"user_id"`
	UserType int    `json:"user_type"`
	jwt.StandardClaims
}

func generateToken(w http.ResponseWriter, email, password string, userID int, userType int) {
	tokenExpiryTime := time.Now().Add(10 * time.Minute)

	claims := &Claims{
		Email:    email,
		Password: password,
		UserID:   userID,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiryTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Print(token)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     tokenName,
		Value:    signedToken,
		Expires:  tokenExpiryTime,
		Secure:   false,
		HttpOnly: true,
	})

}

func resetUserToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     tokenName,
		Value:    "",
		Expires:  time.Now(),
		Secure:   false,
		HttpOnly: true,
	})
}

func Authenticate(next http.HandlerFunc, accessType int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isValidToken := validateUserToken(w, r, accessType)
		if !isValidToken {
			sendUnAuthorizedResponse(w)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func validateUserToken(w http.ResponseWriter, r *http.Request, accessType int) bool {
	isAccessTokenValid, id, email, userID, userType := validateTokenFromCookies(r)
	fmt.Print(id, email, userType, userID, accessType, isAccessTokenValid)

	if isAccessTokenValid {
		isUserValid := userType == accessType
		fmt.Print(isUserValid)
		if isUserValid {
			return true
		}
	}
	return false
}

func validateTokenFromCookies(r *http.Request) (bool, string, string, int, int) {
	if cookie, err := r.Cookie(tokenName); err == nil {
		accessToken := cookie.Value
		accessClaims := &Claims{}
		parsedToken, err := jwt.ParseWithClaims(accessToken, accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err == nil && parsedToken.Valid {
			return true, accessClaims.Email, accessClaims.Password, accessClaims.UserID, accessClaims.UserType
		}
	}
	return false, "", "", -1, -1
}

func GetIDFromCookies(r *http.Request) (bool, int, error) {
	cookie, err := r.Cookie(tokenName)
	if err != nil {
		return false, -1, err
	}

	accessToken := cookie.Value
	accessClaims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(accessToken, accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil && !parsedToken.Valid {
		return false, -1, err
	}
	return true, accessClaims.UserID, nil
}
