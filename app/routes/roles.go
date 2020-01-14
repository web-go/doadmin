package routes

import (
	"github.com/go-rock/rock"
	"github.com/web-go/doadmin/app/models"
	"github.com/web-go/doadmin/pkg/inject"
	"github.com/web-go/doadmin/pkg/utils"
	"github.com/yudai/pp"
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

	utils.Success(c, rock.M{"role": m})
}

func ShowRole(c rock.Context) {
	id := c.MustParamInt("id", 0)
	m := &models.Role{}
	m.ID = uint64(id)

	if err := m.Get(); err != nil {
		utils.Fail(c, err.Error())
		return
	}

	utils.Success(c, rock.M{"role": m})
}

func DeleteRole(c rock.Context) {
	id := c.MustParamInt("id", 0)
	m := &models.Role{}
	m.ID = uint64(id)

	if err := models.DB.Preload("Users").First(m).Error; err != nil {
		utils.Fail(c, err.Error())
		return
	}

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

type AddApiRoleInfo struct {
	Apis []models.Api `json:"apis"`
}

func AddRoleApis(c rock.Context) {
	id := c.MustParamInt("id", 0)
	m := &models.Role{}
	m.ID = uint64(id)
	if err := m.Get(); err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	addApiRoleInfo := &AddApiRoleInfo{}

	if err := c.ShouldBindJSON(addApiRoleInfo); err != nil {
		utils.Error(c, err)
		return
	}

	// 设置role菜单
	if err := models.DB.Model(m).Association("Apis").Replace(addApiRoleInfo.Apis).Error; err != nil {
		utils.Error(c, err)
		return
	}

	inject.Obj.Enforcer.RemoveFilteredPolicy(0, m.Name)
	for _, api := range addApiRoleInfo.Apis {

		if api.Path == "" {
			continue
		}
		inject.Obj.Enforcer.AddPermissionForUser(m.Name, api.Path, api.Method)
	}

	pp.Println(inject.Obj.Enforcer.GetPolicy(), inject.Obj.Enforcer.GetGroupingPolicy())

	c.JSON(200, addApiRoleInfo)
}
