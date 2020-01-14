package models

import (
	"log"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/jinzhu/gorm"
)

type UserModel struct {
	Username             string `json:"username" binding:"required"`
	Password             string `json:"password" binding:"required,eqfield=PasswordConfirmation"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required"`
	Nickname             string `json:"nickname" gorm:"default:'DoadminUser'"`
}

type UpdateUserModel struct {
	Username             string `json:"username" binding:"required"`
	Password             string `json:"password" binding:"eqfield=PasswordConfirmation"`
	PasswordConfirmation string `json:"password_confirmation"`
	Nickname             string `json:"nickname"`
}

type Users []User

type User struct {
	BaseModel
	Username       string `json:"username" binding:"required,uniq"`
	PasswordDigest string `json:"password_digest"`
	Nickname       string `json:"nickname" gorm:"default:'DoadminUser'"`

	Roles    []Role           `json:"roles" gorm:"many2many:sys_roles_users"`
	Enforcer *casbin.Enforcer `json:"-" inject:""`
}

// 表名
func (User) TableName() string {
	return TableName("users")
}

// 新增
func (m *User) Add() error {
	return DB.Create(&m).Error
}

// 修改
func (m *User) Update() error {
	m.UpdatedAt = time.Now()
	return DB.Model(&m).Update(m).Error
}

// 根据Id查询
func (m *User) Get() error {
	err := DB.Preload("Roles.Menus").Preload("Roles.Apis").First(&m).Error
	return err
}

// 根据Id删除
func (m *User) Delete() error {
	// return DB.Delete(&m).Error
	return DB.Transaction(func(tx *gorm.DB) error {
		// 删除关联的角色
		if err := tx.Model(m).Association("Roles").Clear().Error; err != nil {
			return err
		}

		if err := tx.Delete(&m).Error; err != nil {
			return err
		}
		return nil
	})
}

// 加载所有用户角色策略
func (m *User) LoadAllPolicy() error {
	ms := Users{}
	if err := DB.Find(&ms).Error; err != nil {
		return err
	}
	for _, user := range ms {
		if err := m.LoadPolicy(user.ID); err != nil {
			return err
		}
	}
	log.Printf("用户角色关系:%+v", m.Enforcer.GetGroupingPolicy())
	return nil
}

// 加载用户角色策略
func (m *User) LoadPolicy(id uint64) error {
	user := User{}
	user.ID = id
	if err := user.Get(); err != nil {
		return err
	}
	// uid := strconv.FormatUint(id, 10)
	// pp.Println(m, uid)
	m.Enforcer.DeleteRolesForUser(user.Username)
	for _, role := range user.Roles {
		m.Enforcer.AddRoleForUser(user.Username, role.Name)
	}
	log.Printf("更新用户角色关系:%+v", m.Enforcer.GetGroupingPolicy())
	return nil
}
