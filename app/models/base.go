package models

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/web-go/doadmin/modules/config"
)

var DB *gorm.DB
var validate *validator.Validate

type BaseModel struct {
	ID        uint64    `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at" sql:"DEFAULT:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" sql:"DEFAULT:current_timestamp"`
}

//初始化数据库并产生数据库全局变量
func InitDB(config config.Config) *gorm.DB {
	DB = config.DB
	DB.LogMode(true)
	return DB
}

func TableName(name string) string {
	return fmt.Sprintf("%s%s", "sys_", name)
}

//注册数据库表专用
func Migrate(db *gorm.DB) {
	// db.DropTable(User{}, Role{}, Menu{})
	db.AutoMigrate(User{}, Role{}, Menu{}, Api{})
}
