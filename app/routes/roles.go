package routes

import (
	"github.com/go-rock/rock"
	"github.com/web-go/doadmin/app/models"
	"github.com/web-go/doadmin/pkg/utils"
)

func ListRole(c rock.Context) {
	var ms models.Roles
	size := c.MustQueryInt("size", 10)
	repo := models.Repo{
		Ctx:          c,
		Result:       &ms,
		DB:           models.DB.Preload("Menus").Preload("Apis"),
		Pagination:   models.Pagination{PageSize: size},
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
	// if b := models.DB.Where("id = ?", id).First(m).RecordNotFound(); b {
	// 	utils.NotFound(c, "记录不存在")
	// 	return
	// }

	if err := m.Get(); err != nil {
		utils.Fail(c, err.Error())
		return
	}

	c.JSON(200, rock.M{"role": m})
}
