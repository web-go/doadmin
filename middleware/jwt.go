package middleware

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-rock/rock"
)

func JWTAuth() rock.HandlerFunc {
	return func(c rock.Context) {

		token := c.Request().Header.Get("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(401, rock.M{"errors": "未登录或非法访问"})
			return
		}
		j := NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				c.AbortWithStatusJSON(401, rock.M{"errors": "授权已过期"})
				return
			}
			c.AbortWithStatusJSON(401, rock.M{"errors": err.Error()})
			return
		}
		fmt.Println("xiao")
		// pp.Println(inject.Obj.Enforcer.GetGroupingPolicy())
		// pp.Println(claims.Username, c.Request().URL.Path, c.Request().Method)
		// if b, err := inject.Obj.Enforcer.Enforce(claims.Username, c.Request().URL.Path, c.Request().Method); err != nil {
		// 	utils.CodeMsg(c, 403, "登录用户 校验权限失败")
		// 	c.Abort()
		// 	return
		// } else if !b {
		// 	utils.CodeMsg(c, 403, "登录用户 没有权限")
		// 	c.Abort()
		// 	return
		// }
		c.Set("claims", claims)
	}
}

type JWT struct {
	SigningKey []byte
}

type CustomClaims struct {
	ID       uint64
	Username string
	Nickname string
	jwt.StandardClaims
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

var (
	TokenExpired     error  = errors.New("Token is expired")
	TokenNotValidYet error  = errors.New("Token not active yet")
	TokenMalformed   error  = errors.New("That's not even a token")
	TokenInvalid     error  = errors.New("Couldn't handle this token:")
	SignKey          string = "docms"
)

//获取token
func GetSignKey() string {
	return SignKey
}

// 这是SignKey
func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

//创建一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

//解析 token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid

	}

}

// 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
