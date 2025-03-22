package middleware

import (
	"log"
	"myproject/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func RBAC(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the database connection from the context
		db, exists := c.Get("db")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
			c.Abort()
			return
		}

		gormDB, ok := db.(*gorm.DB)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid database connection"})
			c.Abort()
			return
		}

		// Retrieve the user ID from the context
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Validate userID as a UUID
		_, err := uuid.Parse(userID.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		// Fetch the user with roles and permissions
		var user models.User
		if err := gormDB.Preload("Roles.Permissions").Where("id = ?", userID).First(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
			c.Abort()
			return
		}

		// Check if the user has the required permission
		hasPermission := false
		for _, role := range user.Roles {
			for _, perm := range role.Permissions {
				log.Println("Checking permission:", perm.Name)
				if perm.Name == permission {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func SetDBMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}
