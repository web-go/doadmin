package config

import "github.com/jinzhu/gorm"

type Config struct {
	DB     *gorm.DB
	Prefix string
}
