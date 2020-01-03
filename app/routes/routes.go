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
	r.GET("/users/profile", Profile)

	r.GET("/roles", ListRole)
	r.POST("/roles", CreateRole)
	r.GET("/roles/{id}", ShowRole)
	r.PUT("/roles/{id}", UpdateRole)
	r.DELETE("/roles/{id}", DeleteRole)
	r.POST("/roles/{id}/menus", AddRoleMenus)

	r.GET("/apis", ListApi)
	r.POST("/apis", CreateApi)
	r.GET("/apis/{id}", ShowApi)
	r.PUT("/apis/{id}", UpdateApi)
	r.DELETE("/apis/{id}", DeleteApi)

	r.GET("/menus", ListMenu)
	r.POST("/menus", CreateMenu)
	r.GET("/menus/{id}", ShowMenu)
	r.PUT("/menus/{id}", UpdateMenu)
	r.DELETE("/menus/{id}", DeleteMenu)

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
