package routes

import (
	"github.com/go-rock/rock"
	"github.com/web-go/doadmin/app/models"
	"github.com/web-go/doadmin/pkg/utils"
)

func ListMenu(c rock.Context) {
	var ms models.Menus
	size := c.MustQueryInt("size", 10)
	repo := models.Repo{
		Ctx:          c,
		Result:       &ms,
		DB:           models.DB.Where("parent_id = ?", 0),
		Pagination:   models.Pagination{PageSize: size},
		AutoResponse: false,
		ApplyWhere:   true,
	}
	repo.Fetch()
	var err error
	for i := 0; i < len(ms); i++ {
		m := &ms[i]
		err = m.GetBaseChildrenList()
	}
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, rock.M{"data": ms, "pagination": repo.Pagination})
}

func CreateMenu(c rock.Context) {
	menu := &models.Menu{}
	if err := c.ShouldBindJSON(menu); err != nil {
		utils.Error(c, err)
		return
	}
	if err := menu.Add(); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, rock.M{"menu": menu})
}

func UpdateMenu(c rock.Context) {
	id := c.MustParamInt("id", 0)
	m := &models.Menu{}
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

	c.JSON(200, rock.M{"menu": m})
}

func ShowMenu(c rock.Context) {
	id := c.MustParamInt("id", 0)
	m := &models.Menu{}
	m.ID = uint64(id)
	if err := m.Get(); err != nil {
		utils.Fail(c, err.Error())
		return
	}

	c.JSON(200, rock.M{"menu": m})
}

func DeleteMenu(c rock.Context) {
	id := c.MustParamInt("id", 0)
	m := &models.Menu{}
	m.ID = uint64(id)
	if b := models.DB.Where("id = ?", id).First(m).RecordNotFound(); b {
		utils.NotFound(c, "记录不存在")
		return
	}
	if err := m.Delete(); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, rock.M{"msg": "ok"})
}
