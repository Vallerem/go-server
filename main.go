package main

import (
	// "fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	// "github.com/appleboy/gin-jwt"
	t "github.com/me/todo-go-server/src/todos"
	s "github.com/me/todo-go-server/src/shared"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&t.TodoModel{})
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

	v1 := router.Group("/api/v1/todos")
	{
		v1.POST("/", t.CreateTodo)
		v1.GET("/", t.FetchTodos)
		v1.GET("/:id", t.FetchSingleTodo)
		v1.PUT("/:id", t.UpdateTodo)
		v1.DELETE("/:id", t.DeleteTodo)
	}

	router.Run(":" + os.Getenv("PORT"))
}
