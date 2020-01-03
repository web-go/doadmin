package routes

import (
	"github.com/go-rock/rock"
	"github.com/web-go/doadmin/app/models"
	"github.com/web-go/doadmin/pkg/inject"
	"github.com/web-go/doadmin/pkg/utils"
)

func ListRole(c rock.Context) {
	var ms models.Roles
	limit := c.MustQueryInt("limit", 10)
	repo := models.Repo{
		Ctx:          c,
		Result:       &ms,
		DB:           models.DB.Preload("Menus").Preload("Apis"),
		Pagination:   models.Pagination{PageSize: limit},
		AutoResponse: true,
		ApplyWhere:   true,
	}
	repo.Fetch()
	// utils.Success(c, rock.M{"Roles": ms, "pagination": repo.Pagination})
}

func CreateRole(c rock.Context) {
	role := &models.Role{}
	if err := c.ShouldBindJSON(role); err != nil {
		utils.Error(c, err)
		return
	}
	if err := role.Add(); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, rock.M{"role": role})
}

func UpdateRole(c rock.Context) {
	id := c.MustParamInt("id", 0)
	m := &models.Role{}
	m.ID = uint64(id)
	if err := m.Get(); err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	if err := c.ShouldBindJSON(m); err != nil {
		utils.Error(c, err)
		return
	}

	if err := m.Update(); err != nil {
		utils.Error(c, err)
		return
	}

	c.JSON(200, rock.M{"role": m})
}

func ShowRole(c rock.Context) {
	id := c.MustParamInt("id", 0)
	m := &models.Role{}
	m.ID = uint64(id)

	if err := m.Get(); err != nil {
		utils.Fail(c, err.Error())
		return
	}

	c.JSON(200, rock.M{"role": m})
}

func DeleteRole(c rock.Context) {
	id := c.MustParamInt("id", 0)
	m := &models.Role{}
	m.ID = uint64(id)

	if err := models.DB.Preload("Users").First(m).Error; err != nil {
		utils.Fail(c, err.Error())
		return
	}
	// 如果有用户关联，不可删除

	// if len(m.Users) > 0 {
	// 	utils.Fail(c, "删除失败：此角色有用户正在使用禁止删除")
	// 	return
	// }

	name := m.Name
	if err := m.Delete(); err != nil {
		utils.Error(c, err)
		return
	}
	inject.Obj.Enforcer.DeletePermissionsForUser(name)

	utils.Success(c, rock.M{"code": 0})
}

type AddMenuRoleInfo struct {
	Menus []models.Menu `json:"menus"`
}

func AddRoleMenus(c rock.Context) {
	id := c.MustParamInt("id", 0)
	m := &models.Role{}
	m.ID = uint64(id)
	if err := m.Get(); err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	addMenuRoleInfo := &AddMenuRoleInfo{}

	if err := c.ShouldBindJSON(addMenuRoleInfo); err != nil {
		utils.Error(c, err)
		return
	}

	// 设置role菜单
	if err := models.DB.Model(m).Association("Menus").Replace(addMenuRoleInfo.Menus).Error; err != nil {
		utils.Error(c, err)
		return
	}

	c.JSON(200, addMenuRoleInfo)
}
