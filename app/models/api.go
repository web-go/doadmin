package models

import "time"

type Apis []Api

// Api
type Api struct {
	BaseModel
	// Name        string `json:"name" binding:"required,uniq"`
	Path        string `json:"path" binding:"required"`
	Method      string `json:"method" binding:"required"`
	Group       string `json:"group" binding:"required"`
	Description string `json:"description" binding:"required,uniq"`
	Hidden      bool   `json:"hidden"`
}

// 表名
func (Api) TableName() string {
	return TableName("apis")
}

// 新增
func (m *Api) Add() error {
	return DB.Create(&m).Error
}

// 根据Id查询
func (m *Api) Get() error {
	return DB.First(&m).Error
}

// 修改
func (m *Api) Update() error {
	m.UpdatedAt = time.Now()
	return DB.Model(&m).Update(m).Error
}

// 获取列表
func (ms *Apis) List(e Api) error {
	return DB.Find(&ms, e).Error
}

// 获取总数
func (m *Api) Count() (int, error) {
	var count int
	return count, DB.Model(Menu{}).Where(&m).Count(&count).Error
}

// 根据Id删除
func (m *Api) Delete() error {
	return DB.Delete(&m).Error
}
