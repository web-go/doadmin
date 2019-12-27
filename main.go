package main

import (
	"fmt"
	"log"

	"github.com/go-chi/cors"
	"github.com/go-rock/rock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/web-go/doadmin/engine"
	"github.com/web-go/doadmin/modules/config"
	"github.com/web-go/doadmin/pkg/validate"
)

func main() {
	r := rock.Default()
	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	// pp.Println(cors)
	// cors := cors.New(cors.Options{})
	// app.Use(cors.Handler)
	r.Mux().Use(cors.Handler)
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
	r.StaticFS("/static", "./static")
	r.StaticFS("/backend", "./static")
	r.Run(":5000")
}

func home(c rock.Context) {
	fmt.Println("admin home ")
	c.String(200, "xiao")
}
