package auth

import (
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	m "github.com/me/todo-go-server/src/models"
	s "github.com/me/todo-go-server/src/shared"
)

type login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func GinJwtMiddlewareHandler() *jwt.GinJWTMiddleware {
	jwt_secret := os.Getenv("JWT_SECRET")
	return &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte(jwt_secret),
		Timeout:    time.Duration(24*365) * time.Hour,
		MaxRefresh: time.Hour,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*m.User); ok {
				return jwt.MapClaims{
					"email": v.Email,
					"ID":    v.ID,
					"role":  "user", // TODO: v.Role
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			var user m.User
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			db := s.GetDB()
			db.Where("email = ?", loginVals.Email).First(&user)

			if user.CheckPassword(loginVals.Password) != nil {
				return nil, jwt.ErrFailedAuthentication
			} else {
				return &m.User{Email: user.Email, ID: user.ID}, nil
			}
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// if v, ok := data.(*User); ok && v.UserName == "admin" {
			// 	return true
			// }

			// jwtClaims := jwt.ExtractClaims(c)
			// fmt.Println(jwtClaims["userName"])
			// fmt.Println(c.Request.Header.Get("Authorization"))

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
