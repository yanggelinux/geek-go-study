package middleware

import (
	"fmt"
	"geek/internal/pkg/app"
	"geek/internal/pkg/ce"
	"geek/pkg/log"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const XRequestIDKey = "X-RequestID-Key"

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			token string
			code  int
		)

		data := make(map[string]interface{})
		code = ce.SUCCESS
		appG := app.Gin{c}
		s, exist := c.GetQuery("token")
		if exist {
			token = s
		} else {
			token = c.GetHeader("X-Token")
		}
		if token == "" {
			code = ce.ERROR_AUTH_TOKEN_FAILED
		} else {
			_, err := app.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					code = ce.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
				default:
					code = ce.ERROR_AUTH_CHECK_TOKEN_FAILED
				}
			}
		}
		if code != ce.SUCCESS {
			appG.Response(http.StatusOK, code, data)
			c.Abort()
			return
		}
		c.Next()
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			data := make(map[string]interface{})
			code := ce.ERROR
			appG := app.Gin{c}
			if p := recover(); p != nil {
				//painc时打印出堆栈日志
				errMsg := fmt.Sprintf("%+v", p)
				err := errors.New(errMsg)
				log.Logger.Info("panic recover", zap.Any("panic", err))
				appG.Response(http.StatusOK, code, data)
				c.Abort()
			}
		}()
		c.Next()
	}
}

func ReqCostTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		//请求前获取当前时间
		nowTime := time.Now()

		//请求处理
		c.Next()

		//处理后获取消耗时间
		costTime := time.Since(nowTime)
		url := c.Request.URL.String()
		msg := fmt.Sprintf("the request URL %s cost %v", url, costTime)
		log.Logger.Info(msg)
	}
}
