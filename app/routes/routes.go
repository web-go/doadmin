package routes

import (
	"fmt"

	"github.com/go-rock/rock"
	"github.com/web-go/doadmin/middleware"
)

func AdminRoutes() *rock.App {
	r := rock.Default()
	// r.Use(Auth())
	r.GET("/", Home)
	r.POST("/login", Login)

	r.Use(middleware.JWTAuth())
	r.GET("/users", ListUser)
	r.POST("/users", CreateUser)
	r.PUT("/users/{id}", UpdateUser)
	r.DELETE("/users/{id}", DeleteUser)
	return r
}
func Auth() rock.HandlerFunc {
	return func(c rock.Context) {
		fmt.Println("strtat =======")
		c.Next()
		fmt.Println("====end =======")
	}
}
func Home(c rock.Context) {
	c.JSON(200, rock.M{"ok": "ok"})
}
