package main

import (
	// Hola tio asdkasjkdaskd asdasd
	// "fmt"
	"net/http"
	"os"

	// jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	// "github.com/appleboy/gin-jwt"

	a "github.com/me/todo-go-server/src/auth"
	s "github.com/me/todo-go-server/src/shared"
	t "github.com/me/todo-go-server/src/todos"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&t.TodoModel{})
}

var identityKey = "id"

// User demo
type User struct {
	UserName  string
	FirstName string
	LastName  string
}

// func helloHandler(c *gin.Context) {
// 	claims := jwt.ExtractClaims(c)
// 	user, _ := c.Get(identityKey)
// 	c.JSON(200, gin.H{
// 		"userID":   claims["id"],
// 		"userName": user.(*User).UserName,
// 		"text":     "Hello World.",
// 	})
// }

func main() {

	db := s.Init()
	Migrate(db)
	defer db.Close()

	router := gin.Default()

	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	authMiddleware := a.GinJwtMiddlewareHandler()

	auth := router.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	}

	v1 := router.Group("/api/v1/todos")
	v1.Use(authMiddleware.MiddlewareFunc())
	{
		v1.POST("/", t.CreateTodo)
		v1.GET("/", t.FetchTodos)
		v1.GET("/:id", t.FetchSingleTodo)
		v1.PUT("/:id", t.UpdateTodo)
		v1.DELETE("/:id", t.DeleteTodo)
	}

	router.POST("/login", authMiddleware.LoginHandler)

	router.Run(":" + os.Getenv("PORT"))
}
