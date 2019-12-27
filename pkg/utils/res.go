package utils

import (
	"net/http"

	"github.com/go-rock/rock"
	"github.com/web-go/doadmin/pkg/validate"
)

func Success(c rock.Context, data interface{}) {
	c.JSON(http.StatusOK, rock.M{"data": data})
}
func Error(c rock.Context, err error) {
	c.JSON(http.StatusBadRequest, validate.ValidatorError(err))
}
func Fail(c rock.Context, data interface{}) {
	c.JSON(http.StatusBadRequest, rock.M{"errors": data})
}
func NotFound(c rock.Context, data interface{}) {
	c.JSON(http.StatusNotFound, rock.M{"errors": data})
}
