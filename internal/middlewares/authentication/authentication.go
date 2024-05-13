package authentication

import (
	"github.com/dgrijalva/jwt-go"
	"nagarjuna2323/books_api/infrastructure/secrets"
	L "nagarjuna2323/books_api/internal/middlewares/logger"
	mdl "nagarjuna2323/books_api/internal/models"
	"strconv"
	"time"
)

func GenerateToken(user mdl.User) (string, error) {
	// Create the claims
	claims := mdl.Claims{
		UserID: strconv.Itoa(int(user.ID)),
		Email:  user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	jwtSecret := []byte(secrets.BOOKS_DEV_API_SECRET_KEY)
	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	L.BKSLog("D", "Token:"+L.PrintStruct(tokenString), err)
	return tokenString, nil
}
