package authentication

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"nagarjuna2323/books_api/infrastructure/secrets"
	L "nagarjuna2323/books_api/internal/middlewares/logger"
	mdl "nagarjuna2323/books_api/internal/models"
	"net/http"
	"strconv"
	"time"
)

func RefreshToken(c *gin.Context) {
	// Get refresh token from request
	refreshToken := c.GetHeader("Authorization")
	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token is missing"})
		return
	}

	// Parse refresh token
	token, err := jwt.Parse(refreshToken[len("Bearer "):], func(token *jwt.Token) (interface{}, error) {
		return []byte(secrets.BOOKS_DEV_API_SECRET_KEY), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	// Check if the token is expired
	exp := time.Unix(int64(claims["exp"].(float64)), 0)
	if exp.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token has expired"})
		return
	}

	// Get user ID from claims
	userID := claims["user_id"].(string)

	// Mock retrieving user data from CSV file based on user ID
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		L.BKSLog("E", "Failed to convert userID to integer:\n", err)
	}
	user := mdl.User{
		ID: uint(userIDInt),
		// Populate other user details as needed
	}

	// Generate new access token
	accessToken, err := GenerateToken(user)
	if err != nil {
		L.BKSLog("D", "Error Generating Access Token:\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	// Return new access token
	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}
