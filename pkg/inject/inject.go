package inject

import (
	"log"

	"github.com/web-go/doadmin/app/models"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/facebookgo/inject"
)

var text = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act) || r.sub == "admin"
`

type Common struct {
	User *models.User `inject:""`
	Role *models.Role `inject:""`
	Menu *models.Menu `inject:""`
}

// Object 注入对象
type Object struct {
	Common   *Common
	Enforcer *casbin.Enforcer
}

var Obj *Object

// 初始化依赖注入
func init() {
	g := new(inject.Graph)

	log.Println("注入casbin")
	model, _ := model.NewModelFromString(text)
	enforcer, _ := casbin.NewEnforcer(model)
	_ = g.Provide(&inject.Object{Value: enforcer})

	Common := new(Common)
	_ = g.Provide(&inject.Object{Value: Common})

	if err := g.Populate(); err != nil {
		panic("初始化依赖注入发生错误：" + err.Error())
	}

	Obj = &Object{
		Enforcer: enforcer,
		Common:   Common,
	}
	return
}

// 加载casbin策略数据，包括角色权限数据、用户角色数据
func LoadCasbinPolicyData() error {
	c := Obj.Common

	err := c.Role.LoadAllPolicy()
	if err != nil {
		return err
	}
	err = c.User.LoadAllPolicy()
	if err != nil {
		return err
	}
	return nil
}
