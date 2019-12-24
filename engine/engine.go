package engine

import (
	"github.com/go-rock/rock"
	"github.com/web-go/doadmin/app/models"
	"github.com/web-go/doadmin/app/routes"
	"github.com/web-go/doadmin/modules/config"
)

type Engine struct {
	config config.Config
	router *rock.App
}

func Default() *Engine {
	return &Engine{router: routes.AdminRoutes()}
}

// setConfig set the config of engine.
func (eng *Engine) SetConfig(cfg config.Config) *Engine {
	eng.config = cfg
	return eng
}

func (eng *Engine) Use(r *rock.App) {
	r.Mount(eng.config.Prefix, eng.router)
}

func (eng *Engine) Router() *rock.App {
	return eng.router
}

func (eng *Engine) Run() {
	models.Migrate(models.InitDB(eng.config))
}
