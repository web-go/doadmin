package routes

import (
	"github.com/go-rock/rock"
	"github.com/web-go/doadmin/app/models"
	"github.com/web-go/doadmin/middleware"
	"github.com/web-go/doadmin/pkg/utils"
)

func Profile(c rock.Context) {
	claims, _ := c.Get("claims")
	waitUse := claims.(*middleware.CustomClaims)
	m := models.User{}
	m.ID = waitUse.ID
	if err := m.Get(); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	var ms []interface{}
	var roles []string
	if m.Username == "admin" {
		var menus []models.Menu
		models.DB.Find(&menus)
		roles = []string{"超级管理员"}
		for _, menu := range menus {
			tmpMenu := rock.M{"path": menu.Path, "name": menu.Name, "component": menu.Component, "meta": rock.M{"title": menu.Title, "icon": menu.Icon}, "parent_id": menu.ParentID, "id": menu.ID}
			if utils.Contains(ms, tmpMenu) < 0 {
				ms = append(ms, tmpMenu)
			}
		}
	} else {
		for _, role := range m.Roles {
			roles = append(roles, role.Name)
			for _, menu := range role.Menus {
				tmpMenu := rock.M{"path": menu.Path, "name": menu.Name, "component": menu.Component, "meta": rock.M{"title": menu.Title, "icon": menu.Icon}, "parent_id": menu.ParentID, "id": menu.ID}
				if utils.Contains(ms, tmpMenu) < 0 {
					ms = append(ms, tmpMenu)
				}
			}
		}
	}
	// var basems models.Menus
	// for _, m := range ms {
	// 	if m.ParentID == 0 {
	// 		basems = append(basems, m)
	// 	}
	// }

	// c.JSON(200, rock.M{"user": rock.M{"id": m.ID, "username": m.Username, "nickname": m.Nickname}, "menus": ms})
	utils.Success(c, rock.M{"id": m.ID, "name": m.Username, "nickname": m.Nickname, "avatar": "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif", "menus": ms, "roles": roles})
}

// func UserMenus(c rock.Context) {
// 	claims, _ := c.Get("claims")
// 	waitUse := claims.(*middleware.CustomClaims)
// 	m := models.User{}
// 	m.ID = waitUse.ID
// 	if err := m.Get(); err != nil {
// 		utils.Fail(c, err.Error())
// 		return
// 	}

// 	c.JSON(200, rock.M{"user": m})
// }
