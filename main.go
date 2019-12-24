package main

import (
	"fmt"
	"log"

	"github.com/go-rock/rock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/web-go/doadmin/engine"
	"github.com/web-go/doadmin/modules/config"
	"github.com/web-go/doadmin/pkg/validate"
)

func main() {
	r := rock.Default()
	validate.InitBinding()
	db, err := gorm.Open("sqlite3", "docms.db")
	if err != nil {
		log.Printf("DB数据库启动异常%s", err)
	}
	defer db.Close()
	config := config.Config{
		Prefix: "/api/v1/sys",
		DB:     db,
	}
	eng := engine.Default()
	eng.SetConfig(config).Use(r)
	eng.Run()
	admin := eng.Router()
	admin.GET("/home", home)
	r.Run(":4000")
}

func home(c rock.Context) {
	fmt.Println("admin home ")
	c.String(200, "xiao")
}
