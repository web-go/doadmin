package models

import (
	"time"

	"github.com/casbin/casbin/v2"
)

type Roles []Role

// Role
type Role struct {
	BaseModel
	Name string `json:"name" binding:"required,uniq,lte=8"`

	Users    []User           `json:"-" gorm:"many2many:sys_roles_users"`
	Menus    []Menu           `json:"menus" gorm:"many2many:sys_menus_roles"`
	Apis     []Api            `json:"apis" gorm:"many2many:sys_apis_roles"`
	Enforcer *casbin.Enforcer `json:"-" inject:""`
}

// 表名
func (Role) TableName() string {
	return TableName("roles")
}

// 新增
func (m *Role) Add() error {
	return DB.Create(&m).Error
}

// 根据Id查询
func (m *Role) Get() error {
	return DB.Preload("Menus").Preload("Apis").First(&m).Error
}

// 修改
func (m *Role) Update() error {
	m.UpdatedAt = time.Now()
	if err := DB.Model(m).Update(m).Preload("Menus").Preload("Apis").Error; err != nil {
		return err
	}

	return nil
}

// 根据Id删除
func (m *Role) Delete() error {
	err := DB.Delete(&m).Error
	if err != nil {
		return err
	}
	return nil
}

// 加载角色权限策略
func (m *Role) LoadAllPolicy() error {
	ms := Roles{}
	if err := DB.Find(&ms).Error; err != nil {
		return err
	}
	for _, role := range ms {
		if err := m.LoadPolicy(role.ID); err != nil {
			return err
		}
	}
	return nil
}

func (m *Role) LoadPolicy(id uint64) error {
	role := Role{}
	role.ID = id
	if err := role.Get(); err != nil {
		return err
	}
	m.SetRole(role)
	return nil
}

func (m *Role) SetRole(role Role) error {
	// roleID := strconv.Itoa(int(role.ID))
	roleID := role.Name
	m.Enforcer.DeleteRole(roleID)
	for _, api := range role.Apis {

		if api.Path == "" {
			continue
		}
		m.Enforcer.AddPermissionForUser(roleID, api.Path, api.Method)
	}
	return nil
}
