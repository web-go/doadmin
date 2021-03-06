package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Menus []Menu

// Menu
type Menu struct {
	BaseModel
	Name      string `json:"name" binding:"required,uniq"`
	Path      string `json:"path" binding:"required,uniq"`
	Title     string `json:"title" binding:"required"`
	Component string `json:"component" binding:"required"`
	Icon      string `json:"icon"`
	Position  int    `json:"position" binding:"required"`
	Hidden    bool   `json:"hidden"`
	ParentID  uint64 `json:"parent_id"`
	// Children menuSlice `gorm:"foreignkey:ParentID" json:"children"`
	Children []Menu `gorm:"-" json:"children"`
	Roles    []Role `json:"-" gorm:"many2many:sys_menus_roles"`
}

// 表名
func (Menu) TableName() string {
	return TableName("menus")
}

// 新增
func (m *Menu) Add() error {
	return DB.Create(&m).Error
}

// 根据Id查询
func (m *Menu) Get() error {
	return DB.First(&m).Error
}

// 修改
func (m *Menu) Update() error {
	m.UpdatedAt = time.Now()
	return DB.Save(&m).Error
}

// 获取总数
func (m *Menu) Count() (int, error) {
	var count int
	return count, DB.Model(Menu{}).Where(&m).Count(&count).Error
}

// 根据Id删除
func (m *Menu) Delete() error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 清除权限关联菜单
		DB.Model(&m).Association("Roles").Clear()
		// 设置子菜单parent_id 为 0
		if err := DB.Model(Menu{}).Where("parent_id = ?", m.ID).Updates(map[string]interface{}{"parent_id": 0}).Error; err != nil {
			return err
		}
		if err := DB.Delete(&m).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
}

//获取基础路由树
func (m *Menu) GetBaseMenuTree() (err error, menus []Menu) {
	err = DB.Where(" parent_id = ?", 0).Order("position", true).Find(&menus).Error
	for i := 0; i < len(menus); i++ {
		m := &menus[i]
		err = m.GetBaseChildrenList()
	}
	return err, menus
}

func (menu *Menu) GetBaseChildrenList() (err error) {
	err = DB.Where("parent_id = ?", menu.ID).Order("position", true).Find(&menu.Children).Error
	for i := 0; i < len(menu.Children); i++ {
		m := &menu.Children[i]
		err = m.GetBaseChildrenList()
	}
	return err
}
