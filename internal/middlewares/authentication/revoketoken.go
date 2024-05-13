package authentication

import (
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"os"

	mdl "nagarjuna2323/books_api/internal/models"
	"net/http"
	"time"
)

// RevokeToken revokes the provided token and adds it to the blacklist
func RevokeToken(c *gin.Context) {
	// Get token to revoke from request
	token := c.GetHeader("Authorization")
	if token == "" || token == "Bearer <access_token>" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token to revoke is missing"})
		return
	}
	// Add token to blacklist
	blacklistEntry := mdl.TokenBlacklist{
		Token:     token,
		Reason:    "Revoked by user",
		ExpiresAt: time.Now().AddDate(0, 0, 7),
		CreatedAt: time.Now(),
	}

	// Open or create the blacklist CSV file
	file, err := os.OpenFile("blacklist.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open blacklist file"})
		return
	}
	defer file.Close()

	// Write the blacklist entry to the CSV file
	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{blacklistEntry.Token, blacklistEntry.Reason, blacklistEntry.ExpiresAt.String(), blacklistEntry.CreatedAt.String()}
	if err := writer.Write(record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write to blacklist file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token revoked successfully"})
}
