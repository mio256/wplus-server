package util

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	bearer := c.GetHeader("Authorization")

	re := regexp.MustCompile(`Bearer ([^ ]+)`)
	matches := re.FindStringSubmatch(bearer)
	if len(matches) < 2 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
		return
	}
	token := matches[1]

	if _, err := ParseToken(token); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
		})
		c.Abort()
		return
	}

	c.Next()
}
