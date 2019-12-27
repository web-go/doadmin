package routes

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-rock/rock"
	"github.com/web-go/doadmin/app/models"
	"github.com/web-go/doadmin/middleware"
	"github.com/web-go/doadmin/pkg/utils"
)

func Login(c rock.Context) {
	user := &models.LoginModel{}
	if err := c.ShouldBindJSON(user); err != nil {
		utils.Fail(c, "用户名或密码错误")
		return
	}
	if exit, u := user.Login(); exit {
		tokenNext(c, u)
		return
	}

	utils.Fail(c, "用户名或密码错误")
}

//登录以后签发jwt
func tokenNext(c rock.Context, user *models.User) {

	j := &middleware.JWT{SigningKey: []byte("docms")} // 唯一签名
	clams := middleware.CustomClaims{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		StandardClaims: jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000), // 签名生效时间
			// ExpiresAt: int64(time.Now().Unix() + 60*60*24*7), // 过期时间 一周
			ExpiresAt: int64(time.Now().Unix() + 60*60*24*7), // 过期时间 一周
			Issuer:    "doadmin",                             //签名的发行者
		},
	}
	token, err := j.CreateToken(clams)
	if err != nil {
		utils.Fail(c, "获取token失败")
	} else {
		utils.Success(c, rock.M{"user": user, "token": token, "expiresAt": clams.StandardClaims.ExpiresAt * 2000})
	}
}
