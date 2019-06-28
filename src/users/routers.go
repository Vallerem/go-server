package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	m "github.com/me/todo-go-server/src/models"
	s "github.com/me/todo-go-server/src/shared"
)

func UsersRegistration(c *gin.Context) {
	db := s.GetDB()
	var user m.User
	c.BindJSON(&user)

	password := user.PasswordHash
	user.SetPassword(password)

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": err})
	} else {
		c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Account created successfully!", "user": user})
	}
}
