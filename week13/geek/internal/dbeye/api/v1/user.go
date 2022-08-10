package v1

import (
	"fmt"
	"geek/internal/dbeye/pkg/app"
	"geek/internal/dbeye/pkg/ce"
	"geek/internal/dbeye/service"
	"geek/pkg/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type User struct {
}

func NewUser() User {
	return User{}
}

func (u User) GetUser(c *gin.Context) {
	var (
		err  error
		code int
	)
	req := service.GetUserReq{}
	appG := app.Gin{c}
	data := make(map[string]interface{})
	code = ce.SUCCESS
	err = c.ShouldBindQuery(&req)
	if err != nil {
		code = ce.INVALID_PARAMS
		err = fmt.Errorf("参数验证失败:%w", err)
		log.Logger.Error("", zap.Error(err))
		appG.Response(http.StatusOK, code, data)
		return
	}
	us := service.NewUserService()
	err, code, data = us.GetUserInfo(req)
	if err != nil {
		log.Logger.Error("", zap.Error(err))
		appG.Response(http.StatusOK, code, data)
		return
	}
	appG.Response(http.StatusOK, code, data)
}

func (u User) SyncUser(c *gin.Context) {
	var (
		err  error
		code int
	)
	req := service.SyncUserReq{}
	appG := app.Gin{c}
	data := make(map[string]interface{})
	code = ce.SUCCESS
	err = c.ShouldBindJSON(&req)
	if err != nil {
		code = ce.INVALID_PARAMS
		err = fmt.Errorf("参数验证失败:%w", err)
		log.Logger.Error("", zap.Error(err))
		appG.Response(http.StatusOK, code, data)
		return
	}
	us := service.NewUserService()
	err, code = us.SyncUserInfo(req)
	if err != nil {
		log.Logger.Error("", zap.Error(err))
		appG.Response(http.StatusOK, code, data)
		return
	}
	appG.Response(http.StatusOK, code, data)
}
