package app

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"time"
)

var (
	ErrTokenExpired = errors.New("Token is expired")
)

func (ctx *Context) verifyToken(tokenString string) (bool, *jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return ctx.TokenSignerPublicKey, nil
	})
	if err != nil {
		if token != nil && token.Claims != nil {
			if validationError, ok := token.Claims.Valid().(*jwt.ValidationError); ok {
				if validationError.Errors == jwt.ValidationErrorExpired {
					return false, nil, errors.New("Token is expired")
				}
			}
		}
		return false, nil, errors.Wrap(err, "Cannot parse access token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Printf("%v\n", claims)
		return true, &claims, nil
	}
	return false, nil, errors.New("Invalid access token")
}

func (ctx *Context) createToken(username string, userID int64, roles []string) (string, error) {
	ttl := 5 * time.Hour

	var claims jwt.MapClaims
	claims = jwt.MapClaims{
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(ttl).Unix(),
		"uid":  userID,
		"name": username,
		"role": roles,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	tokenString, err := token.SignedString(ctx.TokenSignerPrivateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
