package middleware

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"gopkg.in/guregu/null.v4"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, "Request does not contain an access token")
			c.Abort()
			return
		}
		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(401, "Invalid access token")
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

var jwtKey = viper.GetString("SECRET")

type JWTClaim struct {
	ID       null.Int    `json:"id"`
	Username null.String `json:"username"`
	jwt.StandardClaims
}

func GenerateJWT(id null.Int, username null.String, expiration time.Time) (tokenString string, err error) {
	claims := &JWTClaim{
		ID:       id,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func ValidateToken(signedToken string) (claims *JWTClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("Could not parse claims")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("Token expired")
		return
	}
	return claims, nil
}
