package todos

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	s "github.com/me/todo-go-server/src/shared"
)

// Fetch all Todos
func FetchTodos(c *gin.Context) {
	db := s.GetDB()
	var todos []TodoModel
	if err := db.Find(&todos).Error; err != nil {
		if len(todos) <= 0 {
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Something went wrong"})
	} else {
		var _todos []TransformedTodo
		for _, item := range todos {
			_todos = append(_todos, TransformedTodo{ID: item.ID, Title: item.Title, Completed: item.Completed})
		}
		c.JSON(http.StatusOK, _todos)
	}
}

// Fetch single Todo
func FetchSingleTodo(c *gin.Context) {
	db := s.GetDB()
	var todo TodoModel
	todoID := c.Param("id")

	if err := db.First(&todo, todoID).Error; err != nil {
		if todo.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Something went wrong"})
	} else {
		c.JSON(http.StatusOK, todo)
	}
}

// Create new todo
func CreateTodo(c *gin.Context) {
	db := s.GetDB()
	var todo TodoModel
	c.BindJSON(&todo)
	if err := db.Save(&todo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Something went wrong"})
	} else {
		c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Todo item created successfully!", "todo": todo})
	}
}

// Update :id selected Todo
func UpdateTodo(c *gin.Context) {
	db := s.GetDB()
	var jsonTodo, todo TodoModel
	c.BindJSON(&jsonTodo)
	todoID := c.Param("id")

	if err := db.First(&todo, todoID).Error; err != nil {
		if todo.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Something went wrong"})
	}

	if err := db.Model(&todo).Update(map[string]interface{}{"title": jsonTodo.Title, "Completed": jsonTodo.Completed}).Error; err != nil {
		if todo.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Something went wrong"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo updated successfully!", "todo": todo})
	}
}

// Delete :id selected Todo
func DeleteTodo(c *gin.Context) {
	db := s.GetDB()
	var todo TodoModel
	todoID := c.Param("id")

	if err := db.First(&todo, todoID).Error; err != nil {
		if todo.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Something went wrong"})
	}

	if err := db.Delete(&todo).Error; err != nil {
		if todo.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Something went wrong"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo deleted successfully!", "deletedID": todo.ID})
	}
}
