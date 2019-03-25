package auth

import (
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// User demo
type User struct {
	UserName  string
	FirstName string
	LastName  string
}

func GinJwtMiddlewareHandler() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte("verySecreteKey"),
		Timeout:    time.Hour * 100,
		MaxRefresh: time.Hour,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					"userName": v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
				return &User{
					UserName:  userID,
					LastName:  "Test",
					FirstName: "User",
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// if v, ok := data.(*User); ok && v.UserName == "admin" {
			// 	return true
			// }

			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
				// Could add more stuff here...
			})
		},
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		}
	}
}
