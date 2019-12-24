package models

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
	// Password             string `gorm:"-" json:"password" binding:"required,eqfield=PasswordConfirmation"`
	// PasswordConfirmation string `gorm:"-" json:"password_confirmation" binding:"required"`
	Nickname string `json:"nickname" gorm:"default:'DoadminUser'"`

	// Role     []Role           `json:"role" gorm:"many2many:roles_users"`
	// Enforcer *casbin.Enforcer `json:"-" inject:""`
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
	return DB.Model(&m).Update(m).Error
}

// 根据Id查询
func (m *User) Get() error {
	err := DB.First(&m).Error
	return err
}

// 根据Id删除
func (m *User) Delete() error {
	return DB.Delete(&m).Error
}
