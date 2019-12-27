package routes

import (
	"github.com/go-rock/rock"
	"github.com/web-go/doadmin/app/models"
	"github.com/web-go/doadmin/pkg/utils"
)

func ListApi(c rock.Context) {
	var ms models.Apis
	size := c.MustQueryInt("size", 10)
	repo := models.Repo{
		Ctx:          c,
		Result:       &ms,
		DB:           models.DB,
		Pagination:   models.Pagination{PageSize: size},
		AutoResponse: true,
		ApplyWhere:   true,
	}
	repo.Fetch()
}

func CreateApi(c rock.Context) {
	api := &models.Api{}
	if err := c.ShouldBindJSON(api); err != nil {
		utils.Error(c, err)
		return
	}
	if err := api.Add(); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, rock.M{"api": api})
}

func UpdateApi(c rock.Context) {
	id := c.MustParamInt("id", 0)
	m := &models.Api{}
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

	c.JSON(200, rock.M{"api": m})
}

func ShowApi(c rock.Context) {
	id := c.MustParamInt("id", 0)
	m := &models.Api{}
	m.ID = uint64(id)
	if err := m.Get(); err != nil {
		utils.Fail(c, err.Error())
		return
	}

	c.JSON(200, rock.M{"api": m})
}

func DeleteApi(c rock.Context) {
	id := c.MustParamInt("id", 0)
	m := &models.Api{}
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
