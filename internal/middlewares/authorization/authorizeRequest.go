package authorization

import (
	"encoding/csv"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"nagarjuna2323/books_api/infrastructure/secrets"
	"net/http"
	"os"
	"time"
)

func AuthorizeRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}
		// Extract JWT token
		//0en from Authorization header
		tokenString = tokenString[len("Bearer "):]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secrets.BOOKS_DEV_API_SECRET_KEY), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// Check if the token is valid
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// Check token expiry
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			exp := time.Unix(int64(claims["exp"].(float64)), 0)
			if exp.Before(time.Now()) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
				c.Abort()
				return
			}
		}
		// Check if the token is blacklisted
		if isRevoked := IsTokenRevoked(tokenString); isRevoked {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
			c.Abort()
			return
		}
		// Set user information in the context for further processing
		c.Set("user", token.Claims)
		c.Next()
	}
}

// IsTokenRevoked checks if the given token is blacklisted in the CSV file
func IsTokenRevoked(token string) bool {
	// Open or create the blacklist CSV file
	file, err := os.Open("blacklist.csv")
	if err != nil {
		// If unable to open the CSV file, assume token is not revoked
		return false
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		// If an error occurred while reading the CSV file, assume token is not revoked
		return false
	}

	// Iterate through the records and check if the token exists in the CSV file
	for _, record := range records {
		if record[0] == token {
			// If the token is found in the CSV file, it is revoked
			return true
		}
	}

	// If the token is not found in the CSV file, it is not revoked
	return false
}
