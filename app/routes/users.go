package routes

import (
	"github.com/go-rock/rock"
	"github.com/web-go/doadmin/app/models"
	"github.com/web-go/doadmin/pkg/utils"
)

func ListUser(c rock.Context) {
	var ms models.Users
	limit := c.MustQueryInt("limit", 10)
	repo := models.Repo{
		Ctx:          c,
		Result:       &ms,
		DB:           models.DB.Preload("Roles"),
		Pagination:   models.Pagination{PageSize: limit},
		AutoResponse: true,
		ApplyWhere:   true,
	}
	repo.Fetch()
	// utils.Success(c, rock.M{"users": ms, "pagination": repo.Pagination})
}

func CreateUser(c rock.Context) {
	m := &models.UserModel{}
	if err := c.ShouldBindJSON(m); err != nil {
		utils.Error(c, err)
		return
	}
	password, err := models.GeneratePassword(m.Password)
	if err != nil {
		utils.Error(c, err)
		return
	}
	user := &models.User{Username: m.Username, Nickname: m.Nickname, PasswordDigest: password}
	if err := c.ShouldBindJSON(user); err != nil {
		utils.Error(c, err)
		return
	}
	if err := user.Add(); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, rock.M{"user": user})
}

func UpdateUser(c rock.Context) {
	id := c.Param("id")
	m := &models.UpdateUserModel{}
	if err := c.ShouldBindJSON(m); err != nil {
		utils.Error(c, err)
		return
	}
	user := &models.User{}
	if b := models.DB.Where("id = ?", id).First(user).RecordNotFound(); b {
		utils.NotFound(c, "记录不存在")
		return
	}

	user.Username = m.Username
	user.Nickname = m.Nickname

	if m.Password != "" {
		password, _ := models.GeneratePassword(m.Password)
		user.PasswordDigest = password
	}

	if err := c.ShouldBindJSON(user); err != nil {
		utils.Error(c, err)
		return
	}
	if err := user.Update(); err != nil {
		utils.Error(c, err)
		return
	}
	utils.Success(c, rock.M{"user": user})
}

func DeleteUser(c rock.Context) {
	id := c.MustParamInt("id", 0)
	m := &models.User{}
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
