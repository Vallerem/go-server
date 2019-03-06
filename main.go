package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var db *gorm.DB

func init() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	dbusername := os.Getenv("DB_USERNAME")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	//open a db connection
	var err error
	dbString := fmt.Sprintf("host=localhost port=5432 user=%s dbname=%s password=%s sslmode=disable", dbusername, dbname, dbpassword)

	db, err = gorm.Open("postgres", dbString)
	if err != nil {
		panic(err)
	}

	log.Printf("Connected to DB")

	//Migrate the schema
	db.AutoMigrate(&todoModel{})
}

func main() {

	router := gin.Default()

	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")
	
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	v1 := router.Group("/api/v1/todos")
	{
		v1.POST("/", createTodo)
		v1.GET("/", fetchTodos)
		v1.GET("/:id", fetchSingleTodo)
		v1.PUT("/:id", updateTodo)
		v1.DELETE("/:id", deleteTodo)
	}

	defer db.Close()
	router.Run(":3000")
}

type (
	todoModel struct {
		gorm.Model
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}

	// Formated Todo
	transformedTodo struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
)

// Fetch all Todos
func fetchTodos(c *gin.Context) {
	var todos []todoModel
	db.Find(&todos)

	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	var _todos []transformedTodo
	for _, item := range todos {
		_todos = append(_todos, transformedTodo{ID: item.ID, Title: item.Title, Completed: item.Completed})
	}

	c.JSON(http.StatusOK, _todos)
}

// Fetch single Todo
func fetchSingleTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// Create new todo
func createTodo(c *gin.Context) {
	var todo todoModel
	c.BindJSON(&todo)
	db.Save(&todo)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Todo item created successfully!", "todo": todo})
}

// Update :id selected Todo
func updateTodo(c *gin.Context) {
	var jsonTodo, todo todoModel
	c.BindJSON(&jsonTodo)
	todoID := c.Param("id")

	db.First(&todo, todoID)
	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	db.Model(&todo).Updates(map[string]interface{}{"title": jsonTodo.Title, "Completed": jsonTodo.Completed})

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo updated successfully!", "todo": todo})
}

// Delete :id selected Todo
func deleteTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	db.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo deleted successfully!", "deletedID": todo.ID})
}
