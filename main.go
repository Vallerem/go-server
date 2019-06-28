package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	a "github.com/me/todo-go-server/src/auth"
	m "github.com/me/todo-go-server/src/models"
	r "github.com/me/todo-go-server/src/routes"
	s "github.com/me/todo-go-server/src/shared"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&m.User{}, &m.Todo{})
}

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
		v1.POST("/", r.CreateTodo)
		v1.GET("/", r.FetchTodos)
		v1.GET("/:id", r.FetchSingleTodo)
		v1.PUT("/:id", r.UpdateTodo)
		v1.DELETE("/:id", r.DeleteTodo)
	}

	router.POST("/login", authMiddleware.LoginHandler)
	router.POST("/signup", r.UsersRegistration)

	router.GET("/users/:id", r.GetUser)

	router.Run(":" + os.Getenv("PORT"))
}
